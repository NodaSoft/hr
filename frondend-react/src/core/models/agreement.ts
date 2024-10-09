import * as Yup from 'yup'

import { ContractorTypes } from './contractor'

export type TAgreement = {
  agreementType: number
  balance: number
  balanceHV: string
  contractorDetailsId: number
  contractorId: number
  createDate: Date
  currency: string
  docTemplate: string
  id: number
  isActive: boolean
  isDelete: boolean
  isDefault: boolean
  name: string
  nameHV: string
  number: string
  paymentType: number
  personData: null
  relationType: number
  shopDetailsId: number
  currencyDesignation: string
  creditLimit: number
  creditLimitHV: string
  repaymentPeriod: number | null
}

export type TGetAgreementsListRO = {
  list: TAgreement[]
  total: number
}

export type MixedAgreement = {
  agreementId: number
  balance: number
  balanceHV: string
  contractorId: number
  contractorProfileId: number | null
  isDefault: boolean
  label: string
  relationType: ContractorTypes
  value: string
  isActive: boolean
  isDelete: boolean
  legalPersonId: number
  contractorLegalPersonINN: string
  shopLegalPersonINN: string
  contractorLegalPersonId: number
  currency?: { isoCode: string; designation: string }
  repaymentPeriod: number
}

export type MixedAgreementsResponse = {
  data: MixedAgreement[]
}

export const MixedAgreementShape = Yup.object().shape<MixedAgreement>({
  agreementId: Yup.number().required(),
  contractorId: Yup.number().required(),
  balance: Yup.number(),
  balanceHV: Yup.string(),
  contractorProfileId: Yup.number().nullable().default(null),
  isDefault: Yup.boolean().required(),
  label: Yup.string().required(),
  relationType: Yup.number().required(),
  value: Yup.string().required(),
  isActive: Yup.boolean(),
  isDelete: Yup.boolean(),
  legalPersonId: Yup.number(),
  contractorLegalPersonId: Yup.number(),
  contractorLegalPersonINN: Yup.string(),
  shopLegalPersonINN: Yup.string(),
  currency: Yup.object().shape({
    isoCode: Yup.string(),
    designation: Yup.string(),
  }),
  repaymentPeriod: Yup.number(),
})

export type TAgreementShape = Pick<
  TAgreement,
  'id' | 'name' | 'number' | 'createDate' | 'balance' | 'balanceHV' | 'nameHV'
>

export const AgreementShape = Yup.object().shape<TAgreementShape>({
  id: Yup.number(),
  name: Yup.string(),
  number: Yup.string(),
  nameHV: Yup.string(),
  createDate: Yup.string(),
  balance: Yup.number(),
  balanceHV: Yup.string(),
})
