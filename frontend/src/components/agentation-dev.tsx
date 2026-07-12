import { Agentation } from 'agentation'

const endpoint = import.meta.env.VITE_AGENTATION_ENDPOINT ?? 'http://localhost:4747'

/** Dev-only visual feedback for AI agents. Excluded from production builds. */
export function AgentationDev() {
  if (!import.meta.env.DEV) {
    return null
  }

  return (
    <Agentation
      endpoint={endpoint}
      onSessionCreated={(sessionId) => {
        console.info('[agentation] session:', sessionId)
      }}
    />
  )
}
