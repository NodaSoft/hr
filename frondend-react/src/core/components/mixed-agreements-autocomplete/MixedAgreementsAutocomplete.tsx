import React, { useRef, useState } from 'react'

import {
  Button,
  Classes,
  ControlGroup,
  InputGroupProps2,
  Intent,
  IPopoverProps,
  Spinner,
} from '@blueprintjs/core'
import { IconNames } from '@blueprintjs/icons'
import { Suggest } from '@blueprintjs/select'

import cn from 'classnames'
import { isEqual, isObject } from 'lodash-es'
import { v4 } from 'uuid'

import ContractorsApi from 'core/api/contractors'
import logger from 'core/common/logger'
import { preventSubmit } from 'core/components/form/Form'
import { MixedAgreement } from 'core/models/agreement'
import { ContractorTypes } from 'core/models/contractor'
import ControlGroupStyles from 'core/styles/control-group.module.scss'

import { NoResults, ItemRenderer } from './components'
import Styles from './Styles.module.scss'

export const OPTION_CREATE = 'create'
export type Option = MixedAgreement | typeof OPTION_CREATE

type Value = MixedAgreement | null

export interface IMixedAgreementsAutocompleteProps {
  onSelect: (value: Value) => void
  onClear?: (value: Value) => void
  onFocus?: (evt: React.FocusEvent<HTMLInputElement>) => void
  value: Value
  contractorType?: ContractorTypes
  query?: string
  disabled?: boolean
  minQueryLength?: number
  autoFocus?: boolean
  loading?: boolean
  loadImmediately?: boolean
  showBalance?: boolean
  isClearButtonShow?: boolean
  allowCreate?: boolean
  inputRef?: { current: HTMLInputElement | null }
}

export default function MixedAgreementsAutocomplete({
  onSelect,
  onClear = () => undefined,
  onFocus = () => undefined,
  value = null,
  allowCreate = true,
  contractorType = ContractorTypes.CLIENT,
  disabled = false,
  minQueryLength = 3,
  autoFocus = false,
  loading: forceLoading = false,
  query: initialQuery = '',
  loadImmediately = false,
  showBalance = true,
  isClearButtonShow = true,
  inputRef,
}: IMixedAgreementsAutocompleteProps): JSX.Element {
  const [query, setQuery] = useState<string>(initialQuery)
  const [options, setOptions] = useState<Option[]>([])
  const [loading, setLoading] = useState<boolean>(false)

  const searchTimeout = useRef<NodeJS.Timeout>()
  const localInputRef = useRef<HTMLInputElement | null>(null)
  const allowedEmptyQuery = minQueryLength === 0
  const defaultKey = v4()

  const loadMixedAgreements = async () => {
    try {
      const { data } = await ContractorsApi.getMixedAgreements()

      const agreements = Array.isArray(data) ? data : []

      return agreements
    } catch (e) {
      logger.error(e)
      return []
    }
  }

  const handleSearch = async () => {
    setLoading(true)

    const mixedAgreements = await loadMixedAgreements()

    const canCreate = query.length >= minQueryLength && allowCreate

    const agreementsWithCreateField =
      allowedEmptyQuery && canCreate
        ? [...mixedAgreements, OPTION_CREATE]
        : mixedAgreements

    // @ts-ignore
    setOptions(agreementsWithCreateField)

    if (mixedAgreements.length === 1) {
      onSelect(mixedAgreements[0])
    } else if (autoFocus && localInputRef.current) {
      localInputRef.current.focus()
    }

    setLoading(false)
  }

  React.useEffect(() => {
    if (loadImmediately && !disabled && !options.length) {
      handleSearch()
    }
  }, [loadImmediately, disabled])

  const handleChangeQuery = (nextQuery: string) => {
    // @ts-ignore
    clearTimeout(searchTimeout.current)
    setQuery(nextQuery)
    if (nextQuery.length >= minQueryLength) {
      setOptions([])
      setLoading(true)

      searchTimeout.current = setTimeout(() => {
        handleSearch()
      }, 600)
    } else {
      setLoading(false)
      setOptions([])
    }
  }

  const handleItemSelect = (item: Option) => {
    if (item === OPTION_CREATE) {
      return // Здесь было создание
    }

    onSelect(item)
  }

  const handleInputFocus = () => {
    if (localInputRef.current) {
      localInputRef.current.focus()
    }
  }

  const handleReset = (): void => {
    // @ts-ignore
    clearTimeout(searchTimeout.current)
    setQuery('')
    setLoading(false)
    if (!loadImmediately) {
      setOptions([])
    }
    onClear(value)
    handleInputFocus()
  }

  const inputValueRenderer = (item: Option) => {
    if (isObject(item)) {
      return item.value
    }

    return ''
  }

  const inputClassnames = cn(Classes.FILL, {
    [ControlGroupStyles.ForceSingleControl]: !value || disabled,
    [ControlGroupStyles.ForceLeftControl]: value,
  })

  const spinner =
    loading || forceLoading ? (
      <Spinner intent={Intent.PRIMARY} size={16} />
    ) : undefined

  const inputProps: InputGroupProps2 = {
    onFocus,
    placeholder: '',
    rightElement: spinner,
    onKeyDown: preventSubmit,
    title: value?.value || '',
    leftIcon: IconNames.SEARCH,
    className: inputClassnames,
    inputRef: (input) => {
      localInputRef.current = input

      if (inputRef) {
        inputRef.current = input
      }
    },
  }

  const popoverProps: Partial<IPopoverProps> = {
    wrapperTagName: 'div',
    targetTagName: 'div',
    boundary: 'viewport',
    popoverClassName: Styles.Popover,
  }

  return (
    <>
      <ControlGroup>
        <div className={Classes.FILL}>
          <Suggest<Option>
            query={query}
            items={options}
            disabled={disabled}
            selectedItem={value}
            inputProps={inputProps}
            popoverProps={popoverProps}
            onItemSelect={handleItemSelect}
            onQueryChange={handleChangeQuery}
            inputValueRenderer={inputValueRenderer}
            itemRenderer={(mixedAgreement, itemRendererProps) => {
              const key = isObject(mixedAgreement)
                ? mixedAgreement.agreementId
                : defaultKey

              const isSelected =
                isObject(mixedAgreement) &&
                isEqual(value?.agreementId, mixedAgreement.agreementId)

              return (
                <ItemRenderer
                  {...itemRendererProps}
                  key={key}
                  selected={isSelected}
                  showBalance={showBalance}
                  mixedAgreement={mixedAgreement}
                  // @ts-ignore
                  contractorType={contractorType}
                />
              )
            }}
            noResults={
              <NoResults
                query={query}
                loading={loading}
                minQueryLength={minQueryLength}
              />
            }
          />
        </div>
        {value && !disabled && isClearButtonShow && (
          <Button icon={IconNames.CROSS} onClick={handleReset} tabIndex={-1} />
        )}
      </ControlGroup>
    </>
  )
}
