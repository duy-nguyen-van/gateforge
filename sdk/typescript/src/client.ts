import {
  AdminApi,
  AuthApi,
  Configuration,
  FederationApi,
  HealthApi,
  MFAApi,
  OIDCApi,
  TenantsApi,
  UsersApi,
  WebAuthnApi,
  type HTTPHeaders,
  type Middleware,
  type ResponseContext,
} from './generated/src/index.js'
import { APIError, type Envelope, type Meta } from './errors.js'

/** Returns a bearer access token for authenticated requests. */
export type TokenProvider = () => string | Promise<string | null | undefined> | null | undefined

/** Optional single-flight refresh after a 401. */
export type RefreshAccessToken = () => Promise<string | null | undefined>

export type GateForgeClientOptions = {
  /** Absolute GateForge IAM base URL (e.g. https://iam.example.com). */
  baseUrl: string
  /** Custom fetch implementation (defaults to global fetch). */
  fetch?: typeof fetch
  /** Provider for Authorization: Bearer tokens. */
  getAccessToken?: TokenProvider
  /**
   * Fetch credentials mode. Defaults to `"include"` so SSO cookies
   * (`iam_session`) are sent on same-site / CORS-credentialed requests.
   */
  credentials?: RequestCredentials
  /** Optional default headers merged into every request. */
  headers?: HeadersInit
  /**
   * Called when a request returns 401. If provided, invoked after a failed
   * refresh (or immediately when `refreshAccessToken` is not set).
   */
  onUnauthorized?: () => void | Promise<void>
  /**
   * Single-flight refresh hook. On 401, called once; if it returns a token
   * (or updates whatever backs `getAccessToken`), the request is retried once.
   */
  refreshAccessToken?: RefreshAccessToken
}

function headersInitToRecord(headers?: HeadersInit): HTTPHeaders {
  if (!headers) {
    return {}
  }
  const out: HTTPHeaders = {}
  new Headers(headers).forEach((value, key) => {
    out[key] = value
  })
  return out
}

function isUnauthorizedResponse(response: Response | undefined): boolean {
  return response?.status === 401
}

/**
 * GateForge IAM client wrapping the generated OpenAPI fetch APIs.
 */
export class GateForgeClient {
  readonly baseUrl: string
  readonly fetch: typeof fetch
  readonly getAccessToken?: TokenProvider
  readonly credentials: RequestCredentials
  readonly onUnauthorized?: () => void | Promise<void>
  readonly refreshAccessToken?: RefreshAccessToken

  readonly auth: AuthApi
  readonly users: UsersApi
  readonly admin: AdminApi
  readonly health: HealthApi
  readonly mfa: MFAApi
  readonly webauthn: WebAuthnApi
  readonly oidc: OIDCApi
  readonly federation: FederationApi
  readonly tenants: TenantsApi

  private readonly defaultHeaders: HTTPHeaders
  private readonly configuration: Configuration
  private refreshInFlight: Promise<string | null | undefined> | null = null

