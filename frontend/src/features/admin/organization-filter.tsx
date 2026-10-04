import { useAdminTenants } from '@/features/admin/use-admin-queries'
import { cn } from '@/lib/utils'

export function OrganizationFilter({
  value,
  onChange,
  className,
}: {
  value: string
  onChange: (tenantId: string) => void
  className?: string
}) {
  const tenantsQuery = useAdminTenants({ page: 1, page_size: 100 })
  const tenants = tenantsQuery.data?.data ?? []

  return (
    <select
      aria-label="Organization"
      value={value}
      onChange={(event) => onChange(event.target.value)}
      className={cn(
        'flex h-10 w-full rounded-lg border border-input bg-background px-3 text-sm text-on-surface focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
        className,
      )}
    >
      <option value="">All organizations</option>
      {tenants.map((tenant) => (
        <option key={tenant.id} value={tenant.id}>
          {tenant.name || tenant.domain || tenant.id}
        </option>
      ))}
    </select>
  )
}
