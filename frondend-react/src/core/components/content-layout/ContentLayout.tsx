import React from 'react'

import cn from 'classnames'

import Styles from './ContentLayout.module.scss'

interface IProps {
  children: React.ReactNode
  responsive?: boolean
  fullWidth?: boolean
  fullHeight?: boolean
}

const ContentLayout: React.FC<IProps> = ({
  children,
  responsive = true,
  fullWidth = false,
  fullHeight = false,
}) => {
  const classnames = cn(Styles.Container, {
    [Styles.Responsive]: responsive,
    [Styles.FullWidth]: fullWidth,
    [Styles.FullHeight]: fullHeight,
  })

  return <div className={classnames}>{children}</div>
}

export default ContentLayout
