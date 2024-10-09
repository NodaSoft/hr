/* eslint-disable */
import { replace } from 'connected-react-router'
import { difference as _difference, reduce as _reduce } from 'lodash-es'
import { createSelector } from 'reselect'

import grAPI from 'core/api/goods-receipts'
import {
  asyncAction,
  AsyncActionResponse,
  failSuffix,
  finishSuffix,
  startSuffix,
  TActionData,
  withGlobalPrefix,
} from 'core/common/redux'
import { GR_OPERATION, toPath } from 'core/common/routes'
import { ModuleSelector } from 'core/interfaces/utils'
import {
  TGROperation,
  TGROperationCreateBody,
} from 'core/models/goods-receipt/operation'
import { TRootState } from 'core/store/types'
import { TStateOperationsGR } from 'core/store/types/good-receipts/operations'

export const moduleKey = 'operationsGR'

const reducerPrefix = withGlobalPrefix(moduleKey)

const CREATE_OPERATION = `${reducerPrefix}/CREATE_OPERATION`

export function createOperation(
  operation: TGROperationCreateBody
): AsyncActionResponse<TGROperation> {
  return asyncAction(
    CREATE_OPERATION,
    () => grAPI.createOperation(operation),
    ({ dispatch, result, error }) => {
      if (!error && result && result.id) {
        dispatch(replace(toPath(GR_OPERATION, { opId: result.id })))
      }
    }
  )
}

function moduleSelector<T extends unknown[], R>(
  selector: ModuleSelector<TStateOperationsGR, T, R>
) {
  return (globalState: TRootState, ...args: T) =>
    selector(globalState[moduleKey], ...args)
}

export const selectList = moduleSelector((state) => state.list || {})

export const selectStatusRules = moduleSelector((state) => state.statusRules)

export const operationByIdSelector = createSelector(
  [selectList, (_, opId: number) => opId],
  (list, opId) => list[opId] || null
)

export const selectFinalRule = createSelector(
  selectStatusRules,
  (statusRules) => statusRules.find((rule) => rule.final)
)

export const isOperationReadOnlySelector = createSelector(
  [operationByIdSelector, selectFinalRule],
  (operation, finalRule) => {
    return operation && finalRule ? operation.status === finalRule.id : false
  }
)

export const initialState = {
  list: {},
  byIds: [],
  total: 0,
  creating: false,
  statusRules: [],
}

export const reducer = (
  state = initialState,
  action: TActionData
): TStateOperationsGR => {
  switch (action.type) {
    case startSuffix(CREATE_OPERATION):
      return {
        ...state,
        creating: true,
      }
    case finishSuffix(CREATE_OPERATION): {
      const operation = action.result as TGROperation
      return {
        ...state,
        creating: false,
        byIds: [operation.id, ...state.byIds],
        list: {
          ...state.list,
          [operation.id]: operation,
        },
        total: state.total + 1,
      }
    }
    case failSuffix(CREATE_OPERATION):
      return {
        ...state,
        creating: false,
      }
    default:
      return state
  }
}
