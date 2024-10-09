import { TCommonStatusRules } from 'core/common/status'
import { TGROperation } from 'core/models/goods-receipt/operation'

export type TStateOperationsGR = {
  list: Record<number, TGROperation>
  byIds: number[]
  total: number
  creating: boolean
  statusRules: TCommonStatusRules
}
