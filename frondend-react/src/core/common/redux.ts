import { ActionCreatorWithoutPayload } from '@reduxjs/toolkit'
import Axios from 'axios'

import { AppDispatch } from '../store/types'

/**
 * Возвращает строку с глобальным префиксом
 * @param {string} suffix произвольное имя для добавления префикса
 * @returns {string}
 */
export const withGlobalPrefix = (suffix: string): string =>
  `${process.env.APP_NAME}/${suffix}`

/**
 * Добавляет суффикс типа "Start" для типа экшена
 * @param {string} actionType тип экшена
 * @returns {string}
 */
export const startSuffix = (actionType: string): string => `${actionType}_START`

/**
 * Добавляет суффикс типа "Resume" для типа экшена
 * @param {string} actionType тип экшена
 * @returns {string}
 */
export const resumeSuffix = (actionType: string): string =>
  `${actionType}_RESUME`

/**
 * Добавляет суффикс типа "Resume fail" для типа экшена
 * @param {string} actionType тип экшена
 * @returns {string}
 */
export const resumeFailSuffix = (actionType: string): string =>
  `${actionType}_RESUME_FAIL`

/**
 * Добавляет суффикс типа "Finish" для типа экшена
 * @param {string} actionType тип экшена
 * @returns {string}
 */
export const finishSuffix = (actionType: string): string =>
  `${actionType}_FINISH`

/**
 * Добавляет суффикс типа "Fail" для типа экшена
 * @param {string} actionType тип экшена
 * @returns {string}
 */
export const failSuffix = (actionType: string): string => `${actionType}_FAIL`

export type TActionData = { type: string; [key: string]: unknown }

type TActionDataOrType = string | ActionCreatorWithoutPayload | TActionData

type TRequest<T> = () => Promise<T>

type TActionCallbackParams<T> = {
  dispatch: AppDispatch
  error: unknown
  result: T | null
}

type TActionCallback<T> = (params: TActionCallbackParams<T>) => void

export type AsyncActionResponse<T> = (dispatch: AppDispatch) => Promise<T>

export const asyncAction = <T>(
  actionDataOrType: TActionDataOrType,
  request: TRequest<T>,
  cb?: TActionCallback<T>
): ((dispatch: AppDispatch) => Promise<T>) => {
  const action =
    typeof actionDataOrType === 'string'
      ? { type: actionDataOrType }
      : actionDataOrType

  return async (dispatch: AppDispatch) => {
    try {
      // Начало запроса
      dispatch({ ...action, type: startSuffix(action.type) } as TActionData)
      const result = await request()
      // Успешное выполнение
      dispatch({
        ...action,
        type: finishSuffix(action.type),
        result,
      } as TActionData & { result: T })
      if (cb) cb({ dispatch, error: null, result })
      return result
    } catch (e) {
      // Ошибка при выполнении
      if (!Axios.isCancel(e)) {
        dispatch({
          ...action,
          type: failSuffix(action.type),
          error: e,
        } as TActionData & { error: unknown })
        if (cb) cb({ dispatch, error: e, result: null })
      }
      throw e
    }
  }
}
