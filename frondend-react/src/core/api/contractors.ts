import { delay } from 'core/common/utils'
import { MixedAgreementsResponse } from 'core/models/agreement'

async function getMixedAgreements(): Promise<MixedAgreementsResponse> {
  await delay(3000)

  return {
    data: [
      {
        agreementId: 249,
        relationType: 1,
        legalPersonId: 2100,
        isDefault: false,
        isDelete: false,
        contractorProfileId: 250793,
        contractorId: 251874,
        balanceHV: '0,00',
        balance: 0,
        repaymentPeriod: 0,
        value: '№DOGNUM-249 "test" - физ. лицо Север',
        label: '№DOGNUM-249 "test" - физ. лицо Север',
        isActive: true,
        currency: {
          isoCode: 'RUB',
          designation: '₽',
        },
        contractorLegalPersonId: 0,
        shopLegalPersonINN: '',
        contractorLegalPersonINN: '',
      },
      {
        agreementId: 252,
        relationType: 1,
        legalPersonId: 5172,
        isDefault: false,
        isDelete: false,
        contractorProfileId: 250793,
        contractorId: 251874,
        balanceHV: '-20,00',
        balance: -20,
        repaymentPeriod: 0,
        value: '№DOGNUM-252 "Активный договор" - физ. лицо Север',
        label: '№DOGNUM-252 "Активный договор" - физ. лицо Север',
        isActive: true,
        currency: {
          isoCode: 'RUB',
          designation: '₽',
        },
        contractorLegalPersonId: 0,
        shopLegalPersonINN: '',
        contractorLegalPersonINN: '',
      },
      {
        agreementId: 2654,
        relationType: 1,
        legalPersonId: 2099,
        isDefault: false,
        isDelete: false,
        contractorProfileId: 250793,
        contractorId: 251874,
        balanceHV: '0,00',
        balance: 0,
        repaymentPeriod: 0,
        value: '№DOGNUM-2654 "конвертация УПД" - ОПТ13',
        label: '№DOGNUM-2654 "конвертация УПД" - ОПТ13',
        isActive: true,
        currency: {
          isoCode: 'RUB',
          designation: '₽',
        },
        contractorLegalPersonId: 26487,
        shopLegalPersonINN: '',
        contractorLegalPersonINN: '',
      },
    ],
  }
}

export default { getMixedAgreements }
