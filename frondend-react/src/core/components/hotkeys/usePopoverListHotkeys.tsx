import React, { useCallback, useMemo } from 'react'

import { useHotkeys } from '@blueprintjs/core'

type IProps = {
  isPopoverOpen: boolean
  activeIndex: number | null
  setActiveIndex: (activeIndex: number | null) => void
  onClosePopover?: () => void
  onOpenPopover?: () => void
  itemsLength: number
}

type ReturnProps = {
  handleKeyDown?: (e: React.KeyboardEvent<HTMLElement>) => void
}

export const usePopoverListHotkeys = ({
  isPopoverOpen,
  activeIndex,
  setActiveIndex,
  itemsLength,
  onClosePopover,
  onOpenPopover,
}: IProps): ReturnProps => {
  const setActiveElement = useCallback(
    (direction: 'up' | 'down') => {
      if (activeIndex === null) {
        setActiveIndex(0)
      }

      if (direction === 'down' && activeIndex !== null) {
        const nextIndex =
          activeIndex === itemsLength - 1 ? itemsLength - 1 : activeIndex + 1
        setActiveIndex(nextIndex)
      }

      if (direction === 'up' && activeIndex !== null) {
        const prevIndex = activeIndex === 0 ? 0 : activeIndex - 1
        setActiveIndex(prevIndex)
      }
    },
    [activeIndex]
  )

  const hotkeys = useMemo(
    () => [
      {
        combo: 'down',
        global: true,
        label: 'common.hotkeys.popover.down',
        onKeyDown: () => setActiveElement('down'),
        disabled: !isPopoverOpen,
      },
      {
        combo: 'up',
        global: true,
        label: 'common.hotkeys.popover.up',
        onKeyDown: () => setActiveElement('up'),
        disabled: !isPopoverOpen,
      },
      {
        combo: 'tab',
        global: true,
        label: 'common.hotkeys.popover.close',
        onKeyDown: onClosePopover,
        disabled: !isPopoverOpen,
      },
      {
        combo: 'esc',
        global: true,
        label: 'common.hotkeys.popover.close',
        onKeyDown: onClosePopover,
        disabled: !isPopoverOpen,
      },
      {
        combo: 'enter',
        global: true,
        label: 'common.hotkeys.popover.open',
        onKeyDown: onOpenPopover,
        disabled: isPopoverOpen || !onOpenPopover,
      },
    ],
    [isPopoverOpen, setActiveElement]
  )

  const { handleKeyDown } = useHotkeys(hotkeys)

  return { handleKeyDown }
}
