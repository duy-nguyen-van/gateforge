import { Link } from 'react-router-dom'

import { GateForgeBrand } from '@/components/brand/gateforge-brand'
import { MaterialIcon } from '@/components/icons/material-icon'

const features = [
  { icon: 'fingerprint', title: 'Passkey-First Auth', desc: 'FIDO2 WebAuthn with biometric and hardware key support.' },
  { icon: 'security', title: 'Multi-Factor Security', desc: 'TOTP, recovery codes, and adaptive MFA challenges.' },
  { icon: 'hub', title: 'Enterprise SSO', desc: 'OIDC, SAML federation, and tenant-scoped identity providers.' },
  { icon: 'domain', title: 'Multi-Tenant IAM', desc: 'Enterprise access control with tenant isolation boundaries.' },
]

export function HomePage() {
  return (
    <div className="relative min-h-screen bg-surface text-on-surface">
      <div className="fixed inset-0 -z-10 dot-grid opacity-60" />
      <div className="fixed inset-0 -z-20 bg-gradient-to-br from-surface via-surface-container-low to-primary-container/20" />

      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <GateForgeBrand size="lg" layout="horizontal" showTagline={false} linkTo="/" />
        <div className="flex items-center gap-3">
          <Link
            to="/login"
            className="rounded-lg px-4 py-2 text-sm font-semibold text-on-surface-variant transition-colors hover:text-primary"
          >
            Sign in
          </Link>
          <Link
            to="/register"
            className="rounded-lg bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-all hover:bg-primary-dim"
          >
            Request Access
          </Link>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6 pb-20 pt-12">
        <section className="mb-20 text-center">
          <h1 className="font-headline text-5xl font-extrabold tracking-tight text-on-primary-fixed md:text-6xl">
            Identity & Access,
            <br />
            <span className="text-primary">Built for Enterprise</span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg text-on-surface-variant">
            Secure sign-in, multi-factor authentication, and passkey support — built for enterprise-grade access control.
          </p>
          <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Link
              to="/login"
              className="flex items-center gap-2 rounded-lg bg-primary px-8 py-4 font-headline font-bold text-on-primary shadow-xl shadow-primary/25 transition-all hover:bg-primary-dim"
            >
              <MaterialIcon name="login" filled />
              Sign in to Console
            </Link>
            <Link
              to="/register"
              className="flex items-center gap-2 rounded-lg border border-outline-variant/30 bg-surface-container-lowest px-8 py-4 font-headline font-bold text-on-surface transition-all hover:bg-surface-container-low"
            >
              <MaterialIcon name="person_add" />
              Create Account
            </Link>
          </div>
        </section>

        <section className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {features.map((f) => (
            <div key={f.title} className="rounded-xl bg-surface-container-lowest p-6 shadow-sm ghost-border transition-shadow hover:shadow-md">
              <MaterialIcon name={f.icon} className="mb-4 rounded-lg bg-primary-container p-2 text-primary text-2xl" />
              <h3 className="font-headline text-lg font-bold">{f.title}</h3>
              <p className="mt-2 text-sm text-on-surface-variant">{f.desc}</p>
            </div>
          ))}
        </section>

        <section className="mt-16 flex flex-wrap items-center justify-center gap-8 opacity-60">
          {['FIDO2 Certified', 'AES-256 Encryption', 'SOC2 Type II', 'GDPR Compliant'].map((badge) => (
            <div key={badge} className="flex items-center gap-2">
              <MaterialIcon name="verified" className="text-sm" />
              <span className="font-label text-[10px] font-bold uppercase tracking-widest">{badge}</span>
            </div>
          ))}
        </section>
      </main>
    </div>
  )
}
