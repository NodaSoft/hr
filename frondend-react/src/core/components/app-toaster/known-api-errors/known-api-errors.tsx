import React from 'react'
import ReactDOM from 'react-dom'

import { Alert, Intent } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'

import { TApiError } from 'core/models/api'

import { KnownApiErrors, knownApiErrors } from './errors'

class KnownApiError {
  public isKnownError(apiError?: TApiError): apiError is TApiError {
    const errorCode = this.getErrorCode(apiError)
    if (typeof errorCode === 'number') {
      return !!this.knownApiErrors[errorCode]
    }
    return false
  }

  public errorHandler(apiError: TApiError) {
    const errorCode = this.getErrorCode(apiError)
    if (typeof errorCode === 'number') {
      const getErrorContent = this.knownApiErrors[errorCode]
      const errorContent = getErrorContent(apiError)
      return this.showAlert(errorContent)
    }
  }

  private knownApiErrors: KnownApiErrors

  constructor() {
    this.knownApiErrors = knownApiErrors
  }

  private getErrorCode(apiError?: TApiError) {
    return apiError?.response?.data?.code
  }

  private showAlert(content: React.ReactNode) {
    const containerElement = document.createElement('div')
    document.body.appendChild(containerElement)

    const handleConfirm = () => {
      containerElement.remove()
    }

    // eslint-disable-next-line react/no-deprecated
    ReactDOM.render(
      <Alert
        isOpen
        intent={Intent.DANGER}
        onConfirm={handleConfirm}
        icon={IconNames.WARNING_SIGN}
        portalContainer={containerElement}
        confirmButtonText={'common.actions.ok'}
      >
        {content}
      </Alert>,
      containerElement
    )
  }
}

export default new KnownApiError()
