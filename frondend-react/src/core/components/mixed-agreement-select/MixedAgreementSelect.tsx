import React from 'react'

import {
  Button,
  Classes,
  ControlGroup,
  Intent,
  IRef,
  MenuItem,
  Spinner,
  SpinnerSize,
  Tag,
} from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'
import { IItemRendererProps, Select } from '@blueprintjs/select'

import cn from 'classnames'

import { highlightText } from 'core/common/utils'
import { MixedAgreement } from 'core/models/agreement'
import ControlGroupStyles from 'core/styles/control-group.module.scss'
import SelectStyles from 'core/styles/select.module.scss'

import Styles from './Styles.module.scss'

const FakeFocusElem = React.forwardRef<HTMLDivElement>((_, ref) => (
  <div tabIndex={-1} ref={ref} />
))

const MIN_ITEMS_TO_SEARCH = 10

type MixedAgreementSelectItemProps = {
  mixedAgreement: MixedAgreement
  selected: boolean
  showBalance: boolean
} & IItemRendererProps

function MixedAgreementSelectItem({
  mixedAgreement,
  handleClick,
  modifiers,
  query,
  selected,
  showBalance,
}: MixedAgreementSelectItemProps) {
  if (!modifiers.matchesPredicate) {
    return null
  }

  const { label, balanceHV, balance } = mixedAgreement

  if (!modifiers.matchesPredicate) {
    return null
  }

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
    return (
      <div className={Styles.ItemTag}>
        {showBalance && (
          <Tag minimal={!modifiers.active} intent={intent}>
            {labeEl}
          </Tag>
        )}
      </div>
    )
  }

  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      icon={selected ? IconNames.TICK : IconNames.BLANK}
      onClick={handleClick}
      labelElement={tag(balanceHV, balanceIntent)}
      text={highlightText(label, query)}
    />
  )
}

export interface IMixedAgreementSelectProps {
  disabled?: boolean
  options: MixedAgreement[]
  onSelect: (
    item: MixedAgreement | null,
    event?: React.SyntheticEvent<HTMLElement>
  ) => void
  onClear: (item: MixedAgreement | null) => void
  value: MixedAgreement | null
  isClearButtonShow?: boolean
  loading?: boolean
  showBalance?: boolean
  elementRef?: IRef<HTMLButtonElement> | undefined
}

export default function MixedAgreementSelect({
  value = null,
  options,
  onSelect,
  onClear,
  loading = false,
  disabled = false,
  isClearButtonShow = true,
  showBalance = true,
  elementRef,
}: IMixedAgreementSelectProps): JSX.Element {
  const valueObj = React.useMemo(() => {
    if (value === null) return null
    return (
      options.find((item) => item.agreementId === value.agreementId) || value
    )
  }, [value, options])

  const itemPredicate = React.useCallback(
    (query: string, item: MixedAgreement) =>
      item.label.toLowerCase().includes(query.toLowerCase()),
    []
  )

  const itemRenderer = React.useCallback(
    (
      mixedAgreement: MixedAgreement,
      { handleClick, modifiers, query }: IItemRendererProps
    ) => (
      <MixedAgreementSelectItem
        selected={
          valueObj !== null &&
          valueObj.agreementId === mixedAgreement.agreementId
        }
        key={mixedAgreement.agreementId}
        handleClick={handleClick}
        modifiers={modifiers}
        query={query}
        mixedAgreement={mixedAgreement}
        showBalance={showBalance}
      />
    ),
    [valueObj]
  )

  const handleClear = () => {
    onSelect(null)
    onClear(value)
  }

  const canShowClearButton = isClearButtonShow && value && !disabled

  const selectRef = React.useRef<Select<MixedAgreement>>(null)
  const fakeFocusElemRef = React.useRef<HTMLDivElement>(null)

  const handleClosePopover = (
    e: React.KeyboardEvent<HTMLInputElement | HTMLButtonElement>
  ) => {
    if (
      selectRef.current &&
      selectRef.current.state.isOpen &&
      !e.currentTarget.value &&
      e.key === 'Tab'
    ) {
      selectRef.current.setState({ isOpen: false })
      if (fakeFocusElemRef.current) fakeFocusElemRef.current.focus()
    }
  }

  return (
    <>
      <FakeFocusElem ref={fakeFocusElemRef} />
      <ControlGroup>
        <Select<MixedAgreement>
          disabled={disabled}
          items={options}
          itemRenderer={itemRenderer}
          className={cn(SelectStyles.Select, {
            [SelectStyles.SelectWithReset]: canShowClearButton,
          })}
          filterable={options.length > MIN_ITEMS_TO_SEARCH}
          itemPredicate={itemPredicate}
          popoverProps={{
            wrapperTagName: 'div',
            targetTagName: 'div',
            boundary: 'viewport',
            popoverClassName: Styles.Popover,
          }}
          inputProps={{
            placeholder: '',
            onKeyDown: handleClosePopover,
          }}
          onItemSelect={onSelect}
          ref={selectRef}
        >
          <Button
            text={valueObj ? valueObj.label : ' '}
            fill
            disabled={disabled}
            onKeyDown={handleClosePopover}
            className={cn(
              Classes.ALIGN_LEFT,
              Classes.FILL,
              SelectStyles.SelectButton,
              ControlGroupStyles.ForceLeftControl,
              {
                [ControlGroupStyles.ForceRightControl]: !canShowClearButton,
              }
            )}
            rightIcon={
              loading ? (
                <Spinner size={SpinnerSize.SMALL} />
              ) : (
                IconNames.DOUBLE_CARET_VERTICAL
              )
            }
            elementRef={elementRef}
          />
        </Select>
        {canShowClearButton && (
          <Button icon={IconNames.CROSS} onClick={handleClear} />
        )}
      </ControlGroup>
    </>
  )
}
