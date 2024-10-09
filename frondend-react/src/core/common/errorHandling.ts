import { Classes } from '@blueprintjs/core'

import { TApiError } from 'core/models/api'

export function errorTrace(
  error: TApiError,
  onlyClientMessage = false
): string {
  let message = ''
  if (
    error &&
    error.response &&
    typeof error.response === 'object' &&
    typeof error.response.data === 'object'
  ) {
    const { clientMessage, traceMessage, code } = error.response.data || {}
    if (typeof clientMessage === 'string' && clientMessage.length) {
      message += clientMessage
    }

    if (!onlyClientMessage) {
      if (typeof traceMessage === 'string' && traceMessage.length) {
        message += '<br />'
        message += `<span class=${Classes.MONOSPACE_TEXT}> ${traceMessage} </span>`
      }

      if (typeof code === 'string' || typeof code === 'number') {
        message += '<br />'
        message += `CODE: ${code}`
      }
    }
  }
  return message
}
