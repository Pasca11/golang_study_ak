package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"net/http"
	authv1 "proxy/proto/gen/auth"
	"proxy/proto/gen/geo"
	"strings"
)

type ProxyService interface {
	Login(user *authv1.User) (string, error)
	ValidateToken(stringToken string) (*jwt.Token, error)
	Register(user *authv1.User) error
	GeocodeAnswer(address *geo.Address) (*geo.GetCoords, error)
	SearchAnswer(coordinates *geo.RequestAddressSearch) (*geo.ResponseAddress, error)
}

type ProxyServiceImpl struct {
	gRPCGeo  *grpc.ClientConn
	AuthgRPC *grpc.ClientConn
}

func NewProxyService(client *grpc.ClientConn, auth *grpc.ClientConn) ProxyService {
	return &ProxyServiceImpl{
		gRPCGeo:  client,
		AuthgRPC: auth,
	}
}

func (s *ProxyServiceImpl) Login(user *authv1.User) (string, error) {
	client := authv1.NewAuthClient(s.AuthgRPC)
	token, err := client.Login(context.Background(), user)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (s *ProxyServiceImpl) ValidateToken(token string) (*jwt.Token, error) {
	resp, err := http.Post("/auth/validateToken", "application/json", strings.NewReader(token))
	if err != nil {
		return nil, err
	}
	var resBody jwt.Token
	json.NewDecoder(resp.Body).Decode(&resBody)

	return &resBody, nil
}

func (s *ProxyServiceImpl) Register(user *authv1.User) error {
	//body, err := json.Marshal(user)
	//if err != nil {
	//	return fmt.Errorf("service, json marshal err: %w", err)
	//}
	if user == nil {
		return fmt.Errorf("service. user is nil")
	}
	client := authv1.NewAuthClient(s.AuthgRPC)
	if client == nil {
		return fmt.Errorf("service. client nil")
	}
	_, err := client.Register(context.Background(), user)
	//_, err = http.Post("http://auth:4444/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("service, http post err: %w", err)
	}

	return nil
}

func (s *ProxyServiceImpl) SearchAnswer(coordinates *geo.RequestAddressSearch) (*geo.ResponseAddress, error) {
	client := geo.NewGeoServiceClient(s.gRPCGeo)
	req := &geo.RequestAddressSearch{
		Lat: coordinates.Lat,
		Lng: coordinates.Lng,
	}

	res, err := client.SearchAnswer(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProxyServiceImpl) GeocodeAnswer(address *geo.Address) (*geo.GetCoords, error) {
	client := geo.NewGeoServiceClient(s.gRPCGeo)

	req := &geo.Address{
		HouseNumber: address.HouseNumber,
		Road:        address.Road,
		Suburb:      address.Suburb,
		City:        address.City,
		State:       address.State,
		Country:     address.Country,
	}

	res, err := client.GeoCodeAnswer(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
