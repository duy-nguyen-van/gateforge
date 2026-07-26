import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import { GateForgeClient, APIError } from '../src/index.ts'

describe('GateForgeClient', () => {
  it('getHealth via generated API', async () => {
    const calls = []
    const fetchMock = async (input, init = {}) => {
      calls.push({ url: String(input), init })
      return new Response(
        JSON.stringify({
          meta: { message: 'ok', code: 200 },
          data: {
            status: 'healthy',
            service: 'gateforge-iam',
            version: '0.1.0',
            timestamp: '2026-01-01T00:00:00Z',
          },
        }),
        {
          status: 200,
          headers: { 'content-type': 'application/json' },
        },
      )
    }

    const client = new GateForgeClient({
      baseUrl: 'https://iam.example.com/',
      fetch: fetchMock,
      credentials: 'include',
    })

    const envelope = await client.health.getHealth()
    assert.equal(envelope.data?.status, 'healthy')
    assert.equal(calls.length, 1)
    assert.equal(calls[0].url, 'https://iam.example.com/api/v1/')
    assert.equal(calls[0].init.credentials, 'include')
  })

  it('sends Authorization bearer from getAccessToken', async () => {
    const calls = []
    const fetchMock = async (input, init = {}) => {
      calls.push({ url: String(input), headers: new Headers(init.headers) })
      return new Response(
        JSON.stringify({
          meta: { message: 'ok', code: 200 },
          data: {
            id: 'u1',
            email: 'a@b.co',
            first_name: 'A',
            last_name: 'B',
            email_verified: true,
            is_platform_admin: false,
            mfa_enabled: false,
            created_at: '2026-01-01T00:00:00Z',
            updated_at: '2026-01-01T00:00:00Z',
          },
        }),
        {
          status: 200,
          headers: { 'content-type': 'application/json' },
        },
      )
    }

    const client = new GateForgeClient({
      baseUrl: 'https://iam.example.com',
      fetch: fetchMock,
      getAccessToken: () => 'test-token',
    })

    await client.users.getMe()
    assert.equal(calls[0].headers.get('Authorization'), 'Bearer test-token')
  })

  it('request() throws APIError from error envelope', async () => {
    const fetchMock = async () =>
      new Response(
        JSON.stringify({
          meta: { message: 'Nope', error_code: 'forbidden', code: 403 },
          data: null,
        }),
        {
          status: 403,
          statusText: 'Forbidden',
          headers: { 'content-type': 'application/json' },
        },
      )

    const client = new GateForgeClient({
      baseUrl: 'https://iam.example.com',
      fetch: fetchMock,
    })

    await assert.rejects(
      () => client.request('/api/v1/me'),
      (err) => {
        assert.ok(err instanceof APIError)
        assert.equal(err.status, 403)
        assert.equal(err.errorCode, 'forbidden')
        assert.equal(err.message, 'Nope')
        return true
      },
    )
  })
})
