import {
  GR_OPERATION,
  GR_POSITION_NEW,
  GR_POSITION_NEW_FROM_ORDERS,
  GR_POSITIONS_FROM_NEW_ORDERS,
  toPath,
} from 'core/common/routes'

export const GR_CREATE_NEW_OP_TYPE_LOCAL_STORAGE_KEY =
  'goodsReceipts.createNewOpOption'

export enum CreateNewOpTypes {
  OPEN_OPERATION = 'open',
  ADD_POSITION_BY_ARTICLE = 'byArticle',
  ADD_POSITION_BY_ORDER = 'byOrder',
  ADD_POSITION_BY_NEW_ORDER = 'byNewOrder',
}

export const DEFAULT_CREATE_NEW_TYPE = CreateNewOpTypes.ADD_POSITION_BY_ARTICLE

export const CREATE_NEW_OP_OPTIONS = Object.values(CreateNewOpTypes).map(
  (option) => {
    const names = {
      [CreateNewOpTypes.ADD_POSITION_BY_ARTICLE]:
        'Сохранить и добавить по артикулу',
      [CreateNewOpTypes.ADD_POSITION_BY_NEW_ORDER]:
        'Сохранить и добавить из заказа 2.0',
      [CreateNewOpTypes.ADD_POSITION_BY_ORDER]:
        'Сохранить и добавить из заказа',
      [CreateNewOpTypes.OPEN_OPERATION]: 'Сохранить и открыть просмотр',
    }

    return {
      id: option,
      name: names[option],
    }
  }
)

export const DEFAULT_CREATE_NEW_OPTION = CREATE_NEW_OP_OPTIONS.filter(
  (option) => option.id === DEFAULT_CREATE_NEW_TYPE
)[0]

export const getLinkAfterSuccessCreation = (
  type: CreateNewOpTypes,
  id: string
): string => {
  const successLink = (path: string) =>
    toPath(path, {
      opId: id,
    })

  switch (type) {
    case CreateNewOpTypes.OPEN_OPERATION:
      return successLink(GR_OPERATION)
    case CreateNewOpTypes.ADD_POSITION_BY_ARTICLE:
      return successLink(GR_POSITION_NEW)
    case CreateNewOpTypes.ADD_POSITION_BY_ORDER:
      return successLink(GR_POSITION_NEW_FROM_ORDERS)
    case CreateNewOpTypes.ADD_POSITION_BY_NEW_ORDER:
      return successLink(GR_POSITIONS_FROM_NEW_ORDERS)
    default:
      return successLink(GR_POSITION_NEW)
  }
}
