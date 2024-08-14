package service

import (
	"geoservice/models"
	"net/http"
	"net/rpc"
)

type Option func(service *GeoRPCService)

type GeoRPC interface {
	SearchAnswer(coordinates models.RequestAddressSearch) (int, models.ResponseAddress, error)
	GeocodeAnswer(address models.Address) (int, []models.GetCoords, error)
}

type GeoRPCService struct {
	client *rpc.Client
}

func NewGeoRPCService(opts ...Option) *GeoRPCService {
	service := &GeoRPCService{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (c *GeoRPCService) SearchAnswer(coordinates models.RequestAddressSearch) (int, models.ResponseAddress, error) {
	var res models.ResponseAddress
	err := c.client.Call("GeoService.SearchAnswer", coordinates, &res)
	if err != nil {
		return http.StatusBadRequest, models.ResponseAddress{}, err
	}
	return http.StatusOK, res, nil
}

func (c *GeoRPCService) GeocodeAnswer(address models.Address) (int, []models.GetCoords, error) {
	var res []models.GetCoords
	err := c.client.Call("GeoService.Geocode", address, &res)
	if err != nil {
		return http.StatusInternalServerError, res, err
	}
	return http.StatusOK, res, nil
}
