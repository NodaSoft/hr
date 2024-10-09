export type TStuff = {
  id: number
  firstName: string
  lastName: string
  isDelete: boolean
  name: string
  email?: string
  mobile?: string
  phone?: string
  photoImageName: string
  type: {
    id: number
    name: string
    comment: string
  }
  sip: string | null
}
