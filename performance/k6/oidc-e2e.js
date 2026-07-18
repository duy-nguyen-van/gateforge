import http from 'k6/http';
import { check } from 'k6';
import { sha256 } from 'k6/crypto';
import { Rate, Trend, Counter } from 'k6/metrics';
import { baseURL, clientID, jsonHeaders, perfUser, redirectURI } from './lib/config.js';

/**
 * Secondary: complete OIDC login journey.
 * Metrics: complete logins/s (oidc_complete_logins) and oidc_e2e p95/p99 — not token-issuance RPS.
 */
const oidcE2E = new Trend('oidc_e2e', true);
const oidcFail = new Rate('oidc_fail');
const oidcLogins = new Counter('oidc_complete_logins');

function randomString(len) {
  const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~';
  let out = '';
  for (let i = 0; i < len; i++) {
    out += alphabet[Math.floor(Math.random() * alphabet.length)];
  }
  return out;
}

function pkcePair() {
  const verifier = randomString(64);
  const challenge = sha256(verifier, 'base64rawurl');
  return { verifier, challenge };
}

export const options = {
  scenarios: {
    e2e: {
      executor: 'constant-vus',
      vus: Number(__ENV.PERF_OIDC_VUS || 20),
      duration: __ENV.PERF_OIDC_DURATION || '2m',
      exec: 'oidcFlow',
    },
  },
  thresholds: {
    oidc_fail: ['rate<0.02'],
    oidc_e2e: ['p(95)<1500', 'p(99)<2500'],
  },
};

export function oidcFlow() {
  const user = perfUser();
  const { verifier, challenge } = pkcePair();
  const state = randomString(16);
  const redirect = redirectURI();
  const cid = clientID();
  const t0 = Date.now();

  const jar = http.cookieJar();
  const loginRes = http.post(
    `${baseURL()}/api/v1/login`,
    JSON.stringify({ email: user.email, password: user.password }),
    { headers: jsonHeaders, jar, tags: { name: 'login' } },
  );
  if (!check(loginRes, { 'login 200': (r) => r.status === 200 })) {
    oidcFail.add(1);
    return;
  }

  const authURL =
    `${baseURL()}/authorize?response_type=code` +
    `&client_id=${encodeURIComponent(cid)}` +
    `&redirect_uri=${encodeURIComponent(redirect)}` +
    `&scope=${encodeURIComponent('openid email profile')}` +
    `&state=${encodeURIComponent(state)}` +
    `&code_challenge=${encodeURIComponent(challenge)}` +
    `&code_challenge_method=S256`;

  const authRes = http.get(authURL, {
    redirects: 0,
    jar,
    tags: { name: 'authorize' },
  });
  const location = authRes.headers.Location || authRes.headers.location || '';
  const codeMatch = /[?&]code=([^&]+)/.exec(location);
  if (
    !check(authRes, {
      'authorize redirect': (r) => r.status === 302 || r.status === 303,
      'has code': () => !!codeMatch,
    })
  ) {
    oidcFail.add(1);
    return;
  }
  const code = decodeURIComponent(codeMatch[1]);

  const tokenRes = http.post(
    `${baseURL()}/token`,
    JSON.stringify({
      grant_type: 'authorization_code',
      code,
      client_id: cid,
      code_verifier: verifier,
      redirect_uri: redirect,
    }),
    { headers: jsonHeaders, tags: { name: 'token' } },
  );
  if (!check(tokenRes, { 'token 200': (r) => r.status === 200 })) {
    oidcFail.add(1);
    return;
  }
  let access = '';
  try {
    access = tokenRes.json('access_token');
  } catch (_) {
    oidcFail.add(1);
    return;
  }

  const uiRes = http.get(`${baseURL()}/userinfo`, {
    headers: { Authorization: `Bearer ${access}` },
    tags: { name: 'userinfo' },
  });
  const ok = check(uiRes, { 'userinfo 200': (r) => r.status === 200 });
  oidcE2E.add(Date.now() - t0);
  oidcFail.add(!ok);
  if (ok) oidcLogins.add(1);
}
