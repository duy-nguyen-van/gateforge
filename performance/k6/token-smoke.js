import http from 'k6/http';
import { check, fail } from 'k6';
import { SharedArray } from 'k6/data';
import exec from 'k6/execution';
import { baseURL, codesPath, jsonHeaders } from './lib/config.js';

const codes = new SharedArray('auth-codes-smoke', function () {
  return JSON.parse(open(codesPath()));
});

export const options = {
  scenarios: {
    smoke: {
      executor: 'constant-arrival-rate',
      rate: 50,
      timeUnit: '1s',
      duration: '15s',
      preAllocatedVUs: 20,
      maxVUs: 50,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
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
    { headers: jsonHeaders },
  );
  check(res, { '200': (r) => r.status === 200 });
}
