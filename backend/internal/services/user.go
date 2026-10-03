package services

import (
	"context"
	stderrors "errors"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

// UserService covers registration, login, refresh, profile, and tenant-scoped tokens.
type UserService interface {
	Register(ctx context.Context, req *dtos.RegisterRequest, host string) (*models.User, error)
	AuthenticateUser(ctx context.Context, req *dtos.LoginRequest) (*models.User, error)
	CompleteAuth(ctx context.Context, u *models.User, in TenantResolveInput) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error)
	SelectTenant(ctx context.Context, req *dtos.TenantSelectRequest) (*dtos.LoginResponse, error)
	SwitchTenant(ctx context.Context, userID, tenantID string) (*dtos.LoginResponse, error)
	IssueTokensForUser(ctx context.Context, u *models.User, tenantID string) (*dtos.LoginResponse, error)
	Login(ctx context.Context, req *dtos.LoginRequest, host string) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error)
	Refresh(ctx context.Context, req *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetOneByID(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, userID string, req *dtos.UpdateProfileRequest) (*models.User, error)
	ListMemberships(ctx context.Context, userID string) ([]dtos.TenantSummary, error)
	ListMembershipsPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error)
	RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error
}

type userService struct {
	userRepo         repositories.UserRepository
	membershipRepo   repositories.TenantMembershipRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	tenantCtx        TenantContextService
	cfg              *config.Config
	tokenService     *auth.TokenService
	audit            AuditService
	loginCache       cache.Cache
	sessions         SessionService
	resetMailer      passwordResetMailer
}

// ProvideUserService creates the identity service.
func ProvideUserService(
	userRepo repositories.UserRepository,
	membershipRepo repositories.TenantMembershipRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	tenantCtx TenantContextService,
	cfg *config.Config,
	tokenService *auth.TokenService,
	audit AuditService,
	loginCache cache.Cache,
	sessions SessionService,
	resetMailer passwordResetMailer,
) UserService {
	return &userService{
		userRepo:         userRepo,
		membershipRepo:   membershipRepo,
		refreshTokenRepo: refreshTokenRepo,
		tenantCtx:        tenantCtx,
		cfg:              cfg,
		tokenService:     tokenService,
		audit:            audit,
		loginCache:       loginCache,
		sessions:         sessions,
		resetMailer:      resetMailer,
	}
}

func (s *userService) Register(ctx context.Context, req *dtos.RegisterRequest, host string) (*models.User, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.Register", nil,
		func(ctx context.Context) (*models.User, error) {
			return s.register(ctx, req, host)
		})
}

func (s *userService) register(ctx context.Context, req *dtos.RegisterRequest, host string) (*models.User, error) {
	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	emailLower := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err == nil {
		return nil, errors.CodedConflict(constants.EmailAlreadyRegistered, nil).
			WithOperation("register").
			WithResource("user")
	}
	var appErr *errors.AppError
	if !stderrors.As(err, &appErr) || appErr.Type != errors.ErrorTypeNotFound {
		return nil, err
	}

	resolved, err := s.tenantCtx.Resolve(ctx, TenantResolveInput{
		Host:          host,
		TenantIDParam: req.TenantID,
	})
	if err != nil {
		return nil, err
	}
	if resolved.RequiresPicker || resolved.TenantID == "" {
		resolved.TenantID = s.cfg.DefaultTenantID
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), constants.BcryptCost)
	if err != nil {
		return nil, errors.InternalError("Failed to hash password", err).
			WithOperation("register").
			WithResource("user")
	}

	u := &models.User{
		BaseModel:     models.NewBaseModel(),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         strings.TrimSpace(req.Email),
		EmailVerified: false,
		Status:        constants.UserStatusActive,
	}

	if err := s.userRepo.CreateWithPasswordHash(ctx, u, string(hash)); err != nil {
		return nil, err
	}

	membership := &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    u.ID,
		TenantID:  resolved.TenantID,
		Role:      constants.TenantMembershipRoleMember,
		Status:    constants.TenantMembershipStatusActive,
	}
	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAuthRegister,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		TenantID:     resolved.TenantID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   u.ID,
	})
	return u, nil
}

func (s *userService) Login(ctx context.Context, req *dtos.LoginRequest, host string) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	return monitoring.Observe2(ctx, userTracer, "UserService.Login", nil,
		func(ctx context.Context) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
			u, err := s.AuthenticateUser(ctx, req)
			if err != nil {
				return nil, nil, err
			}
			return s.CompleteAuth(ctx, u, TenantResolveInput{
				Host:          host,
				TenantIDParam: req.TenantID,
				UserID:        u.ID,
			})
		})
}

