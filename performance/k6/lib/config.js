// Shared k6 config for GateForge IAM marketing benches.
export function baseURL() {
  return __ENV.PERF_BASE_URL || 'http://127.0.0.1:3000';
}

export function signerURL() {
  return __ENV.PERF_SIGNER_URL || 'http://127.0.0.1:9091';
}

export function codesPath() {
  return __ENV.PERF_CODES_FILE || '../.data/codes.json';
}

export function passkeysPath() {
  return __ENV.PERF_PASSKEYS_FILE || '../.data/passkeys.json';
}

export function perfUser() {
  return {
    email: __ENV.PERF_USER_EMAIL || 'perf-token@example.com',
    password: __ENV.PERF_USER_PASSWORD || 'perf-token-passphrase-long!!',
  };
}

export function clientID() {
  return __ENV.PERF_CLIENT_ID || 'oidc-dev';
}

export function redirectURI() {
  return __ENV.PERF_REDIRECT_URI || 'http://localhost:3000/callback';
}

export const jsonHeaders = {
  'Content-Type': 'application/json',
  Accept: 'application/json',
};
