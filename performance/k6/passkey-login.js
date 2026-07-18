import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { Rate, Trend } from 'k6/metrics';
import { baseURL, jsonHeaders, passkeysPath, signerURL } from './lib/config.js';

const users = new SharedArray('passkey-users', function () {
  const data = JSON.parse(open(passkeysPath()));
  if (!data.users || data.users.length === 0) {
    throw new Error('no passkey fixtures — run make -C performance seed-passkeys');
  }
  return data.users;
});

const passkeyE2E = new Trend('passkey_e2e', true);
const passkeyFail = new Rate('passkey_fail');

export const options = {
  scenarios: {
    steady: {
      executor: 'constant-vus',
      vus: Number(__ENV.PERF_PASSKEY_VUS || 10),
      duration: __ENV.PERF_PASSKEY_DURATION || '2m',
      exec: 'passkeyLogin',
      tags: { stage: 'steady' },
    },
  },
  thresholds: {
    passkey_fail: ['rate<0.01'],
    // p95 = landing claim; p99 catches rare Redis/DB stalls
    passkey_e2e: ['p(95)<80', 'p(99)<150'],
  },
};

export function passkeyLogin() {
  const user = users[(__VU - 1) % users.length];
  const t0 = Date.now();

  const startRes = http.post(
    `${baseURL()}/api/v1/webauthn/login/start`,
    JSON.stringify({ email: user.email }),
    { headers: jsonHeaders, tags: { name: 'passkey_start' } },
  );
  const started = check(startRes, {
    'start 200': (r) => r.status === 200,
  });
  if (!started) {
    passkeyFail.add(1);
    return;
  }
  const startBody = startRes.json();
  const options = startBody.data.options;
  const sessionToken = startBody.data.session_token;

  const assertRes = http.post(
    `${signerURL()}/assert`,
    JSON.stringify({ email: user.email, options }),
    { headers: jsonHeaders, tags: { name: 'signer_assert' } },
  );
  if (!check(assertRes, { 'signer 200': (r) => r.status === 200 })) {
    passkeyFail.add(1);
    return;
  }
  const credential = assertRes.json('credential');

  const finishRes = http.post(
    `${baseURL()}/api/v1/webauthn/login/finish`,
    JSON.stringify({
      email: user.email,
      session_token: sessionToken,
      credential,
    }),
    { headers: jsonHeaders, tags: { name: 'passkey_finish' } },
  );
  const finished = check(finishRes, {
    'finish 200': (r) => r.status === 200,
    'has access_token': (r) => {
      try {
        return !!r.json('data.access_token');
      } catch (_) {
        return false;
      }
    },
  });
  passkeyE2E.add(Date.now() - t0);
  passkeyFail.add(!finished);
  sleep(0.05);
}
