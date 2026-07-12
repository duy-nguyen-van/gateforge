package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type memCache struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemCache() *memCache {
	return &memCache{data: map[string]string{}}
}

func (c *memCache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	if !ok {
		return "", errors.NotFoundError("cache key", nil)
	}
	return v, nil
}

func (c *memCache) Set(_ context.Context, key, value string, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *memCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *memCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.data[key]
	return ok, nil
}

func (c *memCache) Close() error { return nil }

type auditCapture struct {
	params []domains.AuditRecordParams
}

func (a *auditCapture) Record(_ context.Context, p domains.AuditRecordParams) {
	a.params = append(a.params, p)
}

type userTestRepo struct {
	users    map[string]*models.User
	byEmail  map[string]*models.User
	created  []*models.User
	password map[string]string
}

func newUserTestRepo() *userTestRepo {
	return &userTestRepo{
		users:    map[string]*models.User{},
		byEmail:  map[string]*models.User{},
		password: map[string]string{},
	}
}

func (r *userTestRepo) seed(email, password string) *models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	u := &models.User{
		BaseModel: models.NewBaseModel(),
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Status:    constants.UserStatusActive,
	}
	u.EmailLower = email
	u.PasswordCredential = &models.PasswordCredential{PasswordHash: string(hash)}
	r.users[u.ID] = u
	r.byEmail[email] = u
	r.password[u.ID] = string(hash)
	return u
}

func (r *userTestRepo) CreateWithPasswordHash(_ context.Context, user *models.User, passwordHash string) error {
	user.EmailLower = user.Email
	r.users[user.ID] = user
	r.byEmail[user.EmailLower] = user
	r.password[user.ID] = passwordHash
	user.PasswordCredential = &models.PasswordCredential{PasswordHash: passwordHash}
	r.created = append(r.created, user)
	return nil
}

func (r *userTestRepo) CreateUserOnly(_ context.Context, user *models.User) error {
	user.EmailLower = user.Email
	r.users[user.ID] = user
	r.byEmail[user.EmailLower] = user
	return nil
}

func (r *userTestRepo) GetOneByID(_ context.Context, id string) (*models.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}

func (r *userTestRepo) GetByEmailLower(_ context.Context, emailLower string) (*models.User, error) {
	u, ok := r.byEmail[emailLower]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}

func (r *userTestRepo) Count(_ context.Context) (int64, error) {
	return int64(len(r.users)), nil
}

func (r *userTestRepo) CountPlatformAdmins(_ context.Context) (int64, error) {
	var n int64
	for _, u := range r.users {
		if u.IsPlatformAdmin {
			n++
		}
	}
	return n, nil
}

func (r *userTestRepo) SetPlatformAdmin(_ context.Context, userID string, isAdmin bool) error {
	u, ok := r.users[userID]
	if !ok {
		return errors.NotFoundError("User", nil)
	}
	u.IsPlatformAdmin = isAdmin
	return nil
}

func (r *userTestRepo) UpdateStatus(_ context.Context, userID string, status constants.UserStatus) error {
	u, ok := r.users[userID]
	if !ok {
		return errors.NotFoundError("User", nil)
	}
	u.Status = status
	return nil
}

func (r *userTestRepo) UpdateProfile(_ context.Context, userID string, patch repositories.UserProfilePatch) (*models.User, error) {
	u, ok := r.users[userID]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	if patch.FirstName != nil {
		u.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		u.LastName = *patch.LastName
	}
	return u, nil
}

func (r *userTestRepo) List(_ context.Context, tenantID, search string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	rows := make([]models.User, 0, len(r.users))
	for _, u := range r.users {
		rows = append(rows, *u)
	}
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.User]{Data: page, Pageable: pageable}, nil
}

type refreshTokenTestRepo struct {
	byHash map[string]*models.RefreshToken
	byID   map[string]*models.RefreshToken
	revoked []string
}

func newRefreshTokenTestRepo() *refreshTokenTestRepo {
	return &refreshTokenTestRepo{
		byHash: map[string]*models.RefreshToken{},
		byID:   map[string]*models.RefreshToken{},
	}
}

func (r *refreshTokenTestRepo) Create(_ context.Context, rt *models.RefreshToken) error {
	if rt.ID == "" {
		rt.ID = uuid.Must(uuid.NewV7()).String()
	}
	r.byHash[rt.TokenHash] = rt
	r.byID[rt.ID] = rt
	return nil
}

