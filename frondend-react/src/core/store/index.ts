import { combineReducers, configureStore } from '@reduxjs/toolkit'
import { connectRouter, routerMiddleware } from 'connected-react-router'
import { createBrowserHistory } from 'history'

import { BASENAME_URL } from 'core/common/config'

import {
  moduleKey as commonModuleKey,
  reducer as commonReducer,
} from './modules/common'
import {
  moduleKey as stuffModuleKey,
  reducer as stuffReducer,
} from './modules/stuff'

export const browserHistory = createBrowserHistory({
  basename: BASENAME_URL,
})

const getInitialStore = () => {
  const staticReducers = {
    router: connectRouter(browserHistory),
    [commonModuleKey]: commonReducer,
    [stuffModuleKey]: stuffReducer,
  }

  // @ts-ignore
  function dynamicAddReducer(asyncReducers) {
    return combineReducers({
      ...staticReducers,
      ...asyncReducers,
    })
  }

  const store = configureStore({
    reducer: { ...staticReducers },
    // @ts-ignore
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({ serializableCheck: false }).concat(
        routerMiddleware(browserHistory)
      ),
  })

  // Add a dictionary to keep track of the registered async reducers
  // @ts-ignore
  store.asyncReducers = {}

  // Create an inject reducer function
  // This function adds the async reducer, and creates a new combined reducer
  // @ts-ignore
  store.injectReducer = (key, asyncReducer) => {
    // @ts-ignore
    store.asyncReducers[key] = asyncReducer
    // @ts-ignore
    store.replaceReducer(dynamicAddReducer(store.asyncReducers))
  }

  // Return the modified store
  return store
}

export const store = getInitialStore()
