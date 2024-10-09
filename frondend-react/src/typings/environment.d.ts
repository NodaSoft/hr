declare global {
  namespace NodeJS {
    interface ProcessEnv {
      APP_NAME: 'test-app'
      NODE_ENV: 'development' | 'production'
    }
  }
}
export {}
