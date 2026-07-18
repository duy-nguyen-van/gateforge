import http from 'k6/http';
import { check, fail } from 'k6';
import { SharedArray } from 'k6/data';
import exec from 'k6/execution';
import { Rate, Trend } from 'k6/metrics';
import { baseURL, codesPath, jsonHeaders } from './lib/config.js';

/** Constant-arrival hold to find sustainable token RPS with healthy latency. */
const codes = new SharedArray('auth-codes-hold', function () {
  return JSON.parse(open(codesPath()));
});

const tokenDuration = new Trend('token_duration', true);
const tokenFail = new Rate('token_fail');
const rate = Number(__ENV.PERF_HOLD_RPS || 500);
const duration = __ENV.PERF_HOLD_DURATION || '45s';

export const options = {
  scenarios: {
    hold: {
      executor: 'constant-arrival-rate',
      rate,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: Math.min(rate, 500),
      maxVUs: Math.max(rate * 2, 100),
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    token_fail: ['rate<0.01'],
    // p95 = marketing SLO; p99 = tail / capacity guardrail
    token_duration: ['p(95)<200', 'p(99)<500'],
    http_reqs: [`rate>${rate * 0.9}`],
  },
};

export default function () {
  const offset = Number(__ENV.PERF_CODE_OFFSET || 0);
  const idx = offset + exec.scenario.iterationInTest;
  if (idx >= codes.length) fail('codes exhausted');
  const entry = codes[idx];
  const res = http.post(
    `${baseURL()}/token`,
    JSON.stringify({
      grant_type: 'authorization_code',
      code: entry.code,
      client_id: entry.client_id,
      code_verifier: entry.code_verifier,
      redirect_uri: entry.redirect_uri,
    }),
    { headers: jsonHeaders, tags: { name: 'token' } },
  );
  tokenDuration.add(res.timings.duration);
  tokenFail.add(
    !check(res, {
      '200': (r) => r.status === 200,
      'access_token': (r) => {
        try {
          return !!r.json('access_token');
        } catch (_) {
          return false;
        }
      },
    }),
  );
}
