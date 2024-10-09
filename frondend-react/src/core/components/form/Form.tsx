import * as React from 'react'
import { Field } from 'react-final-form'
import { OnChange } from 'react-final-form-listeners'

import { Icon, Intent, Tag } from '@blueprintjs/core'

import cn from 'classnames'
import { isEqual as _isEqual } from 'lodash-es'

import InlineControls from '../inline-controls/InlineControls'
import Styles from './Form.module.scss'

type ContainerProps = {
  children: React.ReactNode
}

export function preventSubmit(e: React.KeyboardEvent): void {
  if (e.keyCode === 13) e.preventDefault()
}

export function getDirtyFormData(
  values: { [key: string]: unknown },
  dirtyFieldsMap: { [key: string]: boolean } = {}
): { [key: string]: unknown } | null {
  const formData = {}
  const changedFields = Object.keys(dirtyFieldsMap)
  if (changedFields.length === 0) {
    return null
  }
  for (let i = 0; i < changedFields.length; i += 1) {
    const field = changedFields[i]

    // @ts-ignore
    formData[field] = values[field]
  }
  return formData
}

type ConditionProps = {
  when: string
  is?: boolean
  truly?: boolean
  falsy?: boolean
  isNot?: boolean
  isEqual?: boolean
} & ContainerProps

export function Condition(props: ConditionProps): JSX.Element {
  const { when, is, isNot, isEqual, truly, falsy, children } = props
  return (
    <Field
      allowNull
      name={when}
      subscription={{ value: true }}
      render={({ input: { value } }) => {
        let condition
        switch (true) {
          case truly !== undefined && truly && Boolean(value) === true:
            condition = true
            break
          case falsy !== undefined && falsy && Boolean(value) === false:
            condition = true
            break
          case is !== undefined && value === is:
            condition = true
            break
          case isNot !== undefined && value !== isNot:
            condition = true
            break
          case isEqual !== undefined && _isEqual(value, isEqual):
            condition = true
            break
          default:
            condition = false
        }
        return condition ? children : null
      }}
    />
  )
}

Condition.defaultProps = {
  is: undefined,
  isNot: undefined,
  isEqual: undefined,
  truly: undefined,
  falsy: undefined,
}

type WhenFieldChangesProps = {
  field: string
  becomes: unknown
  set: string
  to: unknown
}

export function WhenFieldChanges({
  field,
  becomes,
  set,
  to,
}: WhenFieldChangesProps): JSX.Element {
  return (
    <Field
      name={set}
      subscription={{}}
      render={(
        // No subscription. We only use Field to get to the change function
        { input: { onChange } }
      ) => (
        <OnChange name={field}>
          {(value) => {
            if (value === becomes) {
              onChange(to)
            }
          }}
        </OnChange>
      )}
    />
  )
}

type FormProps = {
  children?: React.ReactNode
  className?: string
  onSubmit?: (e: React.SyntheticEvent<HTMLFormElement>) => void
  fill?: boolean
  formId?: string
}

type NoteOrErrorProps = {
  note?: JSX.Element
  absolute?: boolean
  error?:
    | string
    | {
        message: string
        params: Record<string, string | number | boolean | undefined>
      }
}

export class HTMLForm extends React.PureComponent<FormProps> {
  static defaultProps = {
    children: undefined,
    onSubmit: undefined,
    className: undefined,
    fill: false,
  }

  static Buttons(props: ContainerProps): JSX.Element {
    return (
      <div className={Styles.FormButtons}>
        <InlineControls>{props.children}</InlineControls>
      </div>
    )
  }

  static RequiredSymbol(): JSX.Element {
    return (
      <Icon
        icon="asterisk"
        intent={Intent.DANGER}
        iconSize={6}
        className={Styles.FormAsterisk}
      />
    )
  }

  static NoteOrError(props: NoteOrErrorProps): JSX.Element | null {
    const { error, note, absolute } = props
    if (error) {
      const message = typeof error === 'string' ? error : error.message

      const ErrorTag = () => (
        <Tag minimal intent={Intent.DANGER}>
          {message}
        </Tag>
      )

      return absolute ? (
        <div className={Styles.Absolute}>
          <ErrorTag />
        </div>
      ) : (
        <ErrorTag />
      )
    }
    if (note) return note
    return null
  }

  handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>): void => {
    e.preventDefault()
    const { onSubmit } = this.props
    if (onSubmit !== undefined && typeof onSubmit === 'function') {
      onSubmit.call(null, e)
    }
  }

  render(): JSX.Element {
    const { children, className, fill } = this.props

    const classList = cn(Styles.Form, className, { [Styles.is_fill]: fill })

    return (
      <form
        className={classList}
        onSubmit={this.handleSubmit}
        id={this.props.formId}
      >
        {children}
      </form>
    )
  }
}
