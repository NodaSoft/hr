import React from 'react'

import { Intent, NonIdealState, Spinner } from '@blueprintjs/core'

interface IProps {
  title?: string
  description?: string
  value?: number
}

const NotReadyState: React.FC<IProps> = ({ title, description, value }) => {
  return (
    <NonIdealState
      title={title || 'Подождите'}
      description={description || 'Идет загрузка данных'}
      icon={<Spinner intent={Intent.PRIMARY} value={value} />}
    />
  )
}

export default NotReadyState
