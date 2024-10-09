import { createSlice } from '@reduxjs/toolkit'

import { withGlobalPrefix } from 'core/common/redux'
import { ModuleSelector } from 'core/interfaces/utils'
import { TStuffState } from 'core/store/types/modules/stuff'

import { TRootState } from '../types'

export const moduleKey = 'stuff'
const reducerPrefix = withGlobalPrefix(moduleKey)

function moduleSelector<T extends unknown[], R>(
  selector: ModuleSelector<TStuffState, T, R>
) {
  return (globalState: TRootState, ...args: T) =>
    selector(globalState[moduleKey], ...args)
}

export const selectEmployeeById = moduleSelector((state, id?: number | null) =>
  state.list.find((item) => item.id === id)
)

export const selectEmployees = moduleSelector((state) => state.list)

const initialState: TStuffState = {
  list: [
    {
      id: 25131850,
      name: 'Большой Босс',
      firstName: 'Большой',
      lastName: 'Босс',
      isDelete: false,
      email: 'asdfasdf@gmail.com',
      mobile: '',
      phone: '+738545784',
      photoImageName: '',
      sip: '1123123',
      type: {
        id: 26597,
        name: 'Сотрудник',
        comment: '',
      },
    },
    {
      id: 25131123,
      name: 'Босс',
      firstName: 'Босс',
      lastName: 'Босс',
      isDelete: false,
      email: 'boss@gmail.com',
      mobile: '',
      phone: '+738123545784',
      photoImageName: '',
      sip: '1111',
      type: {
        id: 26597,
        name: 'Сотрудник',
        comment: '',
      },
    },
  ],
}

export const { reducer } = createSlice({
  name: reducerPrefix,
  initialState,
  reducers: {},
})
