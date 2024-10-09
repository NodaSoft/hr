import { AxiosError } from 'axios'

export type TApiError = AxiosError<{
  clientMessage: string
  code?: number | string
  httpCode: number
  info: string | null
  traceMessage: string
  type: number
}>
