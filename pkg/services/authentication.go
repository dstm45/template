// Package services contains all the sqelette
package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dstm45/template/pkg/config"
	"github.com/dstm45/template/pkg/database"
	"github.com/dstm45/template/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = []byte(config.LoadConfig().Secret)

type AuthService interface {
	SignIn(ctx context.Context, email, password string) (*http.Cookie, *http.Cookie, error)
	RotateToken(ctx context.Context, r *http.Request) (*http.Cookie, *http.Cookie, error)
	ParseRefreshToken(tokenString string) (*RefreshTokenClaim, error)
	ParseAccessToken(tokenString string) (*AccessTokenClaim, error)
	DeleteToken(ctx context.Context, accessTokenClaim string) error
}

type RefreshTokenClaim struct {
	Matricule   string    `json:"matricule"`
	Role        string    `json:"role"`
	UUID        uuid.UUID `json:"uuid"`
	TokenFamily uuid.UUID `json:"tokenFamily"`
	jwt.RegisteredClaims
}

type AccessTokenClaim struct {
	UUID        uuid.UUID `json:"uuid"`
	Role        string    `json:"role"`
	TokenFamily uuid.UUID `json:"tokenFamily"`
	jwt.RegisteredClaims
}

type authService struct {
	DB *database.Queries
}

func NewAuthService(db *database.Queries) AuthService {
	service := authService{
		DB: db,
	}
	return &service
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*http.Cookie, *http.Cookie, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		utils.CheckHash(password, "placeholder")
		return nil, nil, errors.New("mauvais Email ou mot de passe")
	}

	err = utils.CheckHash(password, user.PasswordHash)
	if err != nil {
		return nil, nil, errors.New("mauvais Email ou mot de passe")
	}

	tokenFamily, err := s.DB.CreateTokenFamily(ctx, user.Uuid)
	if err != nil {
		return nil, nil, errors.New("arreur lors de la création du taoken")
	}
	userData, err := s.DB.GetUserDataByUUID(ctx, user.Uuid)
	if err != nil {
		return nil, nil, err
	}
	return s.createTokens(ctx, userData, tokenFamily)
}

func (s *authService) RotateToken(ctx context.Context, r *http.Request) (*http.Cookie, *http.Cookie, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, nil, errors.New("unauthorized")
	}
	tokenString := cookie.Value

	claims, err := s.ParseRefreshToken(tokenString)
	if err != nil {
		// If parsing fails, the token is invalid or expired.
		// We could decide to invalidate the family here for extra security,
		// but for now we just deny access.
		return nil, nil, err
	}

	// Check for token reuse
	tokenHash := utils.HashToken(tokenString)
	_, err = s.DB.GetTokenByHash(ctx, tokenHash)
	// If token is not in DB, it might have been used already or is invalid.
	// This is a potential sign of a replay attack.
	if err != nil {
		// Invalidate the entire token family to force re-authentication.
		s.DB.DeleteTokensByFamily(ctx, claims.TokenFamily)
		return nil, nil, errors.New("unauthorized: token reuse suspected")
	}
	user, err := s.DB.GetUserDataByUUID(ctx, claims.UUID)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	tokenFamily, err := s.DB.GetTokenFamilyByUUID(ctx, claims.TokenFamily)
	if err != nil {
		return nil, nil, errors.New("token family not found")
	}

	// Invalidate the used token hash
	err = s.DB.DeleteTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, nil, err
	}

	return s.createTokens(ctx, user, tokenFamily)
}

func (s *authService) createTokens(ctx context.Context, user database.UserPublicDatum, tokenFamily database.TokenFamily) (*http.Cookie, *http.Cookie, error) {
	refreshTokenclaims, accessTokenclaims := s.newTokenClaims(user, &tokenFamily)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenclaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, nil, err
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  refreshTokenclaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   false, // Should be true in production
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, nil, err
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessTokenString,
		Expires:  accessTokenclaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   false, // Should be true in production
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	refreshTokenHash := utils.HashToken(refreshTokenString)
	params := database.CreateTokenParams{
		Family: tokenFamily.Uuid,
		Hash:   refreshTokenHash,
	}
	_, err = s.DB.CreateToken(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	return refreshCookie, accessCookie, nil
}

func (s *authService) newTokenClaims(user database.UserPublicDatum, tokenFamily *database.TokenFamily) (RefreshTokenClaim, AccessTokenClaim) {
	refreshTokenclaims := RefreshTokenClaim{
		Role:        string(user.Role),
		UUID:        user.Uuid,
		TokenFamily: tokenFamily.Uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sixnet",
			Subject:   "auth",
		},
	}
	accessTokenclaims := AccessTokenClaim{
		Role:        string(user.Role),
		UUID:        user.Uuid,
		TokenFamily: tokenFamily.Uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sixnet",
			Subject:   "auth",
		},
	}
	return refreshTokenclaims, accessTokenclaims
}

func (s *authService) ParseRefreshToken(tokenString string) (*RefreshTokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaim{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *authService) ParseAccessToken(tokenString string) (*AccessTokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaim{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenClaims, ok := token.Claims.(*AccessTokenClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return accessTokenClaims, nil
}

func (s *authService) DeleteToken(ctx context.Context, accessTokenString string) error {
	accessTokenClaim, err := s.ParseAccessToken(accessTokenString)
	if err != nil {
		return err
	}
	s.DB.DeleteTokensByFamily(ctx, accessTokenClaim.TokenFamily)
	return nil
}
