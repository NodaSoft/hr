import React from 'react'
import { useField } from 'react-final-form'

import { FormGroup, IRef } from '@blueprintjs/core'

import { HTMLForm, MixedAgreementSelect } from 'core/components'
import { IMixedAgreementSelectProps } from 'core/components/mixed-agreement-select/MixedAgreementSelect'
import { TFieldProps } from 'core/interfaces/forms'
import { MixedAgreement } from 'core/models/agreement'
import { ContractorTypes } from 'core/models/contractor'

type TInputProps = Omit<
  IMixedAgreementSelectProps,
  'onSelect' | 'disabled' | 'value' | 'onClear' | 'options'
>

interface IProps extends TFieldProps {
  contractorId: number
  contractorType: ContractorTypes
  inputProps?: TInputProps
  elementRef?: IRef<HTMLButtonElement> | undefined
  onlyActive?: boolean
  options: MixedAgreement[]
  onChange?: (value: MixedAgreement | null) => void
}

export default function MixedAgreementsSelectField({
  name,
  label = '',
  className,
  inputProps,
  disabled = false,
  required = false,
  allowNull = false,
  elementRef,
  options,
  onChange: externalOnChange,
}: IProps): JSX.Element {
  const {
    input: { onChange, value },
    meta: { error, touched, submitting },
  } = useField<MixedAgreement | null>(name, { allowNull })

  const handleClear = () => {
    onChange(null)
    if (externalOnChange) {
      externalOnChange(null)
    }
  }

  const handleSelect = (selectedValue: MixedAgreement | null) => {
    onChange(selectedValue)
    if (externalOnChange) {
      externalOnChange(selectedValue) // Вызовите внешнюю onChange функцию
    }
  }

  return (
    <FormGroup
      label={label}
      labelInfo={required && <HTMLForm.RequiredSymbol />}
      helperText={<HTMLForm.NoteOrError error={touched ? error : ''} />}
      className={className}
    >
      <MixedAgreementSelect
        {...inputProps}
        options={options}
        value={value}
        onSelect={handleSelect}
        onClear={handleClear}
        disabled={submitting || disabled}
        elementRef={elementRef}
      />
    </FormGroup>
  )
}
