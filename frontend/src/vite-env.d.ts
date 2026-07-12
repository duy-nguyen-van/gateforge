/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_API_PROXY_TARGET: string
  readonly VITE_DEFAULT_TENANT_ID: string
  readonly VITE_AGENTATION_ENDPOINT?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
