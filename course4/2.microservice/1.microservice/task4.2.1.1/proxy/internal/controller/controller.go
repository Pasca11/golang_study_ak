package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proxy/internal/service"
	"proxy/models"
	authv1 "proxy/proto/gen/auth"
	"proxy/proto/gen/geo"
)

type Controller struct {
	service service.ProxyService
}

func NewController(service service.ProxyService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var user authv1.User
	json.NewDecoder(r.Body).Decode(&user)

	token, err := c.service.Login(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tok := models.AuthToken{Token: token}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tok)
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var user authv1.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "err decoding user", http.StatusInternalServerError)
		return
	}
	err = c.service.Register(&user)
	if err != nil {
		http.Error(w, fmt.Errorf("Contrl: %w", err).Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GeocodeAnswer(w http.ResponseWriter, r *http.Request) {
	var address geo.Address
	json.NewDecoder(r.Body).Decode(&address)

	resp, err := c.service.GeocodeAnswer(&address)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) SearchAnswer(w http.ResponseWriter, r *http.Request) {
	var address geo.RequestAddressSearch
	json.NewDecoder(r.Body).Decode(&address)

	res, err := c.service.SearchAnswer(&address)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}
		_, err := c.service.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}
