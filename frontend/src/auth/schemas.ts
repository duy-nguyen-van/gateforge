import { z } from 'zod'

export const loginSchema = z.object({
  email: z.string().email('Enter a valid email'),
  password: z.string().min(1, 'Password is required'),
  remember_me: z.boolean(),
})

export const registerSchema = z.object({
  email: z.string().email('Enter a valid email'),
  password: z.string().min(12, 'Password must be at least 12 characters').max(128),
  first_name: z.string().max(100).optional(),
  last_name: z.string().max(100).optional(),
})

export const forgotPasswordSchema = z.object({
  email: z.string().email('Enter a valid email'),
})

export const resetPasswordSchema = z.object({
  password: z.string().min(12, 'Password must be at least 12 characters').max(128),
})

export const mfaCodeSchema = z.object({
  code: z.string().min(6, 'Enter your 6-digit code or recovery code'),
})

export type ForgotPasswordValues = z.infer<typeof forgotPasswordSchema>
export type ResetPasswordValues = z.infer<typeof resetPasswordSchema>
export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>
export type MfaCodeFormValues = z.infer<typeof mfaCodeSchema>
