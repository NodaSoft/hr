import React, { MouseEventHandler } from 'react'

import {
  AnchorButton,
  IAnchorButtonProps,
  IconName,
  ITooltipProps,
  PopoverPosition,
  Tooltip,
} from '@blueprintjs/core'

import cn from 'classnames'

import Styles from 'core/layout/action-button/ActionButton.module.scss'

export interface IActionButtonProps extends Omit<IAnchorButtonProps, 'icon'> {
  rel?: string
  href?: string
  text?: string
  title?: string
  target?: string
  rounded?: boolean
  tabIndex?: number
  children?: React.ReactNode
  iconName?: IconName | JSX.Element
  tooltipProps?: Omit<ITooltipProps, 'content'>
  onMouseDown?: MouseEventHandler<HTMLAnchorElement>
}
export default function ActionButton({
  iconName,
  tabIndex,
  text = '',
  tooltipProps,
  rounded = false,
  ...buttonProps
}: IActionButtonProps): JSX.Element {
  const classNames = cn(buttonProps.className, { [Styles.Rounded]: rounded })

  const btnProps = {
    ...buttonProps,
    tabIndex,
    icon: iconName,
    className: classNames,
  } as const

  if (text?.length) {
    return (
      <Tooltip
        content={text}
        position={PopoverPosition.BOTTOM}
        disabled={tooltipProps?.disabled || btnProps.loading}
        {...tooltipProps}
      >
        <AnchorButton {...btnProps} />
      </Tooltip>
    )
  }

  return <AnchorButton {...btnProps} />
}
