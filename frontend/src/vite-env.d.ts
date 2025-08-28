/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_DB_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
