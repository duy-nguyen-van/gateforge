import { createContext, useContext, useEffect, type Dispatch, type SetStateAction } from 'react'

type CrumbContextValue = {
  name?: string
  setName: Dispatch<SetStateAction<string | undefined>>
}

export const CrumbContext = createContext<CrumbContextValue | null>(null)

export function useCrumbName() {
  return useContext(CrumbContext)?.name
}

export function useDetailCrumb(name: string | undefined) {
  const setName = useContext(CrumbContext)?.setName

  useEffect(() => {
    if (!setName) {
      return
    }
    setName(name)
    return () => setName(undefined)
  }, [name, setName])
}
