package service

import (
	"auth/internal/repository"
	authv1 "auth/proto/gen/auth"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"time"
)

type AuthService interface {
	authv1.AuthServer
}

type AuthServiceImpl struct {
	authv1.UnimplementedAuthServer
	db     repository.AuthRepository
	secret []byte
}

func RegServer(grpc *grpc.Server, impl AuthService) {
	authv1.RegisterAuthServer(grpc, impl)
}

func NewAuthService(db repository.AuthRepository, secret []byte, srv *grpc.Server) {
	s := &AuthServiceImpl{
		db:     db,
		secret: secret,
	}
	RegServer(srv, s)
}

func (s *AuthServiceImpl) Register(ctx context.Context, user *authv1.User) (*authv1.User, error) {
	_, err := s.db.FindByUsername(user.Username)
	if err == nil {
		return nil, fmt.Errorf("user with username %s already exists", user.Username)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)
	err = s.db.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, user *authv1.User) (*authv1.Token, error) {
	userDB, err := s.db.FindByUsername(user.Username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		return nil, err
	}
	token, err := s.CreateJWT(user)
	if err != nil {
		return nil, err
	}
	return &authv1.Token{Token: token}, nil
}

func (s *AuthServiceImpl) CreateJWT(user *authv1.User) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Second * 300,
		"sub": user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token *authv1.Token) (*authv1.Token, error) {
	t, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		if !token.Valid {
			return nil, fmt.Errorf("invalid token")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return &authv1.Token{Token: t.Raw}, nil
}