func (s *userService) AuthenticateUser(ctx context.Context, req *dtos.LoginRequest) (*models.User, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.AuthenticateUser", nil,
		func(ctx context.Context) (*models.User, error) {
			return s.authenticateUser(ctx, req)
		})
}

func (s *userService) authenticateUser(ctx context.Context, req *dtos.LoginRequest) (*models.User, error) {
	emailLower := strings.ToLower(strings.TrimSpace(req.Email))
	if s.loginLocked(ctx, emailLower) {
		monitoring.AddAuthLogin(ctx, "locked")
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionAuthLoginLocked,
			Result:    constants.AuditResultDenied,
			ActorType: constants.AuditActorTypeUser,
		})
		return nil, errors.CodedUnauthorized(constants.InvalidLoginCredentials, nil).
			WithOperation("login").
			WithResource("user")
	}

	u, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err != nil {
		var appErr *errors.AppError
		if stderrors.As(err, &appErr) && appErr.Type == errors.ErrorTypeNotFound {
			s.audit.Record(ctx, domains.AuditRecordParams{
				Action:    constants.AuditActionAuthLogin,
				Result:    constants.AuditResultFailure,
				ActorType: constants.AuditActorTypeUser,
			})
			return nil, errors.CodedUnauthorized(constants.InvalidLoginCredentials, nil).
				WithOperation("login").
				WithResource("user")
		}
		return nil, err
	}
	if u.PasswordCredential == nil {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionAuthLogin,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
			ActorID:   u.ID,
		})
		return nil, errors.CodedUnauthorized(constants.InvalidLoginCredentials, nil).
			WithOperation("login").
			WithResource("user")
	}
	if u.Status != constants.UserStatusActive {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionAuthLogin,
			Result:    constants.AuditResultDenied,
			ActorType: constants.AuditActorTypeUser,
			ActorID:   u.ID,
		})
		return nil, errors.CodedUnauthorized(constants.AccountNotActive, nil).
			WithOperation("login").
			WithResource("user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordCredential.PasswordHash), []byte(req.Password)); err != nil {
		s.noteLoginFailure(ctx, emailLower, u.ID)
		monitoring.AddAuthLogin(ctx, "failure")
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionAuthLogin,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
			ActorID:   u.ID,
		})
		return nil, errors.CodedUnauthorized(constants.InvalidLoginCredentials, nil).
			WithOperation("login").
			WithResource("user")
	}

	s.clearLoginLock(ctx, emailLower)
	monitoring.AddAuthLogin(ctx, "success")
	return u, nil
}

func (s *userService) CompleteAuth(ctx context.Context, u *models.User, in TenantResolveInput) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	return monitoring.Observe2(ctx, userTracer, "UserService.CompleteAuth",
		[]attribute.KeyValue{attribute.String("user_id", u.ID)},
		func(ctx context.Context) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
			return s.completeAuth(ctx, u, in)
		})
}

func (s *userService) completeAuth(ctx context.Context, u *models.User, in TenantResolveInput) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	in.UserID = u.ID
	resolved, err := s.tenantCtx.Resolve(ctx, in)
	if err != nil {
		return nil, nil, err
	}
	if resolved.RequiresPicker {
		sel, err := s.buildTenantSelection(ctx, u.ID)
		if err != nil {
			return nil, nil, err
		}
		return nil, sel, nil
	}
	out, err := s.IssueTokensForUser(ctx, u, resolved.TenantID)
	if err != nil {
		return nil, nil, err
	}
	return out, nil, nil
}

func (s *userService) buildTenantSelection(ctx context.Context, userID string) (*dtos.TenantSelectionResponse, error) {
	summaries, err := s.ListMemberships(ctx, userID)
	if err != nil {
		return nil, err
	}
	token, exp, err := s.tokenService.SignSelectionToken(userID)
	if err != nil {
		return nil, errors.InternalError("Failed to issue selection token", err).
			WithOperation("tenant_selection").
			WithResource("user")
	}
	return &dtos.TenantSelectionResponse{
		SelectionRequired: true,
		Tenants:           summaries,
		SelectionToken:    token,
		ExpiresIn:         exp,
	}, nil
}

func (s *userService) SelectTenant(ctx context.Context, req *dtos.TenantSelectRequest) (*dtos.LoginResponse, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.SelectTenant", nil,
		func(ctx context.Context) (*dtos.LoginResponse, error) {
			return s.selectTenant(ctx, req)
		})
}

func (s *userService) selectTenant(ctx context.Context, req *dtos.TenantSelectRequest) (*dtos.LoginResponse, error) {
	userID, err := s.tokenService.ParseSelectionToken(req.SelectionToken)
	if err != nil {
		return nil, errors.CodedUnauthorized(constants.InvalidSelectionToken, err).
			WithOperation("tenant_select").
			WithResource("user")
	}
	if err := s.tenantCtx.ValidateMembership(ctx, userID, req.TenantID); err != nil {
		return nil, err
	}
	u, err := s.userRepo.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out, err := s.IssueTokensForUser(ctx, u, req.TenantID)
	if err != nil {
		return nil, err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionTenantSelect,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		TenantID:     req.TenantID,
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   req.TenantID,
	})
	return out, nil
}

