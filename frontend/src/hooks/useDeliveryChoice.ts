import { useCallback, useEffect, useState } from 'react'
import { useWatch, type UseFormReturn } from 'react-hook-form'
import { api, type SavedDetails, type User } from '../api'
import type { CheckoutValues } from '../validation'

export type Delivery = Omit<CheckoutValues, 'email'>

export const EMPTY_DELIVERY: Delivery = { phone: '', label: 'Home', line1: '', line2: '', city: '', state: '', pincode: '' }

export type Choice = number | 'new' // index into the saved list, or the user's own new address

const toDelivery = ({ phone, address }: SavedDetails): Delivery => ({ phone, ...address })

type Form = Pick<UseFormReturn<CheckoutValues>, 'control' | 'getValues' | 'setValue'>

/**
 * Which phone + address the logged-in user is checking out with: one of their saved ones
 * (from previous orders) or a new one they type. Keeps the form fields in sync with the choice.
 *
 * - On login, the most recent saved address is prefilled, unless the user already typed one:
 *   then their input stays selected as "Your new address".
 * - Switching to a saved address keeps the typed one as a draft and restores it on switching back.
 */
export function useDeliveryChoice(user: User | null, { control, getValues, setValue }: Form) {
  const [saved, setSaved] = useState<{ userId: number; list: SavedDetails[] } | null>(null)
  const [choice, setChoice] = useState<Choice>('new')
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState<Delivery | null>(null)
  const [reloadKey, setReloadKey] = useState(0)
  const typedLine1 = useWatch({ control, name: 'line1' })

  const savedList = user && saved?.userId === user.id ? saved.list : []

  const fill = useCallback(
    (d: Delivery) => {
      for (const [field, value] of Object.entries(d) as [keyof Delivery, string][]) {
        setValue(field, value, { shouldDirty: true })
      }
    },
    [setValue],
  )

  useEffect(() => {
    if (!user) return
    let cancelled = false
    api
      .savedDetails()
      .then(({ saved: list }) => {
        if (cancelled) return
        setSaved({ userId: user.id, list })
        const { phone, line1 } = getValues()
        const userAlreadyTyped = phone.trim() !== '' || line1.trim() !== ''
        if (list.length > 0 && !userAlreadyTyped) {
          fill(toDelivery(list[0]))
          setChoice(0)
        }
      })
      .catch(() => {}) // not critical: the user can still fill the form by hand
    return () => {
      cancelled = true
    }
  }, [user, reloadKey, getValues, fill])

  const choose = (next: Choice) => {
    if (next === choice && !editing) return
    if (choice === 'new') {
      const { email: _email, ...typed } = getValues()
      setDraft(typed)
    }
    setChoice(next)
    setEditing(false)
    fill(next === 'new' ? (draft ?? EMPTY_DELIVERY) : toDelivery(savedList[next]))
  }

  const clear = () => {
    setChoice('new')
    setEditing(false)
    setDraft(null)
  }

  /** Start over after an order; reload, since the address just used is now a saved one. */
  const reload = () => {
    clear()
    setReloadKey((n) => n + 1)
  }

  const usingSaved = savedList.length > 0 && choice !== 'new'
  const draftPreview = draft?.line1.trim() ? `${draft.line1}, ${draft.city}` : undefined

  return {
    savedList,
    choice,
    editing,
    /** With a saved address selected its card shows everything, so the fields stay hidden unless editing. */
    showFields: !usingSaved || editing,
    newAddress:
      choice === 'new' ? { typed: typedLine1.trim() !== '' } : { typed: draftPreview !== undefined, preview: draftPreview },
    choose,
    edit: () => setEditing(true),
    clear,
    reload,
  }
}
