import { Outlet, useLocation } from 'react-router-dom'

import { ConsoleSidebar } from '@/components/layout/console-sidebar'
import { ConsoleTopbar } from '@/components/layout/console-topbar'

export function ConsoleLayout() {
  const location = useLocation()

  return (
    <div className="min-h-screen bg-surface text-on-surface selection:bg-primary-container">
      <ConsoleTopbar />
      <ConsoleSidebar />
      <main className="ml-64 min-h-screen px-10 pb-12 pt-24">
        <div key={location.pathname} className="console-page-enter">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