func (s *userService) SwitchTenant(ctx context.Context, userID, tenantID string) (*dtos.LoginResponse, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.SwitchTenant",
		[]attribute.KeyValue{
			attribute.String("user_id", userID),
			attribute.String("tenant_id", tenantID),
		},
		func(ctx context.Context) (*dtos.LoginResponse, error) {
			return s.switchTenant(ctx, userID, tenantID)
		})
}

func (s *userService) switchTenant(ctx context.Context, userID, tenantID string) (*dtos.LoginResponse, error) {
	if err := s.tenantCtx.ValidateMembership(ctx, userID, tenantID); err != nil {
		return nil, err
	}
	u, err := s.userRepo.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out, err := s.IssueTokensForUser(ctx, u, tenantID)
	if err != nil {
		return nil, err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionTenantSwitch,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeTenant,
		ResourceID:   tenantID,
	})
	return out, nil
}

func (s *userService) IssueTokensForUser(ctx context.Context, u *models.User, tenantID string) (*dtos.LoginResponse, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.IssueTokensForUser",
		[]attribute.KeyValue{
			attribute.String("user_id", u.ID),
			attribute.String("tenant_id", tenantID),
		},
		func(ctx context.Context) (*dtos.LoginResponse, error) {
			return s.issueTokensForUser(ctx, u, tenantID)
		})
}

func (s *userService) issueTokensForUser(ctx context.Context, u *models.User, tenantID string) (*dtos.LoginResponse, error) {
	return s.issueAccessAndRefresh(ctx, u, tenantID)
}

func (s *userService) Refresh(ctx context.Context, req *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.Refresh", nil,
		func(ctx context.Context) (*dtos.LoginResponse, error) {
			return s.refresh(ctx, req)
		})
}

func (s *userService) refresh(ctx context.Context, req *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error) {
	raw := strings.TrimSpace(req.RefreshToken)
	if raw == "" {
		return nil, errors.CodedUnauthorized(constants.InvalidRefreshToken, nil).
			WithOperation("refresh").
			WithResource("refresh_token")
	}

	hash := auth.HashOpaqueToken(raw)
	opaqueRefreshToken, refreshTokenHash, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return nil, errors.InternalError("Failed to issue refresh token", err).
			WithOperation("refresh").
			WithResource("refresh_token")
	}
	newRT := &models.RefreshToken{
		OAuthClientID: s.cfg.NativeOAuthClientID,
		TokenHash:     refreshTokenHash,
		Revoked:       false,
		ExpiresAt:     time.Now().UTC().Add(s.cfg.JWTRefreshTTL),
	}
	status, err := s.refreshTokenRepo.Rotate(ctx, hash, newRT)
	if err != nil {
		return nil, err
	}
	if status == repositories.RefreshRotationReuse {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionAuthRefresh,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
			NewValue:  map[string]any{"reuse": true},
		})
		return nil, errors.CodedUnauthorized(constants.InvalidRefreshToken, nil).
			WithOperation("refresh").
			WithResource("refresh_token")
	}
	if status != repositories.RefreshRotationOK {
		return nil, errors.CodedUnauthorized(constants.InvalidRefreshToken, nil).
			WithOperation("refresh").
			WithResource("refresh_token")
	}

	u, err := s.userRepo.GetOneByID(ctx, newRT.UserID)
	if err != nil {
		return nil, err
	}
	if u.Status != constants.UserStatusActive {
		return nil, errors.CodedUnauthorized(constants.AccountNotActive, nil).
			WithOperation("refresh").
			WithResource("user")
	}

	tenantID := newRT.TenantID
	access, _, err := s.tokenService.SignAccessToken(u.ID, tenantID)
	if err != nil {
		return nil, errors.InternalError("Failed to issue token", err).
			WithOperation("refresh").
			WithResource("user")
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAuthRefresh,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   u.ID,
	})

	return &dtos.LoginResponse{
		AccessToken:      access,
		RefreshToken:     opaqueRefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(s.cfg.JWTAccessTTL.Seconds()),
		RefreshExpiresIn: int64(s.cfg.JWTRefreshTTL.Seconds()),
		ActiveTenantID:   tenantID,
	}, nil
}

