import { delay } from 'core/common/utils'
import {
  TGROperation,
  TGROperationCreateBody,
  TGROperationsListRO,
} from 'core/models/goods-receipt/operation'

type TGROperationsRequestFilters = {
  supNumber: string
}

export async function getOperationsList(
  _: TGROperationsRequestFilters // eslint-disable-line
): Promise<TGROperationsListRO> {
  await delay(3000)
  return {
    total: 1,
    list: [
      {
        agreement: null,
        agreementId: 0,
        buyAmountHV: '',
        buyPriceAmount: 3,
        buyPriceAmountHV: '',
        createDate: '',
        creator: '',
        debt: 3,
        debtHV: '',
        id: 0,
        number: '',
        positionsCount: 0,
        repaymentPeriod: 0,
        status: 1,
        sumQuantity: 3,
        sumSOQuantity: 3,
        supBuyAmountHV: '',
        supBuyPriceAmount: 2,
        supBuyPriceAmountHV: '',
        supNumber: '',
        supplier: '',
        supShipmentDate: '',
        worker: '',
      },
    ],
  }
}

export async function createOperation(
  _: TGROperationCreateBody // eslint-disable-line
): Promise<TGROperation> {
  await delay(3000)
  return {
    agreement: null,
    agreementId: 0,
    buyAmountHV: '',
    buyPriceAmount: 3,
    buyPriceAmountHV: '',
    createDate: '',
    creator: '',
    debt: 3,
    debtHV: '',
    id: 0,
    number: '',
    positionsCount: 0,
    repaymentPeriod: 0,
    status: 1,
    sumQuantity: 3,
    sumSOQuantity: 3,
    supBuyAmountHV: '',
    supBuyPriceAmount: 2,
    supBuyPriceAmountHV: '',
    supNumber: '',
    supplier: '',
    supShipmentDate: '',
    worker: '',
  }
}
