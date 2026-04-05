package auth

import (
	"context"

	"go-server/internal/errcode"
	appjwt "go-server/internal/jwt"
	"go-server/internal/modules/captcha"
	"go-server/internal/modules/user"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest) (TokenResponse, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)
}

type service struct {
	userSvc    user.Service
	jwtManager *appjwt.Manager
	blacklist  *appjwt.Blacklist
	captchaSvc captcha.Service
}

func NewService(userSvc user.Service, jwtManager *appjwt.Manager, blacklist *appjwt.Blacklist, captchaSvc captcha.Service) Service {
	return &service{userSvc: userSvc, jwtManager: jwtManager, blacklist: blacklist, captchaSvc: captchaSvc}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (TokenResponse, error) {
	if !s.captchaSvc.Verify(req.CaptchaID, req.CaptchaCode) {
		return TokenResponse{}, errcode.ErrCaptchaInvalid.AsError()
	}

	u, err := s.userSvc.GetByUsername(ctx, req.Username)
	if err != nil {
		// 统一返回"用户名或密码错误"，不暴露用户是否存在
		return TokenResponse{}, errcode.ErrInvalidCredentials.AsError()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return TokenResponse{}, errcode.ErrInvalidCredentials.AsError()
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(u.ID, u.Username)
	if err != nil {
		return TokenResponse{}, errcode.ErrInternalError.AsError()
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(u.ID, u.Username)
	if err != nil {
		return TokenResponse{}, errcode.ErrInternalError.AsError()
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, accessToken string) error {
	if s.blacklist == nil {
		return errcode.ErrInternalError.AsError()
	}
	claims, err := s.jwtManager.Parse(accessToken)
	if err != nil {
		return errcode.ErrUnauthorized.AsError()
	}
	if err := s.blacklist.Add(ctx, accessToken, claims.ExpiresAt.Time); err != nil {
		return errcode.ErrInternalError.AsError()
	}
	return nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error) {
	claims, err := s.jwtManager.Parse(refreshToken)
	if err != nil {
		return TokenResponse{}, errcode.ErrUnauthorized.AsError()
	}
	if claims.TokenType != appjwt.RefreshToken {
		return TokenResponse{}, errcode.ErrUnauthorized.AsError()
	}

	if s.blacklist != nil {
		blocked, err := s.blacklist.IsBlocked(ctx, refreshToken)
		if err != nil || blocked {
			return TokenResponse{}, errcode.ErrUnauthorized.AsError()
		}
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return TokenResponse{}, errcode.ErrInternalError.AsError()
	}
	return TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
