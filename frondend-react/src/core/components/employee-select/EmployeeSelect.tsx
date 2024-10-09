import * as React from 'react'

import { Button, Classes, ControlGroup, MenuItem } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'
import { IItemRendererProps, Select } from '@blueprintjs/select'

import cn from 'classnames'

import { AbstractNamedEntry, IEmployee } from 'core/api/types'
import { highlightText } from 'core/common/utils'
import ControlGroupStyles from 'core/styles/control-group.module.scss'
import SelectStyles from 'core/styles/select.module.scss'

const FakeFocusElem = React.forwardRef<HTMLDivElement>((_, ref) => (
  <div tabIndex={-1} ref={ref} />
))

const MIN_ITEMS_TO_SEARCH = 10

type EmployeeSelectItemProps = {
  employee: IEmployee
  selected: boolean
} & IItemRendererProps

function EmployeeSelectItem({
  employee,
  handleClick,
  modifiers,
  query,
  selected,
}: EmployeeSelectItemProps) {
  if (!modifiers.matchesPredicate) {
    return null
  }
  const text = `${employee.name}`
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      icon={selected ? IconNames.TICK : IconNames.BLANK}
      onClick={handleClick}
      text={highlightText(text, query)}
    />
  )
}

type EmployeeSelectProps = {
  disabled?: boolean
  options: Array<IEmployee>
  onSelect: (
    item: null | IEmployee,
    event?: React.SyntheticEvent<HTMLElement>
  ) => void
  value: IEmployee | null
  isClearButtonShow?: boolean
  inputRef?: { current: HTMLInputElement | null }
  buttonRef?: { current: HTMLButtonElement | null }
}

const EmployeeSelect: React.FC<EmployeeSelectProps> = ({
  value = null,
  options,
  onSelect,
  disabled = false,
  isClearButtonShow = true,
  inputRef,
  buttonRef,
}) => {
  const valueObj = React.useMemo(() => {
    if (value === null) return null
    return options.find((item) => item.id === value.id) || value
  }, [value, options])

  const itemPredicate = React.useCallback(
    (query: string, item: IEmployee) =>
      item.name.toLowerCase().includes(query.toLowerCase()),
    []
  )

  const itemRenderer = React.useCallback(
    (
      employee: IEmployee,
      { handleClick, modifiers, query }: IItemRendererProps
    ) => (
      <EmployeeSelectItem
        selected={valueObj !== null && valueObj.id === employee.id}
        key={employee.id}
        handleClick={handleClick}
        modifiers={modifiers}
        query={query}
        employee={employee}
      />
    ),
    [valueObj]
  )

  const canShowClearButton = isClearButtonShow && value && !disabled

  const selectRef = React.useRef<Select<IEmployee>>(null)
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

  const [activeItem, setActiveItem] = React.useState(valueObj)

  const handleActiveItemChange = (item: AbstractNamedEntry | null) => {
    setActiveItem(item)
  }

  return (
    <>
      <FakeFocusElem ref={fakeFocusElemRef} />
      <ControlGroup>
        <Select
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
          }}
          inputProps={{
            placeholder: '',
            onKeyDown: handleClosePopover,
            inputRef,
          }}
          onItemSelect={onSelect}
          ref={selectRef}
          onActiveItemChange={handleActiveItemChange}
          activeItem={activeItem}
        >
          <Button
            text={valueObj ? valueObj.name : ' '}
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
            rightIcon="double-caret-vertical"
            elementRef={buttonRef}
          />
        </Select>
        {canShowClearButton && (
          <Button icon="cross" onClick={() => onSelect(null)} tabIndex={-1} />
        )}
      </ControlGroup>
    </>
  )
}

export default EmployeeSelect
