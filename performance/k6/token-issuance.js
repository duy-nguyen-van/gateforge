import http from 'k6/http';
import { check, fail } from 'k6';
import { SharedArray } from 'k6/data';
import exec from 'k6/execution';
import { Counter, Rate, Trend } from 'k6/metrics';
import { baseURL, codesPath, jsonHeaders } from './lib/config.js';

/**
 * Primary marketing bench: isolated POST /token throughput with pre-seeded codes.
 * Publish sustained http_reqs from the highest arrival-rate stage with error rate < 1%.
 */
const codes = new SharedArray('auth-codes', function () {
  const path = codesPath();
  const data = JSON.parse(open(path));
  if (!Array.isArray(data) || data.length === 0) {
    throw new Error(`no codes in ${path} — run make -C performance seed-codes`);
  }
  return data;
});

const tokenDuration = new Trend('token_duration', true);
const tokenFail = new Rate('token_fail');
const codesExhausted = new Counter('codes_exhausted');

const holdRPS = Number(__ENV.PERF_HOLD_RPS || 12000);

export const options = {
  scenarios: {
    token_issuance: {
      executor: 'ramping-arrival-rate',
      startRate: 50,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 4000,
      stages: [
        { target: 50, duration: '10s' }, // warm
        { target: 500, duration: '20s' },
        { target: 2000, duration: '20s' },
        { target: 5000, duration: '20s' },
        { target: 10000, duration: '20s' },
        { target: holdRPS, duration: '60s' }, // sustained hold — marketing figure
        { target: Math.floor(holdRPS * 1.25), duration: '20s' }, // probe beyond
      ],
      exec: 'tokenExchange',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    token_fail: ['rate<0.01'],
    // Ramp probe — looser than token-hold; still watch p99 collapse
    token_duration: ['p(95)<500', 'p(99)<1500'],
  },
};

export function tokenExchange() {
  const offset = Number(__ENV.PERF_CODE_OFFSET || 0);
  const idx = offset + exec.scenario.iterationInTest;
  if (idx >= codes.length) {
    codesExhausted.add(1);
    fail('auth codes exhausted — re-seed with a larger -n');
  }
  const entry = codes[idx];
  const url = `${baseURL()}/token`;
  const payload = JSON.stringify({
    grant_type: 'authorization_code',
    code: entry.code,
    client_id: entry.client_id,
    code_verifier: entry.code_verifier,
    redirect_uri: entry.redirect_uri,
  });
  const res = http.post(url, payload, {
    headers: jsonHeaders,
    tags: { name: 'token' },
  });
  tokenDuration.add(res.timings.duration);
  const ok = check(res, {
    'token status 200': (r) => r.status === 200,
    'has access_token': (r) => {
      try {
        return !!r.json('access_token');
      } catch (_) {
        return false;
      }
    },
  });
  tokenFail.add(!ok);
}