func (s *userService) issueAccessAndRefresh(ctx context.Context, u *models.User, tenantID string) (*dtos.LoginResponse, error) {
	access, _, err := s.tokenService.SignAccessToken(u.ID, tenantID)
	if err != nil {
		return nil, errors.InternalError("Failed to issue token", err).
			WithOperation("login").
			WithResource("user")
	}

	opaqueRefreshToken, refreshTokenHash, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return nil, errors.InternalError("Failed to issue refresh token", err).
			WithOperation("login").
			WithResource("refresh_token")
	}

	rt := &models.RefreshToken{
		TenantID:      tenantID,
		UserID:        u.ID,
		OAuthClientID: s.cfg.NativeOAuthClientID,
		TokenHash:     refreshTokenHash,
		Revoked:       false,
		ExpiresAt:     time.Now().UTC().Add(s.cfg.JWTRefreshTTL),
	}

	if err := s.refreshTokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAuthLogin,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   u.ID,
	})

	return &dtos.LoginResponse{
		AccessToken:      access,
		RefreshToken:     opaqueRefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(s.cfg.JWTAccessTTL.Seconds()),
		RefreshExpiresIn: int64(s.cfg.JWTRefreshTTL.Seconds()),
		ActiveTenantID:   tenantID,
	}, nil
}

func (s *userService) GetOneByID(ctx context.Context, userID string) (*models.User, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.GetOneByID",
		[]attribute.KeyValue{attribute.String("user_id", userID)},
		func(ctx context.Context) (*models.User, error) {
			return s.userRepo.GetOneByID(ctx, userID)
		})
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req *dtos.UpdateProfileRequest) (*models.User, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.UpdateProfile",
		[]attribute.KeyValue{attribute.String("user_id", userID)},
		func(ctx context.Context) (*models.User, error) {
			return s.updateProfile(ctx, userID, req)
		})
}

func (s *userService) updateProfile(ctx context.Context, userID string, req *dtos.UpdateProfileRequest) (*models.User, error) {
	if req == nil || (req.FirstName == nil && req.LastName == nil) {
		return nil, errors.ValidationError("At least one field must be provided", nil)
	}

	existing, err := s.userRepo.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing.Status != constants.UserStatusActive {
		return nil, errors.CodedValidation(constants.AccountNotActive, nil)
	}

	patch := repositories.UserProfilePatch{}
	if req.FirstName != nil {
		name := strings.TrimSpace(*req.FirstName)
		if name == "" {
			return nil, errors.ValidationError("First name cannot be empty", nil)
		}
		patch.FirstName = &name
	}
	if req.LastName != nil {
		lastName := strings.TrimSpace(*req.LastName)
		patch.LastName = &lastName
	}

	return s.userRepo.UpdateProfile(ctx, userID, patch)
}

func (s *userService) ListMemberships(ctx context.Context, userID string) ([]dtos.TenantSummary, error) {
	return monitoring.Observe(ctx, userTracer, "UserService.ListMemberships",
		[]attribute.KeyValue{attribute.String("user_id", userID)},
		func(ctx context.Context) ([]dtos.TenantSummary, error) {
			return s.listMemberships(ctx, userID)
		})
}

func (s *userService) listMemberships(ctx context.Context, userID string) ([]dtos.TenantSummary, error) {
	rows, err := s.membershipRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tenantSummariesFromMemberships(rows), nil
}

func (s *userService) ListMembershipsPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error) {
	return monitoring.Observe2(ctx, userTracer, "UserService.ListMembershipsPaginated",
		[]attribute.KeyValue{attribute.String("user_id", userID)},
		func(ctx context.Context) ([]dtos.TenantSummary, *dtos.Pageable, error) {
			return s.listMembershipsPaginated(ctx, userID, pr)
		})
}

func (s *userService) listMembershipsPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error) {
	result, err := s.membershipRepo.ListByUserIDPaginated(ctx, userID, pr)
	if err != nil {
		return nil, nil, err
	}
	return tenantSummariesFromMemberships(result.Data), result.Pageable, nil
}

func tenantSummariesFromMemberships(rows []models.TenantMembership) []dtos.TenantSummary {
	out := make([]dtos.TenantSummary, 0, len(rows))
	for _, m := range rows {
		summary := dtos.TenantSummary{
			ID:   m.TenantID,
			Role: string(m.Role),
		}
		if m.Tenant != nil {
			summary.Name = m.Tenant.Name
			summary.Domain = m.Tenant.Domain
		}
		out = append(out, summary)
	}
	return out
}

func (s *userService) RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error {
	return monitoring.ObserveErr(ctx, userTracer, "UserService.RevokeAllRefreshTokensForUser",
		[]attribute.KeyValue{attribute.String("user_id", userID)},
		func(ctx context.Context) error {
			if userID == "" {
				return nil
			}
			return s.refreshTokenRepo.RevokeAllValidForUser(ctx, userID)
		})
}
