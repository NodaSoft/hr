import React from 'react'

import { MenuItem } from '@blueprintjs/core'

interface IProps {
  query: string
  minQueryLength: number
  loading: boolean
}

export default function NoResults({
  query,
  minQueryLength,
  loading,
}: IProps): JSX.Element {
  switch (true) {
    case query.length >= minQueryLength && !loading:
      return <MenuItem disabled text={'Совпадений не найдено'} />
    case query.length >= minQueryLength && loading:
      return <MenuItem disabled text={'...загрузка'} />
    default:
      return (
        <MenuItem
          disabled
          text={`Введите минимум ${minQueryLength}, чтобы начать поиск`}
        />
      )
  }
}