  constructor(options: GateForgeClientOptions) {
    this.baseUrl = options.baseUrl.replace(/\/+$/, '')
    this.fetch = options.fetch ?? globalThis.fetch.bind(globalThis)
    this.getAccessToken = options.getAccessToken
    this.credentials = options.credentials ?? 'include'
    this.defaultHeaders = headersInitToRecord(options.headers)
    this.onUnauthorized = options.onUnauthorized
    this.refreshAccessToken = options.refreshAccessToken

    const middleware: Middleware[] = [
      {
        post: async (context: ResponseContext): Promise<Response | void> => {
          if (!isUnauthorizedResponse(context.response)) {
            return
          }
          if (!this.refreshAccessToken) {
            await this.onUnauthorized?.()
            return
          }

          // Avoid refresh loops on the refresh call itself.
          if (context.init.headers) {
            const hdrs = new Headers(context.init.headers as HeadersInit)
            if (hdrs.get('X-GateForge-Retry') === '1') {
              await this.onUnauthorized?.()
              return
            }
          }

          const token = await this.runRefresh()
          if (!token && !(await this.getAccessToken?.())) {
            await this.onUnauthorized?.()
            return
          }

          const headers = new Headers(context.init.headers as HeadersInit)
          headers.set('X-GateForge-Retry', '1')
          const nextToken = token ?? (await this.getAccessToken?.())
          if (nextToken) {
            headers.set('Authorization', `Bearer ${nextToken}`)
          } else {
            headers.delete('Authorization')
          }

          return this.fetch(context.url, {
            ...context.init,
            headers,
          })
        },
      },
    ]

    this.configuration = new Configuration({
      basePath: this.baseUrl,
      fetchApi: this.fetch,
      credentials: this.credentials,
      headers: this.defaultHeaders,
      middleware,
      accessToken: this.getAccessToken
        ? async () => {
            const token = await this.getAccessToken?.()
            return token ?? ''
          }
        : undefined,
    })

    this.auth = new AuthApi(this.configuration)
    this.users = new UsersApi(this.configuration)
    this.admin = new AdminApi(this.configuration)
    this.health = new HealthApi(this.configuration)
    this.mfa = new MFAApi(this.configuration)
    this.webauthn = new WebAuthnApi(this.configuration)
    this.oidc = new OIDCApi(this.configuration)
    this.federation = new FederationApi(this.configuration)
    this.tenants = new TenantsApi(this.configuration)
  }

  /** Build an absolute URL for a path under the configured base. */
  url(path: string): string {
    if (/^https?:\/\//i.test(path)) {
      return path
    }
    const normalized = path.startsWith('/') ? path : `/${path}`
    return `${this.baseUrl}${normalized}`
  }

  /**
   * Low-level request helper. Parses GateForge Meta envelopes on error.
   * Applies the same bearer / credentials / 401-refresh behavior as generated APIs.
   */
  async request<T>(path: string, init: RequestInit = {}): Promise<Envelope<T>> {
    const doFetch = async (retried: boolean): Promise<Response> => {
      const headers = new Headers(this.defaultHeaders)
      new Headers(init.headers).forEach((value, key) => {
        headers.set(key, value)
      })

      if (this.getAccessToken) {
        const token = await this.getAccessToken()
        if (token) {
          headers.set('Authorization', `Bearer ${token}`)
        }
      }

      if (init.body && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }

      if (retried) {
        headers.set('X-GateForge-Retry', '1')
      }

      return this.fetch(this.url(path), {
        ...init,
        headers,
        credentials: init.credentials ?? this.credentials,
      })
    }

    let response = await doFetch(false)

    if (response.status === 401 && this.refreshAccessToken) {
      const token = await this.runRefresh()
      if (token || (await this.getAccessToken?.())) {
        response = await doFetch(true)
      } else {
        await this.onUnauthorized?.()
      }
    } else if (response.status === 401) {
      await this.onUnauthorized?.()
    }

    let body: Envelope<T> | undefined
    const contentType = response.headers.get('content-type') ?? ''
    if (contentType.includes('application/json')) {
      body = (await response.json()) as Envelope<T>
    }

    if (!response.ok) {
      const meta: Meta = body?.meta ?? {}
      throw new APIError(meta.message ?? (response.statusText || 'Request failed'), {
        status: response.status,
        errorCode: meta.error_code,
        meta,
        response,
      })
    }

    if (!body) {
      throw new APIError('Empty response body', {
        status: response.status,
        response,
      })
    }

    return body
  }

  private runRefresh(): Promise<string | null | undefined> {
    if (!this.refreshAccessToken) {
      return Promise.resolve(undefined)
    }
    this.refreshInFlight ??= this.refreshAccessToken().finally(() => {
      this.refreshInFlight = null
    })
    return this.refreshInFlight
  }
}
