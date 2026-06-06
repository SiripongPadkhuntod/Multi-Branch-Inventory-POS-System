package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/infrastructure/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthRepo struct {
	byEmail       map[string]*domain.User
	byID          map[uuid.UUID]*domain.User
	refreshTokens map[string]fakeRefreshToken
}

type fakeRefreshToken struct {
	userID    uuid.UUID
	expiresAt time.Time
	revoked   bool
}

func (r *fakeAuthRepo) FindUserByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (r *fakeAuthRepo) FindUserByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (r *fakeAuthRepo) StoreRefreshToken(_ context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	if r.refreshTokens == nil {
		r.refreshTokens = map[string]fakeRefreshToken{}
	}
	r.refreshTokens[tokenHash] = fakeRefreshToken{userID: userID, expiresAt: expiresAt}
	return nil
}

func (r *fakeAuthRepo) IsRefreshTokenActive(_ context.Context, userID uuid.UUID, tokenHash string) (bool, error) {
	token, ok := r.refreshTokens[tokenHash]
	if !ok {
		return false, nil
	}
	return token.userID == userID && !token.revoked && token.expiresAt.After(time.Now()), nil
}

func (r *fakeAuthRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	token, ok := r.refreshTokens[tokenHash]
	if !ok {
		return nil
	}
	token.revoked = true
	r.refreshTokens[tokenHash] = token
	return nil
}

func TestAuthServiceLoginIssuesAccessAndRefreshTokens(t *testing.T) {
	service, user := newTestAuthService(t, "ACTIVE")

	access, refresh, returnedUser, err := service.Login(context.Background(), user.Email, "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if returnedUser.ID != user.ID {
		t.Fatalf("Login() user ID = %s, want %s", returnedUser.ID, user.ID)
	}

	accessClaims, err := service.ParseToken(access)
	if err != nil {
		t.Fatalf("ParseToken(access) error = %v", err)
	}
	if accessClaims.TokenUse != "access" {
		t.Fatalf("access token_use = %q, want access", accessClaims.TokenUse)
	}
	if accessClaims.UserID != user.ID {
		t.Fatalf("access user_id = %s, want %s", accessClaims.UserID, user.ID)
	}

	refreshClaims, err := service.ParseToken(refresh)
	if err != nil {
		t.Fatalf("ParseToken(refresh) error = %v", err)
	}
	if refreshClaims.TokenUse != "refresh" {
		t.Fatalf("refresh token_use = %q, want refresh", refreshClaims.TokenUse)
	}
}

func TestAuthServiceRefreshRotatesTokens(t *testing.T) {
	service, user := newTestAuthService(t, "ACTIVE")
	_, refresh, _, err := service.Login(context.Background(), user.Email, "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	access2, refresh2, refreshedUser, err := service.Refresh(context.Background(), refresh)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if refreshedUser.ID != user.ID {
		t.Fatalf("Refresh() user ID = %s, want %s", refreshedUser.ID, user.ID)
	}

	accessClaims, err := service.ParseToken(access2)
	if err != nil {
		t.Fatalf("ParseToken(new access) error = %v", err)
	}
	if accessClaims.TokenUse != "access" {
		t.Fatalf("new access token_use = %q, want access", accessClaims.TokenUse)
	}

	refreshClaims, err := service.ParseToken(refresh2)
	if err != nil {
		t.Fatalf("ParseToken(new refresh) error = %v", err)
	}
	if refreshClaims.TokenUse != "refresh" {
		t.Fatalf("new refresh token_use = %q, want refresh", refreshClaims.TokenUse)
	}
}

func TestAuthServiceRefreshRejectsReusedRefreshToken(t *testing.T) {
	service, user := newTestAuthService(t, "ACTIVE")
	_, refresh, _, err := service.Login(context.Background(), user.Email, "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if _, _, _, err := service.Refresh(context.Background(), refresh); err != nil {
		t.Fatalf("first Refresh() error = %v", err)
	}
	if _, _, _, err := service.Refresh(context.Background(), refresh); err == nil {
		t.Fatal("second Refresh(reused token) error = nil, want error")
	}
}

func TestAuthServiceRefreshRejectsAccessToken(t *testing.T) {
	service, user := newTestAuthService(t, "ACTIVE")
	access, _, _, err := service.Login(context.Background(), user.Email, "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if _, _, _, err := service.Refresh(context.Background(), access); err == nil {
		t.Fatal("Refresh(access token) error = nil, want error")
	}
}

func TestAuthServiceRefreshRejectsInactiveUser(t *testing.T) {
	service, user := newTestAuthService(t, "INACTIVE")
	refresh, err := service.sign(user, time.Hour, "refresh")
	if err != nil {
		t.Fatalf("sign(refresh) error = %v", err)
	}

	if _, _, _, err := service.Refresh(context.Background(), refresh); err == nil {
		t.Fatal("Refresh(inactive user) error = nil, want error")
	}
}

func newTestAuthService(t *testing.T, status string) (*AuthService, *domain.User) {
	t.Helper()

	branchID := uuid.New()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	user := &domain.User{
		ID:           uuid.New(),
		BranchID:     &branchID,
		Role:         domain.RoleOwner,
		Name:         "Owner",
		Email:        "owner@example.com",
		PasswordHash: string(passwordHash),
		Status:       status,
		CreatedAt:    time.Now(),
	}
	repo := &fakeAuthRepo{
		byEmail:       map[string]*domain.User{user.Email: user},
		byID:          map[uuid.UUID]*domain.User{user.ID: user},
		refreshTokens: map[string]fakeRefreshToken{},
	}

	return &AuthService{
		repo: repo,
		cfg: config.Config{
			JWTSecret:  "test-secret",
			AccessTTL:  time.Minute,
			RefreshTTL: time.Hour,
		},
	}, user
}
