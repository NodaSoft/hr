import { IDatePickerLocaleUtils } from '@blueprintjs/datetime'

import moment from 'moment'
import 'moment-timezone'

import { Locale } from 'core/models/localization'

const locales = {
  [Locale.ru_RU]: async () => {
    await import('moment/locale/ru')
    moment.locale('ru')
  },
  [Locale.en_US]: () => moment.locale(),
  [Locale.uk_UA]: async () => {
    await import('moment/locale/uk')
    moment.locale('uk')
  },
  [Locale.ru_UA]: async () => {
    await import('moment/locale/ru')
    moment.locale('ru')
  },
  [Locale.ru_KZ]: async () => {
    await import('moment/locale/ru')
    moment.locale('ru')
  },
}

async function setLocale(locale: Locale): Promise<void> {
  if (typeof locale !== 'string') {
    throw new TypeError('locale should be a string')
  }

  if (!Object.keys(locales).includes(locale)) {
    throw new Error('locale is unsupported')
  }

  const localeSetter = locales[locale]

  await localeSetter() // eslint-disable-line
}

const localeUtils = {
  formatDay: (day: string): string => moment(day).toISOString(),
  formatMonthTitle: (month: number): string =>
    moment().localeData().months()[month],
  formatWeekdayShort: (weekday: number): string =>
    moment().localeData().weekdaysShort()[weekday],
  formatWeekdayLong: (weekday: number): string =>
    moment().localeData().weekdays()[weekday],
  getMonths: (): string[] => moment().localeData().months(),
  getFirstDayOfWeek: (): number => moment().localeData().firstDayOfWeek(),
} as unknown as IDatePickerLocaleUtils

export { localeUtils, setLocale }
