package controller

import (
	"auth/internal/service"
	"auth/models"
	authv1 "auth/proto/gen/auth"
	"context"
	"encoding/json"
	"net/http"
)

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	ValidateToken(w http.ResponseWriter, r *http.Request)
}

type AuthControllerImpl struct {
	authService service.AuthService
	ctx         context.Context
}

func NewAuthController(authService service.AuthService) AuthController {
	ctrl := &AuthControllerImpl{}
	ctrl.authService = authService
	ctrl.ctx = context.Background()

	return ctrl
}

func (ac *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var req authv1.User

	json.NewDecoder(r.Body).Decode(&req)

	_, err := ac.authService.Register(ac.ctx, &req)
	if err != nil {
		errResp := models.ErrorResponse{Message: err.Error()}
		RenderJSON(w, http.StatusInternalServerError, errResp)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ac *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var req authv1.User

	json.NewDecoder(r.Body).Decode(&req)

	token, err := ac.authService.Login(nil, &req)
	if err != nil {
		resp := models.ErrorResponse{
			Message: err.Error(),
		}
		RenderJSON(w, http.StatusForbidden, resp)
		return
	}
	resp := token
	RenderJSON(w, http.StatusOK, resp)
}

func (ac *AuthControllerImpl) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var req authv1.Token

	json.NewDecoder(r.Body).Decode(&req)

	_, err := ac.authService.ValidateToken(ac.ctx, &req)
	if err != nil {
		errResp := models.ErrorResponse{
			Message: err.Error(),
		}
		RenderJSON(w, http.StatusUnauthorized, errResp)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RenderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(&v)
	if err != nil {
		http.Error(w, "Can`t render result", http.StatusInternalServerError)
	}
}
