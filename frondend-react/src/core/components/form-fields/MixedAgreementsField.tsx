import React from 'react'
import { useField } from 'react-final-form'

import { FormGroup } from '@blueprintjs/core'

import { isFunction } from 'lodash-es'

import { HTMLForm, MixedAgreementsAutocomplete } from 'core/components'
import { TFieldProps } from 'core/interfaces/forms'
import { MixedAgreement } from 'core/models/agreement'

import { IMixedAgreementsAutocompleteProps } from '../mixed-agreements-autocomplete/MixedAgreementsAutocomplete'

type TInputProps = Omit<
  IMixedAgreementsAutocompleteProps,
  'onSelect' | 'disabled' | 'value' | 'onClear'
>

interface IProps extends TFieldProps {
  inputProps?: TInputProps
  onChange?: (value: MixedAgreement | null) => void
}

export default function MixedAgreementsField({
  name,
  note,
  onChange,
  className,
  label = '',
  inputProps,
  disabled = false,
  required = false,
  allowNull = false,
}: IProps): JSX.Element {
  const {
    input: { onChange: onInputChange, value, onFocus },
    meta: { error, touched, submitting },
  } = useField<MixedAgreement | null>(name, { allowNull })

  const handleChangeInput = (newValue: MixedAgreement | null) => {
    onInputChange(newValue)

    if (isFunction(onChange)) {
      onChange(newValue)
    }
  }

  const handleClear = () => {
    handleChangeInput(null)
  }

  return (
    <FormGroup
      label={label}
      labelInfo={required && <HTMLForm.RequiredSymbol />}
      helperText={
        <HTMLForm.NoteOrError note={note} error={touched ? error : ''} />
      }
      className={className}
    >
      <MixedAgreementsAutocomplete
        {...inputProps}
        value={value}
        onFocus={onFocus}
        onClear={handleClear}
        onSelect={handleChangeInput}
        disabled={submitting || disabled}
      />
    </FormGroup>
  )
}
