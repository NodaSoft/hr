import React from 'react'

import { Button, IconSize } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'

import { copyStringToClipboard } from 'core/common/utils'
import { AppToaster } from 'core/components'

import { TICK_TIMEOUT, useTickIcon } from './useTickIcon'

interface IProps {
  number: string
  small?: boolean
  minimal?: boolean
  className?: string
  iconSize?: IconSize
}

export default function CopyNumberButton({
  number,
  className,
  small = true,
  minimal = true,
  iconSize = IconSize.STANDARD,
}: IProps): JSX.Element {
  const { iconWithTick, onIconClick } = useTickIcon({
    iconSize,
    icon: IconNames.DUPLICATE,
  })

  const handleButtonClick = (
    event: React.MouseEvent<HTMLElement, MouseEvent>
  ) => {
    event.stopPropagation()

    const copyToClipboard = () => copyStringToClipboard(number)

    onIconClick(copyToClipboard)

    AppToaster.success({
      message: 'Скопировано',
      timeout: TICK_TIMEOUT,
    })
  }

  const handleDoubleClick = (
    event: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    event.stopPropagation()
  }

  const htmlTitle = 'common.actions.copy.number'

  return (
    <Button
      icon={iconWithTick}
      small={small}
      tabIndex={-1}
      minimal={minimal}
      title={htmlTitle}
      className={className}
      onClick={handleButtonClick}
      onDoubleClick={handleDoubleClick}
    />
  )
}
