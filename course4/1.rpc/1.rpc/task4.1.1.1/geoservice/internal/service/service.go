package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"geoservice/internal/cache"
	"net/http"
	"time"

	repository "geoservice/internal/repository"
	models "geoservice/models"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
)

type ServiceOption func(*GeoService)

type GeoService struct {
	Database repository.UserRepository
	Token    *jwtauth.JWTAuth
}

func WithToken(token *jwtauth.JWTAuth) ServiceOption {
	return func(c *GeoService) {
		c.Token = token
	}
}

func NewGeoService(options ...ServiceOption) (*GeoService, error) {
	cachedDB, err := cache.NewCache(context.Background())
	if err != nil {
		return nil, err
	}

	service := &GeoService{Database: cachedDB}

	for _, option := range options {
		option(service)
	}

	return service, nil
}

type GeoServicer interface {
	RegisterUser(user models.User) (int, error)
	LoginUser(user models.User) (int, string, error)
	GetByID(id string) (models.User, error)
}

func (c *GeoService) RegisterUser(user models.User) (int, error) {
	_, ok, _ := c.Database.GetByName(context.Background(), user.Username)

	if ok {
		return http.StatusInternalServerError, fmt.Errorf("username already exist")
	}

	passwordHash := hashPassword([]byte(user.Password))
	user.Password = passwordHash
	err := c.Database.Create(context.Background(), user)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (c *GeoService) LoginUser(user models.User) (int, string, error) {
	databaseUser, ok, _ := c.Database.GetByName(context.Background(), user.Username)

	if !ok {
		return http.StatusForbidden, "", fmt.Errorf("user dont exist")
	}

	passwordHash := hashPassword([]byte(user.Password))
	if passwordHash != databaseUser.Password {
		return http.StatusForbidden, "", fmt.Errorf("invalid password")
	}

	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      jwtauth.ExpireIn(time.Hour),
	}
	_, tokenString, _ := c.Token.Encode(claims)

	return http.StatusOK, tokenString, nil
}

func (c *GeoService) GetByID(id string) (models.User, error) {
	return c.Database.GetByID(context.Background(), id)
}

func hashPassword(password []byte) string {
	hash := sha256.New()
	hash.Write(password)
	return hex.EncodeToString(hash.Sum(nil))
}
