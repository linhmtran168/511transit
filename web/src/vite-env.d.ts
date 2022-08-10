/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ImportMetaEnv {
  readonly VITE_LOCAL_WS_HOST: string
  readonly VITE_TRIP_UPDATE_INTERVAL: number
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}