import * as Yup from 'yup'

export type TPerson = {
  id: number
  name: string
}

export const PersonSchema = Yup.object().shape<TPerson>({
  id: Yup.number(),
  name: Yup.string(),
})
