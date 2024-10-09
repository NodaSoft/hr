import * as React from 'react'

import { Callout, IconName, Intent } from '@blueprintjs/core'

import cn from 'classnames'

import { isApiError } from 'core/common/api'
import { errorTrace } from 'core/common/errorHandling'

import ErrorTrace from '../error-trace/ErrorTrace'
import Styles from './ErrorAlert.module.scss'

type IProps = {
  error: unknown
  text?: string
  icon?: IconName
  children?: React.ReactNode
  className?: string
}

const ErrorAlert: React.FC<IProps> = ({
  error,
  text,
  icon,
  children,
  className,
}) => {
  const trace = React.useMemo(
    () => isApiError(error) && errorTrace(error),
    [error]
  )

  return (
    <div className={cn(Styles.Container, className)}>
      <Callout intent={Intent.DANGER} icon={icon}>
        {trace ? (
          <>
            <p>{text}</p>
            {children}
            <ErrorTrace trace={trace} />
          </>
        ) : (
          <>
            {text}
            {children}
          </>
        )}
      </Callout>
    </div>
  )
}

export default ErrorAlert
