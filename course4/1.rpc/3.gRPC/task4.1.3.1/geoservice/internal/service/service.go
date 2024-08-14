package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"geoservice/internal/cache"
	"io"
	"net/http"
	"strings"
	"time"

	repository "geoservice/internal/repository"
	models "geoservice/models"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
)

type ServiceOption func(*Service)

type Service struct {
	Database repository.UserRepository
	Token    *jwtauth.JWTAuth
}

func WithToken(token *jwtauth.JWTAuth) ServiceOption {
	return func(c *Service) {
		c.Token = token
	}
}

func NewService(options ...ServiceOption) (*Service, error) {
	cachedDB, err := cache.NewCache(context.Background())
	if err != nil {
		return nil, err
	}

	service := &Service{Database: cachedDB}

	for _, option := range options {
		option(service)
	}

	return service, nil
}

type Servicer interface {
	RegisterUser(user models.User) (int, error)
	LoginUser(user models.User) (int, string, error)
	GetByID(id string) (models.User, error)
}

func (c *Service) RegisterUser(user models.User) (int, error) {
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

func (c *Service) LoginUser(user models.User) (int, string, error) {
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

func (c *Service) SearchAnswer(coordinates models.RequestAddressSearch) (int, models.ResponseAddress, error) {
	var address models.ResponseAddress
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", coordinates.Lat, coordinates.Lng)
	resp, err := http.Get(url)

	if err != nil {
		return http.StatusInternalServerError, address, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, address, err
	}

	err = json.Unmarshal(body, &address)
	if err != nil {
		return http.StatusInternalServerError, address, err
	}

	return http.StatusOK, address, nil
}

func (c *Service) GeocodeAnswer(address models.Address) (int, []models.GetCoords, error) {
	parts := []string{}
	parts = append(parts, strings.Split(address.House_number, " ")...)
	parts = append(parts, strings.Split(address.Road, " ")...)
	parts = append(parts, strings.Split(address.Suburb, " ")...)
	parts = append(parts, strings.Split(address.City, " ")...)
	parts = append(parts, strings.Split(address.State, " ")...)
	parts = append(parts, strings.Split(address.Country, " ")...)

	var sb strings.Builder
	for _, i := range parts {
		if i != "" {
			sb.WriteString("+")
			sb.WriteString(i)
		}
	}

	request := "https://nominatim.openstreetmap.org/search?q=" + strings.Trim(sb.String(), "+") + "&format=json"
	var coords []models.GetCoords

	resp, err := http.Get(request)
	if err != nil {
		return http.StatusInternalServerError, coords, err
	}

	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, coords, err
	}

	err = json.Unmarshal(answer, &coords)
	if err != nil {
		return http.StatusInternalServerError, coords, err
	}

	return http.StatusOK, coords, nil
}

func (c *Service) GetByID(id string) (models.User, error) {
	return c.Database.GetByID(context.Background(), id)
}

func hashPassword(password []byte) string {
	hash := sha256.New()
	hash.Write(password)
	return hex.EncodeToString(hash.Sum(nil))
}
