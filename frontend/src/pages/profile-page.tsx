import { useEffect, useState, type FormEvent } from 'react'

import { updateProfile } from '@/api/client'
import { ApiError } from '@/api/types'
import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { MaterialIcon } from '@/components/icons/material-icon'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/hooks/use-auth'

type ProfileField = {
  label: string
  value: string
  mono?: boolean
}

export function ProfilePage() {
  const { user, refreshProfile } = useAuth()
  const [isEditing, setIsEditing] = useState(false)
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [isSaving, setIsSaving] = useState(false)

  useEffect(() => {
    if (!success) {
      return
    }
    const timeoutId = globalThis.setTimeout(() => setSuccess(null), 5000)
    return () => globalThis.clearTimeout(timeoutId)
  }, [success])

  if (!user) {
    return null
  }

  const profile = user

  const displayName = profile.first_name
    ? `${profile.first_name} ${profile.last_name ?? ''}`.trim()
    : profile.email

  const isDirty =
    firstName.trim() !== (profile.first_name ?? '') || lastName.trim() !== (profile.last_name ?? '')

  function startEdit() {
    setFirstName(profile.first_name ?? '')
    setLastName(profile.last_name ?? '')
    setError(null)
    setSuccess(null)
    setIsEditing(true)
  }

  function cancelEdit() {
    setFirstName(profile.first_name ?? '')
    setLastName(profile.last_name ?? '')
    setError(null)
    setIsEditing(false)
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    const trimmedFirst = firstName.trim()
    if (!trimmedFirst) {
      setError('First name is required.')
      return
    }

    setIsSaving(true)
    try {
      await updateProfile({
        first_name: trimmedFirst,
        last_name: lastName.trim(),
      })
      await refreshProfile()
      setSuccess('Profile updated successfully.')
      setIsEditing(false)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not update profile.')
    } finally {
      setIsSaving(false)
    }
  }

  const personalFields: ProfileField[] = [
    { label: 'First name', value: profile.first_name || '—' },
    { label: 'Last name', value: profile.last_name || '—' },
  ]

  const accountFields: ProfileField[] = [
    { label: 'Email', value: profile.email },
    { label: 'User ID', value: profile.id, mono: true },
    { label: 'Email verified', value: profile.email_verified ? 'Yes' : 'No' },
    { label: 'Member since', value: new Date(profile.created_at).toLocaleDateString() },
  ]

  const displayFields = isEditing ? accountFields : [...personalFields, ...accountFields]

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">Profile</h1>
          <p className="mt-1 text-on-surface-variant">Manage your GateForge identity profile and account details.</p>
        </div>
        <div className="flex shrink-0 flex-wrap gap-3">
          {isEditing ? (
            <>
              <button
                type="button"
                onClick={cancelEdit}
                disabled={isSaving}
                className="flex items-center gap-2 rounded-xl bg-surface-container-high px-5 py-2.5 text-sm font-bold text-on-surface ghost-border transition-opacity hover:opacity-90 disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                form="profile-form"
                disabled={isSaving || !isDirty || !firstName.trim()}
                className="flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-opacity hover:opacity-90 disabled:opacity-50"
              >
                {isSaving ? 'Saving…' : 'Save changes'}
              </button>
            </>
          ) : (
            <button
              type="button"
              onClick={startEdit}
              className="flex items-center gap-2 rounded-xl bg-surface-container-high px-5 py-2.5 text-sm font-bold text-on-surface ghost-border transition-opacity hover:opacity-90"
            >
              <MaterialIcon name="edit" className="text-sm" />
              Edit
            </button>
          )}
        </div>
      </header>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        <div className="flex items-center gap-4 border-b border-surface-container px-8 py-6">
          <DefaultAvatar
            seed={profile.email}
            name={displayName === profile.email ? undefined : displayName}
            size="lg"
            title={profile.email}
          />
          <div>
            <h2 className="font-headline text-xl font-bold">{displayName}</h2>
            <p className="text-sm text-on-surface-variant">{profile.email}</p>
          </div>
        </div>

        {error ? (
          <div className="border-b border-surface-container px-8 py-4">
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </div>
        ) : null}
        {!isEditing && success ? (
          <div className="content-soft-in border-b border-surface-container px-8 py-4">
            <Alert variant="success" className="alert-soft-life">
              <AlertDescription>{success}</AlertDescription>
            </Alert>
          </div>
        ) : null}

        {isEditing ? (
          <div className="content-soft-in border-b border-surface-container px-8 py-6">
            <h3 className="font-headline text-lg font-bold">Personal information</h3>
            <p className="mt-1 text-sm text-on-surface-variant">Update how your name appears across GateForge.</p>
            <form id="profile-form" onSubmit={(e) => void handleSubmit(e)} className="mt-6 grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="profile-first-name">First name</Label>
                <Input
                  id="profile-first-name"
                  required
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  autoComplete="given-name"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="profile-last-name">Last name</Label>
                <Input
                  id="profile-last-name"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  autoComplete="family-name"
                />
              </div>
            </form>
          </div>
        ) : null}

        <dl key={isEditing ? 'editing' : 'view'} className="content-soft-in grid gap-0 sm:grid-cols-2">
          {displayFields.map((field) => (
            <div key={field.label} className="border-b border-surface-container px-8 py-5 sm:border-r">
              <dt className="text-[10px] font-bold uppercase tracking-widest text-on-surface-variant">{field.label}</dt>
              <dd className={`mt-1 font-medium ${field.mono ? 'break-all font-mono text-sm' : ''}`}>{field.value}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  )
}
