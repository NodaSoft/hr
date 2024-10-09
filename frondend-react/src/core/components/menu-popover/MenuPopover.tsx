import React, { useEffect, useRef, useState } from 'react'

import { Menu, Position } from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import { TOption } from 'core/models/common'

import { usePopoverListHotkeys } from '../hotkeys/usePopoverListHotkeys'
import { MenuPopoverItem } from './MenuPopoverItem'
import Styles from './Styles.module.scss'

interface IProps<T> {
  menuItem?: React.FC<T>
  selectedItem: number | string | null
  onItemClick: (item: T) => void
  menuItems: Array<T>
  children: React.FC<{ ref: React.MutableRefObject<HTMLAnchorElement | null> }>
  position?: Position
  isPopoverOpen: boolean
  setPopoverOpen: React.Dispatch<React.SetStateAction<boolean>>
}

export const MenuPopover = <T extends TOption>({
  menuItems,
  menuItem,
  selectedItem,
  onItemClick,
  children,
  position = Position.BOTTOM,
  isPopoverOpen,
  setPopoverOpen,
}: IProps<T>): JSX.Element => {
  const popoverInteraction = (
    state: boolean,
    e?: React.SyntheticEvent<HTMLElement>
  ) => {
    if (e?.isTrusted) {
      setPopoverOpen(state)
    }
  }

  const [activeIndex, setActiveIndex] = useState<number | null>(null)

  const childComponentRef = useRef<HTMLAnchorElement | null>(null)
  const childComponentFocusTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  const setFocusOnChildComponent = () => {
    childComponentFocusTimeoutRef.current = setTimeout(() => {
      if (childComponentRef.current) {
        childComponentRef.current.focus()
      }
    }, 100)
  }

  const onClosePopover = () => {
    setPopoverOpen(false)
    setActiveIndex(null)
    setFocusOnChildComponent()
  }

  const onOpenPopover = () => {
    if (document.activeElement === childComponentRef.current) {
      setPopoverOpen(true)
    }
  }

  const { handleKeyDown } = usePopoverListHotkeys({
    itemsLength: menuItems.length,
    isPopoverOpen,
    activeIndex,
    setActiveIndex,
    onClosePopover,
    onOpenPopover,
  })

  const clearTimeouts = () => {
    if (childComponentFocusTimeoutRef.current) {
      clearTimeout(childComponentFocusTimeoutRef.current)
    }
  }

  useEffect(() => {
    return () => {
      clearTimeouts()
    }
  }, [])

  return menuItems.length ? (
    <Popover2
      position={position}
      isOpen={isPopoverOpen}
      onInteraction={popoverInteraction}
      canEscapeKeyClose={false}
      onOpened={onOpenPopover}
      onClosed={onClosePopover}
      content={
        <Menu className={Styles.Menu} onKeyDown={handleKeyDown}>
          {menuItems.map((item, index) => (
            <MenuPopoverItem<typeof item>
              key={item.id}
              isActive={index === activeIndex}
              item={item}
              isOpen={isPopoverOpen}
              onClick={onItemClick}
              content={menuItem ? menuItem(item) : item.name}
              isSelected={item.id === selectedItem}
              closePopover={onClosePopover}
            />
          ))}
        </Menu>
      }
      shouldReturnFocusOnClose={false}
    >
      {children({ ref: childComponentRef })}
    </Popover2>
  ) : (
    <></>
  )
}
