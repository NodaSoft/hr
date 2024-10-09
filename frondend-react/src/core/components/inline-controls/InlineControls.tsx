import React from 'react'

import cn from 'classnames'

import { FlexAlignItems } from 'core/types/styles'

import Styles from './InlineControls.module.scss'

interface IProps {
  children?: React.ReactNode
  fill?: boolean
  text?: boolean
  className?: string
  alignment?: FlexAlignItems
}

const InlineControls: React.FC<IProps> = ({
  children,
  fill = false,
  alignment = FlexAlignItems.START,
  text = false,
  className,
}) => {
  const AlignmentClass = React.useMemo(() => {
    switch (alignment) {
      case 'baseline':
        return Styles.AlignBaseLine
      case 'center':
        return Styles.AlignCenter
      case 'end':
        return Styles.AlignEnd
      default:
        return Styles.AlignStart
    }
  }, [alignment])
  return (
    <div
      className={cn(Styles.InlineControls, AlignmentClass, className, {
        [Styles.Fill]: fill,
        [Styles.Text]: text,
      })}
    >
      {children}
    </div>
  )
}

export default InlineControls
