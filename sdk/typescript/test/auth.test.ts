import assert from 'node:assert/strict'
import { describe, it } from 'node:test'

import {
  buildAuthorizeUrl,
  createPKCE,
  generateNonce,
  generateState,
  parseCallbackParams,
} from '../src/auth.ts'

describe('auth helpers', () => {
  it('createPKCE returns S256 verifier/challenge', async () => {
    const pkce = await createPKCE()
    assert.equal(pkce.method, 'S256')
    assert.ok(pkce.verifier.length > 20)
    assert.ok(pkce.challenge.length > 20)
    assert.notEqual(pkce.verifier, pkce.challenge)
  })

  it('buildAuthorizeUrl appends /authorize and query params', () => {
    const url = buildAuthorizeUrl('https://iam.example.com', {
      clientId: 'app',
      redirectUri: 'https://app.example.com/cb',
      scope: 'openid profile',
      state: 's1',
      nonce: 'n1',
      codeChallenge: 'challenge',
    })
    const parsed = new URL(url)
    assert.equal(parsed.origin, 'https://iam.example.com')
    assert.equal(parsed.pathname, '/authorize')
    assert.equal(parsed.searchParams.get('client_id'), 'app')
    assert.equal(parsed.searchParams.get('redirect_uri'), 'https://app.example.com/cb')
    assert.equal(parsed.searchParams.get('response_type'), 'code')
    assert.equal(parsed.searchParams.get('scope'), 'openid profile')
    assert.equal(parsed.searchParams.get('state'), 's1')
    assert.equal(parsed.searchParams.get('nonce'), 'n1')
    assert.equal(parsed.searchParams.get('code_challenge'), 'challenge')
    assert.equal(parsed.searchParams.get('code_challenge_method'), 'S256')
  })

  it('generateState and generateNonce are url-safe', () => {
    const state = generateState()
    const nonce = generateNonce()
    assert.match(state, /^[A-Za-z0-9_-]+$/)
    assert.match(nonce, /^[A-Za-z0-9_-]+$/)
    assert.notEqual(state, nonce)
  })

  it('parseCallbackParams reads code/state and errors', () => {
    const ok = parseCallbackParams(
      'https://app.example.com/cb?code=abc&state=xyz',
    )
    assert.deepEqual(ok, {
      code: 'abc',
      state: 'xyz',
      error: undefined,
      error_description: undefined,
    })

    const err = parseCallbackParams(
      new URL(
        'https://app.example.com/cb?error=access_denied&error_description=nope&state=s',
      ),
    )
    assert.equal(err.error, 'access_denied')
    assert.equal(err.error_description, 'nope')
    assert.equal(err.state, 's')
    assert.equal(err.code, undefined)
  })
})
