/**
 * API schema types for iam-backend.
 *
 * The backend publishes Swagger 2.x (`docs/swagger.yaml`). openapi-typescript
 * requires OpenAPI 3.x, so types are hand-maintained in `../types.ts` and
 * re-exported here for a stable import path.
 */
export type {
  ApiEnvelope,
  ApiMeta,
  LoginRequest,
  LoginResponse,
  LoginResult,
  MFAChallengeVerifyRequest,
  MFALoginChallengeResponse,
  MFARecoveryCodesResponse,
  MFATOTPSetupResponse,
  MFATOTPVerifyRequest,
  RefreshTokenRequest,
  RegisterRequest,
  UserResponse,
  WebauthnLoginFinishRequest,
  WebauthnLoginStartRequest,
  WebauthnLoginStartResponse,
  WebauthnRegisterFinishRequest,
  WebauthnRegisterStartRequest,
  WebauthnRegisterStartResponse,
} from '../types'

export { ApiError, isMfaChallenge } from '../types'