func (r *refreshTokenTestRepo) FindValidByTokenHash(_ context.Context, tokenHash string) (*models.RefreshToken, error) {
	rt, ok := r.byHash[tokenHash]
	if !ok || rt.Revoked {
		return nil, errors.NotFoundError("Refresh token", nil)
	}
	return rt, nil
}

func (r *refreshTokenTestRepo) RevokeByID(_ context.Context, id string) error {
	if rt, ok := r.byID[id]; ok {
		rt.Revoked = true
	}
	return nil
}

func (r *refreshTokenTestRepo) RevokeAllValidForUser(_ context.Context, userID string) error {
	r.revoked = append(r.revoked, userID)
	for _, rt := range r.byID {
		if rt.UserID == userID {
			rt.Revoked = true
		}
	}
	return nil
}

func (r *refreshTokenTestRepo) RevokeAndCreate(_ context.Context, oldID string, newRT *models.RefreshToken) error {
	if err := r.RevokeByID(context.Background(), oldID); err != nil {
		return err
	}
	return r.Create(context.Background(), newRT)
}

func (r *refreshTokenTestRepo) UsageByClientRecordID(_ context.Context, clientRecordID string) (*repositories.RefreshTokenUsage, error) {
	return &repositories.RefreshTokenUsage{TotalIssued: 1, ActiveCount: 1}, nil
}

type sessionTestRepo struct {
	sessions map[string]*models.Session
	deleted  []string
}

func newSessionTestRepo() *sessionTestRepo {
	return &sessionTestRepo{sessions: map[string]*models.Session{}}
}

func (r *sessionTestRepo) Create(_ context.Context, s *models.Session) error {
	if s.ID == "" {
		s.ID = models.NewBaseModel().ID
	}
	r.sessions[s.ID] = s
	return nil
}

func (r *sessionTestRepo) GetValidByID(_ context.Context, id string) (*models.Session, error) {
	s, ok := r.sessions[id]
	if !ok {
		return nil, errors.NotFoundError("Session", nil)
	}
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.NotFoundError("Session", nil)
	}
	return s, nil
}

func (r *sessionTestRepo) DeleteByID(_ context.Context, id string) error {
	r.deleted = append(r.deleted, id)
	delete(r.sessions, id)
	return nil
}

func (r *sessionTestRepo) DeleteAllByUserID(_ context.Context, userID string) error {
	for id, s := range r.sessions {
		if s.UserID == userID {
			delete(r.sessions, id)
		}
	}
	return nil
}

func (r *sessionTestRepo) CountActive(_ context.Context) (int64, error) {
	return int64(len(r.sessions)), nil
}

func (r *sessionTestRepo) CountActiveByUserID(_ context.Context, userID string) (int64, error) {
	var n int64
	for _, s := range r.sessions {
		if s.UserID == userID {
			n++
		}
	}
	return n, nil
}

type authCodeTestRepo struct {
	byCode map[string]*models.AuthorizationCode
}

func newAuthCodeTestRepo() *authCodeTestRepo {
	return &authCodeTestRepo{byCode: map[string]*models.AuthorizationCode{}}
}

func (r *authCodeTestRepo) Create(_ context.Context, row *models.AuthorizationCode) error {
	r.byCode[row.Code] = row
	return nil
}

func (r *authCodeTestRepo) TakeByCode(_ context.Context, code string) (*models.AuthorizationCode, error) {
	row, ok := r.byCode[code]
	if !ok {
		return nil, errors.NotFoundError("Authorization code", nil)
	}
	return row, nil
}

func (r *authCodeTestRepo) DeleteByCode(_ context.Context, code string) error {
	delete(r.byCode, code)
	return nil
}

type mfaTOTPTestRepo struct {
	byUser map[string]*models.UserMFATOTP
}

func newMFATOTPTestRepo() *mfaTOTPTestRepo {
	return &mfaTOTPTestRepo{byUser: map[string]*models.UserMFATOTP{}}
}

func (r *mfaTOTPTestRepo) GetByUserID(_ context.Context, userID string) (*models.UserMFATOTP, error) {
	row, ok := r.byUser[userID]
	if !ok {
		return nil, errors.NotFoundError("UserMFATOTP", nil)
	}
	return row, nil
}

