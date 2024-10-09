import Axios, { AxiosError } from 'axios'
import { FORM_ERROR } from 'final-form'
import { isObject } from 'lodash-es'

import logger from 'core/common/logger'
import { AppToaster } from 'core/components'
import { TApiError } from 'core/models/api'

export const POST_REQUEST = 'post'

// Функция проверяющая отменен ли предыдуший запрос или нет. Помогает не выкидывать много тостов при ошибке
export function isNotCancelRequest(e: unknown): boolean {
  return !Axios.isCancel(e)
}

export function isApiError(value: unknown): value is TApiError {
  return isObject(value) && 'response' in value && isNotCancelRequest(value)
}

export type TFormError = {
  'FINAL_FORM/form-error': AxiosError<unknown>
}

// Обработчик ошибок при отправке формы
export function formSubmitErrorHandler(e: TApiError): TFormError | undefined {
  if (isNotCancelRequest(e)) {
    logger.error(e)
    return { [FORM_ERROR]: e }
  }
}

type TApiErrorHandlerParams = {
  error: TApiError
  clientMessage?: string
}

// Общий обработчик ошибок API // TODO: подумать, сейчас плохая реализация
export function apiErrorHandler({
  error,
  clientMessage,
}: TApiErrorHandlerParams): void {
  if (isNotCancelRequest(error)) {
    AppToaster.error({
      message: clientMessage,
      error,
    })
    logger.error(error)
  }
}
