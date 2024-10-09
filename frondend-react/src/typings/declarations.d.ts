declare interface File {
  description?: string
  preview?: string
}

declare module '*.scss' {
  const content: { readonly [className: string]: string }
  export default content
}

declare module 'moment/locale/ru'
declare module 'moment/locale/uk'
declare module 'moment/locale/kk'
