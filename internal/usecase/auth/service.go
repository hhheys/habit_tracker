package auth

import (
	"context"
	"habit-tracker/internal/domain"
	useruc "habit-tracker/internal/usecase/user"
	"time"
)

type Service struct {
	users           UserRepository
	refreshSessions RefreshSessionRepository
	passwordHasher  PasswordHasher
	tokenHasher     TokenHasher
	accessTokens    AccessTokenGenerator
	refreshTokens   RefreshTokenGenerator
}

func NewService(
	users UserRepository,
	refreshSessions RefreshSessionRepository,
	passwordHasher PasswordHasher,
	tokenHasher TokenHasher,
	accessTokens AccessTokenGenerator,
	refreshTokens RefreshTokenGenerator,
) *Service {
	return &Service{
		users:           users,
		refreshSessions: refreshSessions,
		passwordHasher:  passwordHasher,
		tokenHasher:     tokenHasher,
		accessTokens:    accessTokens,
		refreshTokens:   refreshTokens,
	}
}

func (s *Service) GenerateAuthTokens(ctx context.Context, subject *TokenSubject, session SessionInfoInput) (*TokenOutput, error) {
	tokens, sessionInfo, err := s.buildAuthTokens(subject, session)
	if err != nil {
		return nil, err
	}
	if err := s.refreshSessions.Create(ctx, sessionInfo); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *Service) buildAuthTokens(subject *TokenSubject, session SessionInfoInput) (*TokenOutput, *domain.RefreshSession, error) {
	accessToken, err := s.accessTokens.GenerateToken(subject)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.refreshTokens.GenerateToken(subject)
	if err != nil {
		return nil, nil, err
	}

	tokenHash, err := s.tokenHasher.HashToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	sessionInfo := domain.RefreshSession{
		UserID:    subject.UserID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		UserAgent: session.UserAgent,
		IPAddress: session.UserIP,
		Revoked:   false,
	}

	return &TokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, &sessionInfo, nil
}

func (s *Service) Register(ctx context.Context, user *RegisterInput, session SessionInfoInput) (*RegisterOutput, error) {
	// Validate user input
	exists, err := s.users.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	exists, err = s.users.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := s.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	u := &domain.User{
		Username: user.Username,
		Email:    user.Email,
		Password: passwordHash,
		Role:     domain.UserRoleDefault,
	}

	err = s.users.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	tokens, err := s.GenerateAuthTokens(ctx, &TokenSubject{UserID: u.ID, Role: u.Role}, session)
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		User: useruc.Output{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		},
		Authorization: TokenOutput{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil

}

func (s *Service) Login(ctx context.Context, user *LoginInput, session SessionInfoInput) (*LoginOutput, error) {
	var (
		u   *domain.User
		err error
	)
	if user.Email != "" {
		u, err = s.users.GetByEmail(ctx, user.Email)
	} else {
		u, err = s.users.GetByUsername(ctx, user.Username)
	}
	if err != nil {
		return nil, err
	}

	ok := s.passwordHasher.ComparePasswordHash(u.Password, user.Password)
	if !ok {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.GenerateAuthTokens(ctx, &TokenSubject{UserID: u.ID, Role: u.Role}, session)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		User: useruc.Output{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		},
		Tokens: TokenOutput{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, input *RefreshTokenInput, session SessionInfoInput) (*TokenOutput, error) {
	tokenHash, err := s.tokenHasher.HashToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	oldSession, err := s.refreshSessions.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	if oldSession.Revoked {
		return nil, domain.ErrSessionRevoked
	}
	if time.Now().After(oldSession.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	u, err := s.users.GetByID(ctx, oldSession.UserID)
	if err != nil {
		return nil, err
	}

	tokens, replacement, err := s.buildAuthTokens(&TokenSubject{UserID: u.ID, Role: u.Role}, session)
	if err != nil {
		return nil, err
	}
	if err := s.refreshSessions.Rotate(ctx, tokenHash, replacement); err != nil {
		return nil, err
	}

	return &TokenOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
