import React from 'react'
import { Provider as ReduxProvider } from 'react-redux'
import { Redirect, Route, Switch } from 'react-router-dom'

import { ConnectedRouter as RouterProvider } from 'connected-react-router'
import { History } from 'history'

import { GR_INDEX } from 'core/common/routes'
import { NotReadyState } from 'core/components'
import { store as AppStore } from 'core/store'

const GoodsReceipts = React.lazy(() => import('goods-receipts/Root'))

type Props = {
  store: typeof AppStore
  history: History<unknown>
}

const App = ({ store, history }: Props) => {
  return (
    <ReduxProvider store={store}>
      <RouterProvider history={history}>
        <React.Suspense fallback={<NotReadyState />}>
          <Switch>
            <Route path={GR_INDEX} component={GoodsReceipts} />
            <Redirect to={GR_INDEX} />
          </Switch>
        </React.Suspense>
      </RouterProvider>
    </ReduxProvider>
  )
}

export default App
