import { createSlice } from '@reduxjs/toolkit'

import { withGlobalPrefix } from 'core/common/redux'

import { ModuleSelector } from '../../interfaces/utils'
import { TRootState } from '../types'
import { TCommonState } from '../types/modules/common'

export const moduleKey = 'common'
const reducerPrefix = withGlobalPrefix(moduleKey)

function moduleSelector<T extends unknown[], R>(
  selector: ModuleSelector<TCommonState, T, R>
) {
  return (globalState: TRootState, ...args: T) =>
    selector(globalState[moduleKey], ...args)
}

export const selectCurrentEmployeeId = moduleSelector(
  (state) => state.currentEmployeeId || null
)

const initialState: TCommonState = {
  currentEmployeeId: 25131850,
}

const commonSlice = createSlice({
  name: reducerPrefix,
  initialState,
  reducers: {},
})

export const { reducer } = commonSlice
