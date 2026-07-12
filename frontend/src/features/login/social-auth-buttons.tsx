const SOCIAL = {
  google:
    'https://lh3.googleusercontent.com/aida-public/AB6AXuAzfFGxBeFpV7ZYzpVVIBzryInTvOelJsYFTqoUANCdi3UKDVIYJOXOXKGCGxr6emCpIqDi_P7ggKshPeW3G1NpnfzBYykR1us3D_S_Cd1Mc08zsaOvZczW-Orqyd-K8Mrkl56m4QRyp7UcFMFWJjmpzJt_EPkO_OvP2Snk-__2QFEqILCz6iqtbDAlL0YSw6u-Y_ECyOMOGQcZTV0usRtzPWCM9J5jHchPtoc8hoMxGzQAJsXNyOZPfFvCozRrFEL1qihSF1J_jsk',
  linkedin:
    'https://lh3.googleusercontent.com/aida-public/AB6AXuBSD8yejcfc8AGj9rJBEDKx1MjlFuBYiJrDyv-0ghFmNwEWSkdLFVZKHYbShTeZUviMK1yy4b3dEmJ6xjhOpJMVT4Fn11K02SC3cA7oXaV-waWuSlLsx_bPwo2vVOl6PYHrzPDKCqCeCwd0A9ElmBa4fR18REsTh99U_0KO5zW7XOTTMH-hCofSkFZ-fyO6fx_59_nrDpkCbFnwsG1x4t1bR7AU4Iapz_axpCEd11OZPUp-Uw0zOktmzRoZGe9rKsoatZz91kQF74w',
  facebook:
    'https://lh3.googleusercontent.com/aida-public/AB6AXuBdBKsuwFn-9gud195W5Cce8MXM6rDZsEryg0BAMZCxNGIJ4qjapoE9GAkbKmxstthq0yM4PZN2w5sOCPkBjF9LvCORux_pTgSYxVnI9K_PvVY2spZi9RAjs67QkxE8CcX3ek04P6m1wVMtG5SFoUytoMnn4F7UXc2E36Yky34PLdo1jYfQkBwVbp3oP76FM0l_Gt4OxMFzymeEr1T6J00P-xIsxH5uTSHLeXLLIYb-v1J45Rto3V5XHPOfBVKLOGxRjNSd9cOUyAc',
} as const

export function SocialAuthButtons({
  googleHref,
  googleEnabled,
}: {
  googleHref?: string
  googleEnabled?: boolean
}) {
  return (
    <div className="space-y-5">
      <div className="flex items-center gap-3">
        <div className="flex-grow border-t border-outline-variant/10" />
        <span className="font-label text-[9px] uppercase tracking-tighter text-outline">Third-party verification</span>
        <div className="flex-grow border-t border-outline-variant/10" />
      </div>
      <div className="flex justify-center gap-4">
        {googleEnabled && googleHref ? (
          <a
            href={googleHref}
            className="group rounded-circle bg-surface-container-low p-3 transition-colors hover:bg-surface-variant"
            title="Google Identity"
          >
            <img
              src={SOCIAL.google}
              alt="Google"
              className="h-5 w-5 grayscale opacity-70 transition-all group-hover:grayscale-0 group-hover:opacity-100"
            />
          </a>
        ) : (
          <button
            type="button"
            className="group rounded-circle bg-surface-container-low p-3 transition-colors hover:bg-surface-variant"
            title="Google Identity"
          >
            <img
              src={SOCIAL.google}
              alt="Google"
              className="h-5 w-5 grayscale opacity-70 transition-all group-hover:grayscale-0 group-hover:opacity-100"
            />
          </button>
        )}
        <button
          type="button"
          className="group rounded-circle bg-surface-container-low p-3 transition-colors hover:bg-surface-variant"
          title="LinkedIn Professional"
        >
          <img
            src={SOCIAL.linkedin}
            alt="LinkedIn"
            className="h-5 w-5 grayscale opacity-70 transition-all group-hover:grayscale-0 group-hover:opacity-100"
          />
        </button>
        <button
          type="button"
          className="group rounded-circle bg-surface-container-low p-3 transition-colors hover:bg-surface-variant"
          title="Meta ID"
        >
          <img
            src={SOCIAL.facebook}
            alt="Facebook"
            className="h-5 w-5 grayscale opacity-70 transition-all group-hover:grayscale-0 group-hover:opacity-100"
          />
        </button>
      </div>
    </div>
  )
}
