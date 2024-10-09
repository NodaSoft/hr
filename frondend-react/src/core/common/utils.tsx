import React from 'react'

import { set } from 'lodash'
import { ValidationError } from 'yup'

export function copyStringToClipboard(str: string): void {
  const el = document.createElement('textarea')
  el.value = str
  if (document.body) {
    document.body.appendChild(el)
  }
  el.select()
  document.execCommand('copy')
  if (document.body) {
    document.body.removeChild(el)
  }
}

type TFormErrors = Record<string, string>

export function getValidateFormErrors(validationErrors: unknown): TFormErrors {
  const errors = {}

  if (validationErrors instanceof ValidationError) {
    for (let i = 0; i < validationErrors.inner.length; i += 1) {
      const error = validationErrors.inner[i]
      set(errors, error.path, error.message)
    }
  }

  return errors
}

export function delay(time: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, time)
  })
}

function escapeRegExpChars(text: string) {
  return text.replace(/([.*+?^=!:${}()|[\]/\\])/g, '\\$1')
}

export function highlightText(text: string, query: string): React.ReactNode[] {
  let lastInd = 0
  const words = query
    .split(/\s+/)
    .filter((word) => word.length > 0)
    .map(escapeRegExpChars)
  if (words.length === 0) {
    return [text]
  }
  const regexp = new RegExp(words.join('|'), 'gi')
  const tokens: React.ReactNode[] = []
  let hasMatches = true
  while (hasMatches) {
    const match = regexp.exec(text)
    if (!match) {
      hasMatches = false
      break
    }
    const { length } = match[0]
    const before = text.slice(lastInd, regexp.lastIndex - length)
    if (before.length > 0) {
      tokens.push(before)
    }
    lastInd = regexp.lastIndex
    tokens.push(<strong key={lastInd}>{match[0]}</strong>)
  }
  const rest = text.slice(lastInd)
  if (rest.length > 0) {
    tokens.push(rest)
  }
  return tokens
}
