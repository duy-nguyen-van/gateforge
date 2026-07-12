import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'

import { buttonVariants, type ButtonVariantProps } from '@/components/ui/button-variants'
import { cn } from '@/lib/utils'

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<'button'> &
  ButtonVariantProps & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : 'button'
  return <Comp className={cn(buttonVariants({ variant, size, className }))} {...props} />
}

export { Button }
