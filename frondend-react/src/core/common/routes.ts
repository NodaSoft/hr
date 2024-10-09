import pathToRegexp, { PathFunction } from 'path-to-regexp'

import { API_BASENAME_URL, BASENAME_URL } from './config'
import logger from './logger'

export const GR_INDEX = '/goodsReceipts'
export const GR_OPERATION = '/goodsReceipts/operations/:opId(\\d+)'
export const GR_OPERATION_NEW = '/goodsReceipts/operations/new'
export const GR_POSITION_NEW =
  '/goodsReceipts/operations/:opId(\\d+)/positions/new'
export const GR_POSITION_NEW_FROM_ORDERS =
  '/goodsReceipts/operations/:opId(\\d+)/positions/new/fromOrders'
export const GR_POSITIONS_FROM_NEW_ORDERS =
  '/goodsReceipts/operations/:opId(\\d+)/positions/new/fromNewOrders'

const memoizedPaths: Record<string, PathFunction> = {}

type TRouteParams = Record<string, string | number | boolean | null | undefined>

/**
 * генерация параметризованых маршрутов
 * @param {string} path путь формата path-to-regexp
 * @param {{withBase, withAPIBase ...routeParams}} options withBase || withAPIBase флаг генерации пути с BASENAME || withAPIBase, routeParams параметры маршрута path-to-regexp
 */
export function toPath(
  path: string,
  { withBase, withAPIBase, ...routeParams }: TRouteParams = {}
): string {
  try {
    if (!memoizedPaths[path]) memoizedPaths[path] = pathToRegexp.compile(path)
    // проверки на переданный флаг
    if (withBase) return `${BASENAME_URL}${memoizedPaths[path](routeParams)}`
    if (withAPIBase)
      return `${API_BASENAME_URL}${memoizedPaths[path](routeParams)}`
    return memoizedPaths[path](routeParams)
  } catch (e) {
    logger.error(e)
    return ''
  }
}
