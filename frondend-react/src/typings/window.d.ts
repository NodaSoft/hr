import { YM } from 'core/models'
import { GlobalAppData } from 'core/models/app'
import { GlobalMessengerAppData } from 'core/models/messengerApp'

declare global {
  interface Window {
    [process.env.APP_NAME]: GlobalAppData
  }
}
export {}
