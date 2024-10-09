import React from 'react'

import cn from 'classnames'

import Styles from 'core/components/page-header/PageHeader.module.scss'

interface IProps {
  className?: string
  children: React.ReactNode
}

const PageHeader: React.FC<IProps> = ({ children, className }) => {
  const classList = cn(className, Styles.PageHeader)
  return <div className={classList}>{children}</div>
}

export default PageHeader
