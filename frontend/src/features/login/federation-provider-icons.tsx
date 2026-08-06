import type { ComponentType } from 'react'
import { FingerprintIcon } from 'lucide-react'

import { GoogleIcon } from './social-provider-icons'

const FEDERATION_ICONS: Record<string, ComponentType<{ className?: string }>> = {
  google: GoogleIcon,
}

interface FederationProviderIconProps {
  provider: string
  className?: string
}

export function FederationProviderIcon({ provider, className }: FederationProviderIconProps) {
  const Icon = FEDERATION_ICONS[provider]
  if (Icon) {
    return <Icon className={className} />
  }
  return <FingerprintIcon className={className ?? 'h-6 w-6 text-primary'} aria-hidden />
}
