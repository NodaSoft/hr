import React from 'react'

import { TApiError } from 'core/models/api'

type ErrorContent = (apiError: TApiError) => React.ReactNode
export type KnownApiErrors = Record<number, ErrorContent>

export const knownApiErrors: KnownApiErrors = {}