func (r *mfaTOTPTestRepo) GetActiveByUserID(_ context.Context, userID string) (*models.UserMFATOTP, error) {
	row, ok := r.byUser[userID]
	if !ok || !row.Enabled {
		return nil, nil
	}
	return row, nil
}

func (r *mfaTOTPTestRepo) UpsertPending(_ context.Context, row *models.UserMFATOTP) error {
	r.byUser[row.UserID] = row
	return nil
}

func (r *mfaTOTPTestRepo) MarkVerifiedAndEnabled(_ context.Context, userID string) error {
	row, ok := r.byUser[userID]
	if !ok {
		return errors.NotFoundError("UserMFATOTP", nil)
	}
	row.Enabled = true
	return nil
}

func (r *mfaTOTPTestRepo) Disable(_ context.Context, userID string) error {
	delete(r.byUser, userID)
	return nil
}

func (r *mfaTOTPTestRepo) CountEnabled(_ context.Context) (int64, error) {
	var n int64
	for _, row := range r.byUser {
		if row.Enabled {
			n++
		}
	}
	return n, nil
}

type mfaRecoveryTestRepo struct {
	byUser map[string][]models.UserMFARecoveryCode
	used   []string
}

func newMFARecoveryTestRepo() *mfaRecoveryTestRepo {
	return &mfaRecoveryTestRepo{byUser: map[string][]models.UserMFARecoveryCode{}}
}

func (r *mfaRecoveryTestRepo) ReplaceAllForUser(_ context.Context, userID string, rows []*models.UserMFARecoveryCode) error {
	out := make([]models.UserMFARecoveryCode, 0, len(rows))
	for _, row := range rows {
		out = append(out, *row)
	}
	r.byUser[userID] = out
	return nil
}

func (r *mfaRecoveryTestRepo) FindUnusedByUserID(_ context.Context, userID string) ([]models.UserMFARecoveryCode, error) {
	return r.byUser[userID], nil
}

func (r *mfaRecoveryTestRepo) MarkUsed(_ context.Context, id string) error {
	r.used = append(r.used, id)
	return nil
}

type fedIdentityTestRepo struct {
	byProviderSub map[string]*models.FederatedIdentity
	created       []*models.FederatedIdentity
}

func newFedIdentityTestRepo() *fedIdentityTestRepo {
	return &fedIdentityTestRepo{byProviderSub: map[string]*models.FederatedIdentity{}}
}

func (r *fedIdentityTestRepo) key(provider, sub string) string {
	return provider + ":" + sub
}

func (r *fedIdentityTestRepo) GetByProviderSubject(_ context.Context, provider, subject string) (*models.FederatedIdentity, error) {
	fi, ok := r.byProviderSub[r.key(provider, subject)]
	if !ok {
		return nil, errors.NotFoundError("FederatedIdentity", nil)
	}
	return fi, nil
}

func (r *fedIdentityTestRepo) ListByUserID(_ context.Context, _ string) ([]models.FederatedIdentity, error) {
	return nil, nil
}

func (r *fedIdentityTestRepo) Create(_ context.Context, link *models.FederatedIdentity) error {
	r.byProviderSub[r.key(link.Provider, link.Subject)] = link
	r.created = append(r.created, link)
	return nil
}

type mockFedProvider struct {
	id          string
	name        string
	configured  bool
	enabled     bool
	redirectURL string
	claims      *domains.OIDCUserClaims
	exchangeErr error
}

func (m *mockFedProvider) ID() string          { return m.id }
func (m *mockFedProvider) DisplayName() string { return m.name }

func (m *mockFedProvider) OAuthConfiguredForTenant(context.Context, string) (bool, error) {
	return m.configured, nil
}

func (m *mockFedProvider) AuthorizeRedirectURL(context.Context, string, string, string) (string, error) {
	if m.redirectURL == "" {
		return "https://idp.example/authorize", nil
	}
	return m.redirectURL, nil
}

func (m *mockFedProvider) ExchangeAuthorizationCode(context.Context, string, string, string) (*domains.OIDCUserClaims, error) {
	if m.exchangeErr != nil {
		return nil, m.exchangeErr
	}
	return m.claims, nil
}
