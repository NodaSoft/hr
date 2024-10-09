import React from 'react'
import ReactDOM from 'react-dom'

import { FocusStyleManager } from '@blueprintjs/core'

import 'core-js/features/array/flat-map'
import 'intersection-observer'

import { setLocale as setDateLocale } from 'core/common/date'
import logger from 'core/common/logger'
import { Locale } from 'core/models/localization'
import { browserHistory, store } from 'core/store'

import App from './App'
import './index.scss'

const bootstrap = () => {
  FocusStyleManager.onlyShowFocusOnTabs()
  setDateLocale(Locale.ru_RU)

  const run = () => {
    try {
      const entryPointDOM = document.getElementById('root')

      // eslint-disable-next-line react/no-deprecated
      ReactDOM.render(
        React.createElement(App, {
          store,
          history: browserHistory,
        }),
        entryPointDOM
      )
    } catch (e) {
      logger.error(e)
    }
  }

  run()
}

bootstrap()
