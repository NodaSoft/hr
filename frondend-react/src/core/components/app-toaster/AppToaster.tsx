import React from 'react'

import { Intent, IToastProps, Position, Toaster } from '@blueprintjs/core'

import DOMPurify from 'dompurify'

import { errorTrace } from 'core/common/errorHandling'
import { TApiError } from 'core/models/api'

import KnownApiError from './known-api-errors/known-api-errors'

interface IToastError extends Omit<IToastProps, 'message'> {
  error?: TApiError
  message?: string
}

const AppToaster = () => {
  const Toast = Toaster.create({
    position: Position.TOP,
    // className: Styles.AppToaster
  })

  // Все ссылки открываем в новой вкладке
  DOMPurify.addHook('afterSanitizeAttributes', (node) => {
    if ('target' in node) {
      node.setAttribute('target', '_blank')
      node.setAttribute('rel', 'noopener noreferrer')
    }
  })

  const renderMessage = (message: IToastProps['message']) => {
    if (typeof message === 'string') {
      return (
        <div
          dangerouslySetInnerHTML={{
            __html: DOMPurify.sanitize(message),
          }}
        />
      )
    }

    return message
  }

  return {
    show: ({ message, ...args }: IToastProps) => {
      if (!message) return
      return Toast.show({ message: renderMessage(message), ...args })
    },
    success: ({ message = 'common.message.success', ...args }: IToastProps) => {
      return Toast.show({
        message: renderMessage(message),
        intent: Intent.SUCCESS,
        ...args,
      })
    },
    warn: ({ message, ...args }: IToastProps) => {
      if (!message) return
      return Toast.show({
        message: renderMessage(message),
        intent: Intent.WARNING,
        ...args,
      })
    },
    error: ({
      error,
      message = 'critical.errors.default.title',
      ...args
    }: IToastError) => {
      if (KnownApiError.isKnownError(error)) {
        return KnownApiError.errorHandler(error)
      }

      return Toast.show({
        message: renderMessage((error && errorTrace(error)) || message),
        intent: Intent.DANGER,
        ...args,
      })
    },
  }
}

export default AppToaster()
