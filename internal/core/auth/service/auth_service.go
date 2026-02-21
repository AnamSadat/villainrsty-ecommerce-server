package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"math/big"
	"net/url"
	"time"

	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthService struct {
	userRepo          ports.UserRepository
	refreshTokenRepo  ports.RefreshTokenRepository
	passwordResetRepo ports.PasswordResetTokenRepository
	twoFactorOTPRepo  ports.TwoFactorOTPRepository
	emailSender       ports.EmailSender
	hasher            ports.PasswordHasher
	tokenHasher       ports.TokenHasher
	jwtService        ports.JWTService
	logger            *slog.Logger
	resetURL          string
	resetTTL          time.Duration
	twoFactorOTPTTL   time.Duration
}

func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	passwordResetRepo ports.PasswordResetTokenRepository,
	twoFactorOTPRepo ports.TwoFactorOTPRepository,
	emailSender ports.EmailSender,
	hasher ports.PasswordHasher,
	tokenHasher ports.TokenHasher,
	jwtService ports.JWTService,
	logger *slog.Logger,
	resetURL string,
	resetTTL time.Duration,
	twoFactorOTPTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		hasher:            hasher,
		tokenHasher:       tokenHasher,
		passwordResetRepo: passwordResetRepo,
		twoFactorOTPRepo:  twoFactorOTPRepo,
		emailSender:       emailSender,
		refreshTokenRepo:  refreshTokenRepo,
		jwtService:        jwtService,
		logger:            logger,
		resetURL:          resetURL,
		resetTTL:          resetTTL,
		twoFactorOTPTTL:   twoFactorOTPTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string, rememberMe bool) (*models.User, string, string, error) {
	if email == "" || password == "" {
		return nil, "", "", errors.New(errors.ErrValidation, "email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	if user == nil {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	if !s.hasher.Verify(user.Password, password) {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrInternal, "failed to generate access token")
	}

	refreshTokenString, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrInternal, "failed to generate refresh token")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshTokenString)
	s.logger.Info("refresh token debug (login)",
		"refresh_len", len(refreshTokenString),
		"hash_len", len(tokenHash),
		"hash_prefix", tokenHash[:8],
	)

	if err != nil {
		s.logger.Warn("token hash", "error", err)
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to hash token", err)
	}

	ttl := 24 * time.Hour
	if rememberMe {
		ttl = 30 * 24 * time.Hour
		s.logger.Info("ttl", "ttl sekarang: ", ttl)
	}

	refreshToken := models.NewRefreshToken(user.ID, tokenHash, ttl)
	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to save refresh token", err)
	}

	return user, accessToken, refreshTokenString, nil
}

func (s *AuthService) LoginWith2FA(ctx context.Context, email, password string, _ bool) (string, error) {
	if email == "" || password == "" {
		return "", errors.New(errors.ErrValidation, "email and password is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	if !s.hasher.Verify(user.Password, password) {
		return "", errors.New(errors.ErrUnauthorized, "invalid email or password")
	}

	challengeID, err := generateSecureToken(24)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate challenge", err)
	}

	otpCode, err := generateNumericOTP(6)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate otp", err)
	}

	otpHash, err := s.tokenHasher.Hash(otpCode)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to hash otp", err)
	}

	otp := &models.TwoFactorOTP{
		ID:          models.NewID(),
		UserID:      user.ID,
		ChallengeID: challengeID,
		CodeHash:    otpHash,
		ExpiresAt:   time.Now().Add(s.twoFactorOTPTTL),
		CreatedAt:   time.Now(),
	}

	if err := s.twoFactorOTPRepo.Save(ctx, otp); err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to save otp", err)
	}

	if err := s.emailSender.SendLoginOTP(ctx, user.Email, otpCode); err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to send otp", err)
	}

	return challengeID, nil
}

