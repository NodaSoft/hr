export type TStatusRule = number | string | number[] | null

export type TCommonStatusRule = {
  id: number | string
  nextTo?: TStatusRule
  backTo?: TStatusRule
  cancelable?: boolean
  cancel?: boolean
  final?: boolean
}

// Общая модель статусов. Идентичная во всех модулях
export type TCommonStatusRules = TCommonStatusRule[]
