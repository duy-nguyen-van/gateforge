import type { ComponentType } from 'react'

import { MaterialIcon } from '@/components/icons/material-icon'

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
  return <MaterialIcon name="fingerprint" className={className ?? 'text-primary text-2xl'} />
}