func (s *AuthService) VerifyLogin2FA(ctx context.Context, challengeID, otpCode string, rememberMe bool) (*models.User, string, string, error) {
	if challengeID == "" || otpCode == "" {
		return nil, "", "", errors.New(errors.ErrValidation, "challenge_id and otp_code is required")
	}

	otp, err := s.twoFactorOTPRepo.GetByChallengeID(ctx, challengeID)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid or expired otp")
	}

	if !otp.IsValid() {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid or expired otp")
	}

	hash, err := s.tokenHasher.Hash(otpCode)
	if err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to hash otp", err)
	}

	if hash != otp.CodeHash {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid or expired otp")
	}

	if err := s.twoFactorOTPRepo.MarkUsed(ctx, otp.ID); err != nil {
		s.logger.Info("isi markused", "used_at", otp.ID.String())
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to mark otp used", err)
	}

	user, err := s.userRepo.GetByID(ctx, otp.UserID.String())
	s.logger.Info("isi otp id", "otp", otp.ID.String())
	if err != nil {
		s.logger.Info("error di getbyid verify", "error", err)
		return nil, "", "", errors.New(errors.ErrUnauthorized, "user not found")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to generate access token", err)
	}

	refreshTokenString, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", errors.Wrap(errors.ErrUnauthorized, "failed to generate refresh token", err)
	}

	refreshHash, err := s.tokenHasher.Hash(refreshTokenString)
	if err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}

	ttl := 24 * time.Hour
	if rememberMe {
		ttl = 30 * 24 * time.Hour
	}

	refreshToken := models.NewRefreshToken(user.ID, refreshHash, ttl)
	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, "", "", errors.Wrap(errors.ErrInternal, "failed to save refresh token", err)
	}

	return user, accessToken, refreshTokenString, nil
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New(errors.ErrValidation, "email and password are required")
	}

	exists, err := s.userRepo.Exist(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(errors.ErrConflict, "email already registered")
	}

	user := models.NewUser(email, password, name)
	if !user.IsPasswordValid(password) {
		return nil, errors.New(errors.ErrValidation, "password must contain uppercase, lowercase and number")
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password", err)
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New(errors.ErrValidation, "refresh token is required")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshToken)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}

	dbToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.New(errors.ErrUnauthorized, "refresh token not found")
	}

	if err := s.refreshTokenRepo.Revoke(ctx, dbToken.ID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to revoke token", err)
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", errors.New(errors.ErrValidation, "refresh token is required")
	}

	user, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New(errors.ErrUnauthorized, "invalid refresh token")
	}

	tokenHash, err := s.tokenHasher.Hash(refreshToken)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}
	s.logger.Info("refresh token debug (incoming)",
		"incoming_len", len(refreshToken),
		"hash_len", len(tokenHash),
		"hash_prefix", tokenHash[:8],
	)

	dbToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return "", "", errors.New(errors.ErrUnauthorized, "refresh token not found")
	}

	if !dbToken.IsValid() {
		return "", "", errors.New(errors.ErrUnauthorized, "refresh token is expired or revoked")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to generate access token", err)
	}

	newRefreshTokenString, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to generate refresh token", err)
	}

	newTokenHash, err := s.tokenHasher.Hash(newRefreshTokenString)
	if err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to hash refresh token", err)
	}

	newRefreshToken := models.NewRefreshToken(user.ID, newTokenHash, 7*24*time.Hour)
	if err := s.refreshTokenRepo.Save(ctx, newRefreshToken); err != nil {
		return "", "", errors.Wrap(errors.ErrInternal, "failed to save refresh token", err)
	}

	if err := s.refreshTokenRepo.Revoke(ctx, dbToken.ID); err != nil {
		s.logger.Error("failed to revoke old refresh token", "[ERROR]:", err)
	}

	return accessToken, newRefreshTokenString, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	if token == "" {
		return nil, errors.New(errors.ErrValidation, "token is required")
	}

	user, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.New(errors.ErrUnauthorized, "invalid token")
	}

	return user, nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	if email == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.IsKind(err, errors.ErrNotFound) {
			return nil
		}
		return err
	}

	rawToken, err := generateSecureToken(32)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to generate reset token", err)
	}

	tokenHash, err := s.tokenHasher.Hash(rawToken)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to hash reset token", err)
	}

	resetToken := &models.PasswordResetToken{
		ID:        models.NewID(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.resetTTL),
		CreatedAt: time.Now(),
	}

	if err := s.passwordResetRepo.Save(ctx, resetToken); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save password reset token", err)
	}

	resetLink, err := buildResetLink(s.resetURL, rawToken)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to build reset link", err)
	}
	s.logger.Info("sending reset email", "to", email, "link", resetLink)

	if err := s.emailSender.SendPasswordReset(ctx, user.Email, resetLink); err != nil {
		s.logger.Error("failed to send reset email", "to", email, "err", err)
		return errors.Wrap(errors.ErrInternal, "failed to send reset email", err)
	}
	s.logger.Info("reset email sent", "to", email)

	return nil
}

func (s *AuthService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if token == "" {
		return errors.New(errors.ErrValidation, "token is required")
	}

	if newPassword == "" {
		return errors.New(errors.ErrValidation, "new password is required")
	}

	tokenHash, err := s.tokenHasher.Hash(token)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to hash reset token", err)
	}

	dbToken, err := s.passwordResetRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsKind(err, errors.ErrNotFound) {
			return errors.New(errors.ErrUnauthorized, "invalid reset token")
		}
		return err
	}

	if !dbToken.IsValid() {
		return errors.New(errors.ErrUnauthorized, "invalid reset token")
	}

	if !models.NewUser("temp@mail.com", newPassword, "temp").IsPasswordValid(newPassword) {
		return errors.New(errors.ErrValidation, "password must contain uppercase, lowercase and number")
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to hash password", err)
	}

	if err := s.userRepo.UpdateUserPassword(ctx, dbToken.UserID, hashedPassword); err != nil {
		return err
	}

	if err := s.passwordResetRepo.MarkUsed(ctx, dbToken.ID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to mark reset token used", err)
	}
	return nil
}

func generateSecureToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateNumericOTP(digits int) (string, error) {
	if digits <= 0 {
		return "", errors.New(errors.ErrValidation, "invalid otp digits")
	}

	result := make([]byte, digits)
	for i := 0; i < digits; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}

	return string(result), nil
}

func buildResetLink(baseURL, token string) (string, error) {
	if baseURL == "" {
		return "", errors.New(errors.ErrInternal, "reset base url is empty")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
