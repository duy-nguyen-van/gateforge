import { Navigate, Route, Routes } from 'react-router-dom'

import { AuthLayout } from '@/components/layout/auth-layout'
import { ConsoleLayout } from '@/components/layout/console-layout'
import { LoginForm } from '@/features/login/login-form'
import { MfaChallengeForm } from '@/features/mfa/mfa-challenge-form'
import { RegisterForm } from '@/features/register/register-form'
import { AuditLogsPage } from '@/pages/console/audit-logs-page'
import { ClientsPage } from '@/pages/console/clients-page'
import { ClientDetailPage } from '@/pages/console/client-detail-page'
import { DashboardPage } from '@/pages/console/dashboard-page'
import { IdentityProvidersPage } from '@/pages/console/identity-providers-page'
import { LoginHistoryPage } from '@/pages/console/login-history-page'
import { TenantsPage } from '@/pages/console/tenants-page'
import { TenantDetailPage } from '@/pages/console/tenant-detail-page'
import { UsersPage } from '@/pages/console/users-page'
import { FederationCompletePage } from '@/pages/federation-complete-page'
import { HomePage } from '@/pages/home-page'
import { LogoutPage } from '@/pages/logout-page'
import { ProfilePage } from '@/pages/profile-page'
import { SecurityPage } from '@/pages/security-page'
import { SelectTenantPage } from '@/pages/select-tenant-page'
import { GuestRoute, ProtectedRoute, AdminRoute } from '@/routes/guards'

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />

      <Route element={<GuestRoute />}>
        <Route
          path="/login"
          element={
            <AuthLayout variant="revamp">
              <LoginForm />
            </AuthLayout>
          }
        />
        <Route
          path="/register"
          element={
            <AuthLayout>
              <RegisterForm />
            </AuthLayout>
          }
        />
        <Route
          path="/login/federation/complete"
          element={
            <AuthLayout variant="revamp">
              <FederationCompletePage />
            </AuthLayout>
          }
        />
      </Route>

      <Route
        path="/mfa/challenge"
        element={
          <AuthLayout>
            <MfaChallengeForm />
          </AuthLayout>
        }
      />

      <Route
        path="/select-tenant"
        element={
          <AuthLayout variant="revamp">
            <SelectTenantPage />
          </AuthLayout>
        }
      />

      <Route element={<ProtectedRoute />}>
        <Route path="/logout" element={<LogoutPage />} />
        <Route element={<ConsoleLayout />}>
          <Route path="/settings/profile" element={<ProfilePage />} />
          <Route path="/settings/security" element={<SecurityPage />} />
          <Route element={<AdminRoute />}>
            <Route path="/console" element={<DashboardPage />} />
            <Route path="/console/users" element={<UsersPage />} />
            <Route path="/console/clients" element={<ClientsPage />} />
            <Route path="/console/clients/:clientId" element={<ClientDetailPage />} />
            <Route path="/console/tenants" element={<TenantsPage />} />
            <Route path="/console/tenants/:tenantId" element={<TenantDetailPage />} />
            <Route path="/console/identity-providers" element={<IdentityProvidersPage />} />
            <Route path="/console/audit-logs" element={<AuditLogsPage />} />
            <Route path="/console/login-history" element={<LoginHistoryPage />} />
          </Route>
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
