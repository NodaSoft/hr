import React, { useCallback, useMemo } from 'react'

import { MenuItem, useHotkeys } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'

import { TOption } from 'core/models/common'

import Styles from './Styles.module.scss'

interface IProps<T> {
  isOpen: boolean
  onClick: (item: T) => void
  closePopover: () => void
  item: T
  isSelected: boolean
  isActive: boolean
  content: React.ReactElement | string | null
}

export const MenuPopoverItem = <T extends TOption>({
  isActive,
  isOpen,
  item,
  isSelected,
  content,
  onClick,
  closePopover,
}: IProps<T>): JSX.Element => {
  const handleItemClick = () => {
    if (!isSelected) onClick(item)
    closePopover()
  }

  const handleEnterClick = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      if (isActive) {
        handleItemClick()
      }
    },
    [isActive, handleItemClick]
  )

  const hotkeys = useMemo(
    () => [
      {
        combo: 'enter',
        global: true,
        label: { content },
        onKeyDown: handleEnterClick,
        disabled: !isOpen,
      },
    ],
    [isOpen, handleEnterClick]
  )

  const { handleKeyDown } = useHotkeys(hotkeys)

  return (
    <MenuItem
      onKeyDown={handleKeyDown}
      active={isActive}
      className={Styles.MenuItem}
      onClick={handleItemClick}
      icon={isSelected ? IconNames.TICK : IconNames.BLANK}
      text={content}
    />
  )
}
