import React from 'react'

import { Breadcrumbs } from '@blueprintjs/core'

export default function NewOperationBreadcrumbs() {
  const items = [
    {
      text: 'Приемка',
    },
    {
      text: 'Создание приемки',
      current: true,
    },
  ]

  return <Breadcrumbs items={items} />
}
