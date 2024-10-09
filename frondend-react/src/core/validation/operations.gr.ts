import * as Yup from 'yup'

import { getValidateFormErrors } from 'core/common/utils'
import { MixedAgreement, MixedAgreementShape } from 'core/models/agreement'
import {
  TGROperation,
  TGROperationUpdateBody,
} from 'core/models/goods-receipt/operation'
import { PersonSchema, TPerson } from 'core/models/person'

import {
  CREATE_DATE_VALIDATION_EMPTY_TITLE,
  CREATOR_VALIDATION_EMPTY_TITLE,
  NUMBER_VALIDATION_EMPTY_TITLE,
  SUPPLIER_NUMBER_VALIDATION_MAX_TITLE,
  SUPPLIER_VALIDATION_EMPTY_TITLE,
  WORKER_VALIDATION_EMPTY_TITLE,
} from './constants'

type TGROperationSchema = {
  id: number
  worker: TPerson | null
  creator: TPerson | null
  manualNumber?: boolean
  number?: string
  createDate: string | null
  mixedAgreement: MixedAgreement
  supNumber?: string | null
  supShipmentDate?: string | null
}

const supNumberMaxLength = 32

export const GROperationSchema = Yup.object().shape({
  id: Yup.number(),
  worker: PersonSchema.nullable()
    .default(null)
    .required(WORKER_VALIDATION_EMPTY_TITLE),
  creator: PersonSchema.nullable()
    .default(null)
    .required(CREATOR_VALIDATION_EMPTY_TITLE),
  manualNumber: Yup.boolean().default(true),
  number: Yup.string().when('manualNumber', {
    is: true,
    then: Yup.string().required(NUMBER_VALIDATION_EMPTY_TITLE),
    otherwise: Yup.string().strip(true),
  }),
  mixedAgreement: MixedAgreementShape.nullable()
    .default(null)
    .required(SUPPLIER_VALIDATION_EMPTY_TITLE),
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  createDate: Yup.string()
    .required(CREATE_DATE_VALIDATION_EMPTY_TITLE)
    .nullable()
    .default(null),
  supNumber: Yup.string()
    .nullable()
    .max(supNumberMaxLength, SUPPLIER_NUMBER_VALIDATION_MAX_TITLE),
  supShipmentDate: Yup.string().nullable().default(null),
})

export const GROperationAPISchema = GROperationSchema.from(
  'worker.id',
  'workerId'
)
  .from('creator.id', 'creatorId')
  .from('mixedAgreement.contractorId', 'supplierId')
  .from('mixedAgreement.agreementId', 'agreementId')
  .shape({
    mixedAgreement: Yup.object().strip(true),
    manualNumber: Yup.boolean().strip(true),
    worker: Yup.object().strip(true),
    creator: Yup.object().strip(true),
    number: Yup.string().when('$values.manualNumber', {
      is: true,
      then: Yup.string(),
      otherwise: Yup.string().strip(true),
    }),
    createDate: Yup.string(),
    agreementId: Yup.number(),
    workerId: Yup.number(),
    creatorId: Yup.number(),
    supplierId: Yup.number(),
    id: Yup.number(),
    supNumber: Yup.string().nullable(),
    supShipmentDate: Yup.string().nullable(),
  })

export function getInitialValues(
  operation: Partial<TGROperation> = {},
  context: { [key: string]: boolean } = {}
): Partial<TGROperationSchema> {
  const formValues = GROperationSchema.cast(
    {
      ...operation,
    },
    {
      stripUnknown: true,
      context: { values: operation, ...context },
    }
  )
  // @ts-ignore
  return {
    ...formValues,
    manualNumber: formValues?.id !== undefined && !!formValues?.number,
  }
}

export function normalize(
  operation: Partial<TGROperation> = {},
  context: Partial<TGROperation> = {}
): TGROperationUpdateBody {
  const casted = GROperationAPISchema.cast(operation, {
    stripUnknown: true,
    context: { values: operation, ...context },
  }) as TGROperationUpdateBody

  return {
    ...casted,
  }
}

export async function validate(
  operation: Partial<TGROperation> = {},
  context: Record<string, unknown> = {}
): Promise<{ [key: string]: string } | undefined> {
  try {
    await GROperationSchema.validate(operation, {
      stripUnknown: true,
      abortEarly: false,
      context: { values: operation, ...context },
    })
    return undefined
  } catch (e) {
    return getValidateFormErrors(e)
  }
}
