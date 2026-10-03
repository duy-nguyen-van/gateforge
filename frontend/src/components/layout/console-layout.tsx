import { useState } from 'react'
import { Outlet, useLocation } from 'react-router'

import { ConsoleBreadcrumb, ConsoleCrumbProvider } from '@/components/layout/console-breadcrumb'
import { ConsoleSidebar } from '@/components/layout/console-sidebar'
import { ConsoleTopbar } from '@/components/layout/console-topbar'
import { isAdminMfaRequired } from '@/features/admin/admin-mfa'
import { AdminMfaEnrollment } from '@/features/admin/admin-mfa-enrollment'
import { useAdminStats } from '@/features/admin/use-admin-queries'

export function ConsoleLayout() {
  const location = useLocation()
  const statsQuery = useAdminStats()
  const needsMfa = isAdminMfaRequired(statsQuery.error)
  const [menuPath, setMenuPath] = useState<string | null>(null)
  const mobileOpen = menuPath === location.pathname
  const closeMenu = () => setMenuPath(null)

  return (
    <ConsoleCrumbProvider>
      <div className="min-h-[100dvh] bg-surface text-on-surface selection:bg-primary-container">
        <a
          href="#console-main"
          className="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-[60] focus:rounded-lg focus:bg-primary focus:px-4 focus:py-2 focus:text-primary-foreground"
        >
          Skip to content
        </a>
        <ConsoleTopbar mobileOpen={mobileOpen} onOpenMenu={() => setMenuPath(location.pathname)} />
        {mobileOpen ? (
          <button
            type="button"
            aria-label="Close menu"
            className="fixed inset-0 z-30 bg-on-surface/20 lg:hidden"
            onClick={closeMenu}
          />
        ) : null}
        <ConsoleSidebar mobileOpen={mobileOpen} onClose={closeMenu} />
        <main id="console-main" className="min-h-[100dvh] px-4 pb-12 pt-24 lg:ml-64 lg:px-10">
          <ConsoleBreadcrumb />
          <div key={needsMfa ? 'admin-mfa' : location.pathname} className="console-page-enter">
            {needsMfa ? <AdminMfaEnrollment /> : <Outlet />}
          </div>
        </main>
      </div>
    </ConsoleCrumbProvider>
  )
}
