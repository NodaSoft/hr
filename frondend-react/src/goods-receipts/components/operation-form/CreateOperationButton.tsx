import React, { useState } from 'react'

import { Button, ButtonGroup, Intent, IRef } from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'

import appStorage from 'core/common/storage'
import { MenuPopover } from 'core/components'
import { ActionButton } from 'core/layout'

import {
  CREATE_NEW_OP_OPTIONS,
  CreateNewOpTypes,
  DEFAULT_CREATE_NEW_OPTION,
  DEFAULT_CREATE_NEW_TYPE,
  GR_CREATE_NEW_OP_TYPE_LOCAL_STORAGE_KEY,
} from './common'

interface IProps {
  loading: boolean
  disabled: boolean
  elementRef: IRef<HTMLButtonElement>
  onSubmit: () => void
}

type TOption = {
  id: CreateNewOpTypes
  name: string
}

export const CreateOperationButton = ({
  onSubmit,
  loading,
  disabled,
  elementRef,
}: IProps): JSX.Element => {
  const [isPopoverOpen, setPopoverOpen] = useState<boolean>(false)

  const canAddPosByOrder = (id: CreateNewOpTypes) =>
    id !== CreateNewOpTypes.ADD_POSITION_BY_ORDER
  const canAddPosByNewOrder = (id: CreateNewOpTypes) =>
    id !== CreateNewOpTypes.ADD_POSITION_BY_NEW_ORDER

  const filteredOptions = CREATE_NEW_OP_OPTIONS.filter(
    (option) => canAddPosByOrder(option.id) && canAddPosByNewOrder(option.id)
  )

  const savedCreateType =
    appStorage.getItem(GR_CREATE_NEW_OP_TYPE_LOCAL_STORAGE_KEY) ||
    DEFAULT_CREATE_NEW_TYPE

  const initialOption =
    filteredOptions.find((option) => option.id === savedCreateType) ||
    DEFAULT_CREATE_NEW_OPTION

  const [currentOption, setCurrentOption] = useState<TOption>(initialOption)

  const handleSubmit = (type: string) => {
    appStorage.setItem(GR_CREATE_NEW_OP_TYPE_LOCAL_STORAGE_KEY, type)
    onSubmit()
  }

  const onItemClick = (option: TOption) => {
    handleSubmit(option.id)
    setCurrentOption(option)
  }

  return (
    <ButtonGroup>
      <Button
        intent={Intent.PRIMARY}
        onClick={() => handleSubmit(currentOption.id)}
        loading={loading}
        disabled={disabled}
        text={currentOption.name}
        elementRef={elementRef}
      />
      <MenuPopover<TOption>
        menuItems={filteredOptions}
        onItemClick={onItemClick}
        selectedItem={currentOption.id}
        isPopoverOpen={isPopoverOpen}
        setPopoverOpen={setPopoverOpen}
      >
        {({ ref }) => (
          <ActionButton
            elementRef={ref}
            intent={Intent.PRIMARY}
            iconName={IconNames.CHEVRON_DOWN}
            disabled={disabled || loading}
            tooltipProps={{
              openOnTargetFocus: isPopoverOpen,
            }}
          />
        )}
      </MenuPopover>
    </ButtonGroup>
  )
}
