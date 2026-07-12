import { Link } from 'react-router-dom'
import type { ReactNode } from 'react'

import gateforgeIcon from '@/assets/brand/gateforge-icon.webp'
import gateforgeLogoHorizontal from '@/assets/brand/gateforge-logo-horizontal.webp'
import gateforgeLogoStacked from '@/assets/brand/gateforge-logo-stacked.webp'
import { cn } from '@/lib/utils'

interface GateForgeBrandProps {
  size?: 'sm' | 'md' | 'lg'
  showTagline?: boolean
  className?: string
  linkTo?: string | null
  variant?: 'gateforge' | 'sovereign'
  layout?: 'stacked' | 'horizontal' | 'mark'
}

const sizes = {
  sm: {
    mark: 'h-9 w-9',
    stacked: 'h-16 w-auto',
    horizontal: 'h-9 w-auto',
    title: 'text-xl',
    tagline: 'text-xs',
  },
  md: {
    mark: 'h-14 w-14',
    stacked: 'h-32 w-auto',
    horizontal: 'h-11 w-auto',
    title: 'text-3xl',
    tagline: 'text-sm',
  },
  lg: {
    mark: 'h-16 w-16',
    stacked: 'h-40 w-auto',
    horizontal: 'h-14 w-auto',
    title: 'text-3xl',
    tagline: 'text-sm',
  },
}

const brandCopy = {
  gateforge: { title: 'GateForge', tagline: 'Identity. Secured. Simplified.' },
  sovereign: { title: 'GateForge', tagline: 'Identity. Secured. Simplified.' },
}

function BrandMarkContent({
  size,
  showTagline,
  className,
  title,
  tagline,
}: {
  size: keyof typeof sizes
  showTagline: boolean
  className?: string
  title: string
  tagline: string
}) {
  const s = sizes[size]
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <img
        src={gateforgeIcon}
        alt=""
        decoding="async"
        aria-hidden
        className={cn(s.mark, 'object-contain')}
      />
      <div className="text-left">
        <p className={cn('font-headline font-extrabold tracking-tight text-on-primary-fixed', s.title)}>
          {title}
        </p>
        {showTagline ? (
          <p className={cn('font-label uppercase tracking-widest text-on-surface-variant', s.tagline)}>
            {tagline}
          </p>
        ) : null}
      </div>
    </div>
  )
}

function BrandStackedContent({
  size,
  showTagline,
  className,
  title,
}: {
  size: keyof typeof sizes
  showTagline: boolean
  className?: string
  title: string
}) {
  const s = sizes[size]
  return (
    <div className={cn('flex flex-col items-center text-center', className)}>
      <img
        src={showTagline ? gateforgeLogoStacked : gateforgeIcon}
        alt={title}
        decoding="async"
        className={cn(showTagline ? s.stacked : s.mark, 'object-contain')}
      />
    </div>
  )
}

export function GateForgeBrand({
  size = 'md',
  showTagline = true,
  className,
  linkTo = '/',
  variant = 'gateforge',
  layout = 'stacked',
}: Readonly<GateForgeBrandProps>) {
  const s = sizes[size]
  const copy = brandCopy[variant]

  let content: ReactNode
  if (layout === 'horizontal') {
    content = (
      <img
        src={gateforgeLogoHorizontal}
        alt={copy.title}
        decoding="async"
        className={cn(s.horizontal, 'object-contain', className)}
      />
    )
  } else if (layout === 'mark') {
    content = (
      <BrandMarkContent
        size={size}
        showTagline={showTagline}
        className={className}
        title={copy.title}
        tagline={copy.tagline}
      />
    )
  } else {
    content = (
      <BrandStackedContent
        size={size}
        showTagline={showTagline}
        className={className}
        title={copy.title}
      />
    )
  }

  if (linkTo) {
    return (
      <Link
        to={linkTo}
        className={cn(
          'inline-flex w-fit shrink-0 items-center',
          layout === 'stacked' && 'mx-auto block',
        )}
      >
        {content}
      </Link>
    )
  }

  return content
}
