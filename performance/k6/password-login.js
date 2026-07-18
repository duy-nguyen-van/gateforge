import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { baseURL, jsonHeaders, perfUser } from './lib/config.js';

const loginDuration = new Trend('password_login_duration', true);
const loginFail = new Rate('password_login_fail');
const refreshDuration = new Trend('password_refresh_duration', true);

export const options = {
  scenarios: {
    login: {
      executor: 'constant-vus',
      vus: Number(__ENV.PERF_LOGIN_VUS || 25),
      duration: __ENV.PERF_LOGIN_DURATION || '2m',
      exec: 'passwordLogin',
    },
  },
  thresholds: {
    password_login_fail: ['rate<0.01'],
    // bcrypt-bound; looser than token path — still watch the tail
    password_login_duration: ['p(95)<300', 'p(99)<600'],
  },
};

export function passwordLogin() {
  const user = perfUser();
  const res = http.post(
    `${baseURL()}/api/v1/login`,
    JSON.stringify({ email: user.email, password: user.password }),
    { headers: jsonHeaders, tags: { name: 'password_login' } },
  );
  loginDuration.add(res.timings.duration);
  const ok = check(res, {
    'login 200': (r) => r.status === 200,
    'access_token': (r) => {
      try {
        return !!r.json('data.access_token');
      } catch (_) {
        return false;
      }
    },
  });
  loginFail.add(!ok);
  if (!ok) return;

  const refresh = res.json('data.refresh_token');
  if (refresh) {
    const rr = http.post(
      `${baseURL()}/api/v1/refresh`,
      JSON.stringify({ refresh_token: refresh }),
      { headers: jsonHeaders, tags: { name: 'password_refresh' } },
    );
    refreshDuration.add(rr.timings.duration);
    check(rr, { 'refresh 200': (r) => r.status === 200 });
  }
  sleep(0.05);
}
