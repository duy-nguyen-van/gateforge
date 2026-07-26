import { GateForgeClient } from '@gateforge/sdk'

import { getAccessToken, refreshTokens } from '@/auth/token-store'
import { getApiBaseUrl } from '@/lib/utils'

/**
 * Shared GateForge SDK client for the SPA.
 * Dev uses Vite proxy with a relative base (empty string); production may set VITE_API_BASE_URL.
 */
export const gateforge = new GateForgeClient({
  baseUrl: getApiBaseUrl(),
  credentials: 'include',
  getAccessToken: () => getAccessToken(),
  refreshAccessToken: async () => {
    await refreshTokens()
    return getAccessToken()
  },
})
