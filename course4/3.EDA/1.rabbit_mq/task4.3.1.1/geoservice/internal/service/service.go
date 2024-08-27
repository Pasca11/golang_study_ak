package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"geoservice-rabbit/internal/cache"
	"go.uber.org/ratelimit"
	"io"
	"net/http"
	"strings"
	"time"

	"geoservice-rabbit/internal/repository"
	"geoservice-rabbit/models"

	"github.com/go-chi/jwtauth"
)

var lims = make(map[string]ratelimit.Limiter)

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
	SearchAnswer(token string, coordinates models.RequestAddressSearch) (int, *models.ResponseAddress, error)
	GeocodeAnswer(address models.Address) (int, []models.GetCoords, error)
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
	databaseUser, ok, err := c.Database.GetByName(context.Background(), user.Username)

	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("failed to get user: %w", err)
	}
	if !ok {
		return http.StatusForbidden, "", fmt.Errorf("user doesn`t exist")
	}

	passwordHash := hashPassword([]byte(user.Password))
	if passwordHash != databaseUser.Password {
		return http.StatusForbidden, "", fmt.Errorf("invalid password")
	}

	claims := map[string]interface{}{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}

	_, tokenString, _ := c.Token.Encode(claims)

	return http.StatusOK, tokenString, nil
}

func (c *GeoService) SearchAnswer(token string, coordinates models.RequestAddressSearch) (int, *models.ResponseAddress, error) {
	lim, err := c.getLimiter(token)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to get limiter: %w", err)
	}
	take := tryTake(lim)
	if !take {
		return http.StatusTooManyRequests, nil, fmt.Errorf("too many requests")
	}

	var address models.ResponseAddress
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", coordinates.Lat, coordinates.Lng)
	resp, err := http.Get(url)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	err = json.Unmarshal(body, &address)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, &address, nil
}

func (c *GeoService) GeocodeAnswer(address models.Address) (int, []models.GetCoords, error) {
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

func (c *GeoService) GetByID(id string) (models.User, error) {
	return c.Database.GetByID(context.Background(), id)
}

func hashPassword(password []byte) string {
	hash := sha256.New()
	hash.Write(password)
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *GeoService) getLimiter(tokenS string) (ratelimit.Limiter, error) {
	token, err := s.Token.Decode(tokenS)
	if err != nil {
		return nil, err
	}
	u, ok := token.Get("username")
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	username := u.(string)

	res, ok := lims[username]
	if !ok {
		res = ratelimit.New(5, ratelimit.Per(time.Minute))
		lims[username] = res
	}
	return res, nil
}

func tryTake(lim ratelimit.Limiter) bool {
	timer := time.NewTimer(time.Millisecond * 10)
	finish := make(chan struct{})
	go func() {
		defer close(finish)
		lim.Take()
		finish <- struct{}{}
	}()
	select {
	case <-timer.C:
		return false
	case <-finish:
		return true
	}
}
