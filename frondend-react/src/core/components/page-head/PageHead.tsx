import * as React from 'react'

import Styles from './PageHead.module.scss'

interface IProps {
  children?: React.ReactNode
}

const PageHead: React.FC<IProps> = ({ children }) => {
  return <div className={Styles.Head}>{children}</div>
}

export default PageHead
