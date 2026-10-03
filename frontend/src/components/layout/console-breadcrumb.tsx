import { useMemo, useState, type ReactNode } from 'react'
import { Link, useLocation } from 'react-router'

import { accountNavItems, consoleNavItems } from '@/components/layout/console-nav'
import { CrumbContext, useCrumbName } from '@/components/layout/use-detail-crumb'

export function ConsoleCrumbProvider({ children }: { children: ReactNode }) {
  const [name, setName] = useState<string | undefined>()
  const value = useMemo(() => ({ name, setName }), [name])
  return <CrumbContext.Provider value={value}>{children}</CrumbContext.Provider>
}

function currentSection(pathname: string) {
  const items = [...accountNavItems, ...consoleNavItems]
  return items
    .filter((item) =>
      item.end ? pathname === item.to : pathname === item.to || pathname.startsWith(`${item.to}/`),
    )
    .sort((a, b) => b.to.length - a.to.length)[0]
}

export function ConsoleBreadcrumb() {
  const location = useLocation()
  const detailName = useCrumbName()
  const section = currentSection(location.pathname)
  if (!section) {
    return null
  }

  const isDetail = location.pathname !== section.to
  const name = isDetail ? detailName : undefined

  return (
    <nav aria-label="Breadcrumb" className="mb-6 text-sm text-on-surface-variant">
      {name ? (
        <ol className="flex flex-wrap items-center gap-2">
          <li>
            <Link to={section.to} className="hover:text-primary">
              {section.label}
            </Link>
          </li>
          <li aria-hidden>/</li>
          <li className="font-medium text-on-surface">{name}</li>
        </ol>
      ) : (
        <span className="font-medium text-on-surface">{section.label}</span>
      )}
    </nav>
  )
}
