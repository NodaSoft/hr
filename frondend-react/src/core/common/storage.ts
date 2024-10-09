import { get, set, unset } from 'lodash-es'

import logger from './logger'

const prefix = 'testApp'

class AppStorage {
  state: Record<string, string | undefined>

  constructor() {
    let i = 0
    this.state = {}
    while (i < localStorage.length) {
      const storeKey = localStorage.key(i)
      if (typeof storeKey === 'string') {
        if (storeKey.indexOf(prefix) === 0) {
          const value = localStorage.getItem(storeKey)
          set(this.state, storeKey, value)
        }
        i += 1
      }
    }
    window.addEventListener('storage', (e) => {
      if (typeof e.key === 'string') {
        if (e.newValue === null) {
          unset(this.state, e.key)
        } else {
          set(this.state, e.key, e.newValue)
        }
      }
    })
  }

  getItem(path: string) {
    const prefixedPath = `${prefix}.${path}`
    return get(this.state, prefixedPath)
  }

  setItem(path: string, value: unknown) {
    const prefixedPath = `${prefix}.${path}`
    set(this.state, prefixedPath, value)
    try {
      localStorage.setItem(prefixedPath, String(value))
    } catch (e) {
      logger.error(e)
    }
  }

  removeItem(path: string) {
    const prefixedPath = `${prefix}.${path}`
    unset(this.state, prefixedPath)
    try {
      localStorage.removeItem(prefixedPath)
    } catch (e) {
      logger.error(e)
    }
  }
}

export default new AppStorage()
