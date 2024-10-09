import * as React from 'react'

import { Button, Classes, Collapse } from '@blueprintjs/core'

import DOMPurify from 'dompurify'

interface IProps {
  trace: string | null
}

const ErrorTrace: React.FC<IProps> = ({ trace }) => {
  const [isOpen, setIsOpen] = React.useState(false)
  const toggleOpen = () => setIsOpen(!isOpen)
  return (
    <>
      <Button
        onClick={toggleOpen}
        minimal
        icon={isOpen ? 'chevron-up' : 'chevron-down'}
      >
        {isOpen ? 'Скрыть' : 'Показать'}
      </Button>
      <Collapse className={Classes.RUNNING_TEXT} isOpen={isOpen}>
        <div
          dangerouslySetInnerHTML={{
            __html: DOMPurify.sanitize(trace || ''),
          }}
        />
      </Collapse>
    </>
  )
}

export default ErrorTrace
