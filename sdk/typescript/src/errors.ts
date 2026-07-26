/** Meta matches the GateForge IAM JSON envelope metadata. */
export type Meta = {
  error_code?: string
  message?: string
  code?: number
  page?: number
  page_size?: number
  total?: number
}

/** Standard GateForge IAM response wrapper. */
export type Envelope<T> = {
  meta: Meta
  data: T
}

/** Base error for GateForge SDK failures. */
export class GateForgeError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'GateForgeError'
  }
}

export type APIErrorOptions = {
  status: number
  errorCode?: string
  meta?: Meta
  response?: Response
}

/**
 * Thrown when an HTTP response indicates failure.
 * Surfaces `meta.error_code` / `meta.message` from the envelope when present.
 */
export class APIError extends GateForgeError {
  readonly status: number
  readonly errorCode?: string
  readonly meta: Meta
  readonly response?: Response

  constructor(message: string, options: APIErrorOptions) {
    super(message)
    this.name = 'APIError'
    this.status = options.status
    this.errorCode = options.errorCode
    this.meta = options.meta ?? {}
    this.response = options.response
  }
}
