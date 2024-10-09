import React from 'react'

import { Intent, MenuItem, Tag } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'
import { IItemRendererProps } from '@blueprintjs/select'

import { highlightText } from 'core/common/utils'

import { Option } from '../MixedAgreementsAutocomplete'

interface IProps extends IItemRendererProps {
  selected: boolean
  showBalance: boolean
  mixedAgreement: Option
}

export default function ItemRenderer({
  query,
  selected,
  modifiers,
  handleClick,
  showBalance,
  mixedAgreement,
}: IProps): JSX.Element | null {
  if (!modifiers.matchesPredicate) {
    return null
  }

  // @ts-ignore
  const { label, balanceHV, balance } = mixedAgreement

  const balanceIsNegative = balance < 0
  const balanceIsPositive = balance > 0

  const balanceIntent = (() => {
    switch (true) {
      case balanceIsNegative: {
        return Intent.DANGER
      }
      case balanceIsPositive: {
        return Intent.SUCCESS
      }
      case modifiers.active: {
        return Intent.PRIMARY
      }
      default: {
        return Intent.NONE
      }
    }
  })()

  const tag = (labeEl: React.ReactNode, intent: Intent) => {
    if (!showBalance) return

    return (
      <Tag minimal={!modifiers.active} intent={intent}>
        {labeEl}
      </Tag>
    )
  }

  const getIcon = () => {
    switch (true) {
      case selected: {
        return IconNames.TICK
      }
      default: {
        return IconNames.BLANK
      }
    }
  }

  const icon = getIcon()

  return (
    <MenuItem
      icon={icon}
      onClick={handleClick}
      active={modifiers.active}
      disabled={modifiers.disabled}
      text={highlightText(label, query)}
      labelElement={tag(balanceHV, balanceIntent)}
    />
  )
}
