import * as React from 'react'

import cn from 'classnames'

import Styles from './PageLayout.module.scss'

interface IProps {
  children?: React.ReactNode
  overflowed?: boolean
  id?: string
  className?: string
}

const PageLayout: React.FC<IProps> = ({
  id,
  children,
  className,
  overflowed = false,
}) => {
  return (
    <div
      id={id}
      className={cn(className, Styles.Page, {
        [Styles.PageOverflowed]: overflowed,
      })}
    >
      {children}
    </div>
  )
}

export default PageLayout
