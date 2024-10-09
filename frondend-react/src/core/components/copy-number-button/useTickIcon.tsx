import React, { useEffect, useRef, useState } from 'react'

import { Icon } from '@blueprintjs/core'
import { IconName, IconNames } from '@blueprintjs/icons'

interface IProps {
  icon: IconName
  iconSize?: number
}

interface IReturnProps {
  iconWithTick: JSX.Element
  onIconClick: (callback: () => void) => void
}

export const TICK_TIMEOUT = 1000

export const useTickIcon = ({ iconSize, icon }: IProps): IReturnProps => {
  const [isShowTick, setShowTick] = useState<boolean>(false)
  const iconName = isShowTick ? IconNames.TICK : icon

  const tickTimeout = useRef(0)

  const onIconClick = (callback: () => void) => {
    if (isShowTick) return

    callback()
    setShowTick(true)

    tickTimeout.current = window.setTimeout(() => {
      setShowTick(false)
    }, TICK_TIMEOUT)
  }

  const iconWithTick = <Icon icon={iconName} size={iconSize} />

  useEffect(() => {
    return () => {
      window.clearTimeout(tickTimeout.current)
    }
  }, [])

  return {
    iconWithTick,
    onIconClick,
  }
}
