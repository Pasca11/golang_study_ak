package service

import (
	"context"
	"geoservice/models"
	geov1 "geoservice/protos/gen"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type GeoServicer interface {
	SearchAnswer(coordinates models.RequestAddressSearch) (int, models.ResponseAddress, error)
	GeocodeAnswer(address models.Address) (int, []models.GetCoords, error)
}

type GeoService struct {
	client *grpc.ClientConn
}

func NewGeoService(address string) (*GeoService, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, err
	}
	service := &GeoService{
		client: conn,
	}
	return service, nil
}

func (s *GeoService) SearchAnswer(coordinates models.RequestAddressSearch) (int, models.ResponseAddress, error) {
	var address models.ResponseAddress

	client := geov1.NewGeoServiceClient(s.client)
	req := &geov1.RequestAddressSearch{
		Lat: coordinates.Lat,
		Lng: coordinates.Lng,
	}

	res, err := client.SearchAnswer(context.Background(), req)
	if err != nil {
		return http.StatusInternalServerError, address, err
	}

	address = models.ResponseAddress{
		Address: models.Address{
			House_number: res.Address.HouseNumber,
			Road:         res.Address.Road,
			Suburb:       res.Address.Suburb,
			City:         res.Address.City,
			State:        res.Address.State,
			Country:      res.Address.Country,
		},
	}

	return http.StatusOK, address, nil
}

func (s *GeoService) GeocodeAnswer(address models.Address) (int, []models.GetCoords, error) {
	client := geov1.NewGeoServiceClient(s.client)

	req := &geov1.Address{
		HouseNumber: address.House_number,
		Road:        address.Road,
		Suburb:      address.Suburb,
		City:        address.City,
		State:       address.State,
		Country:     address.Country,
	}

	res, err := client.GeoCodeAnswer(context.Background(), req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	var coords []models.GetCoords

	for _, v := range res.Coords {
		coords = append(coords, models.GetCoords{
			Lat: strconv.FormatFloat(v.Lat, 'f', -1, 64),
			Lon: strconv.FormatFloat(v.Lng, 'f', -1, 64),
		})
	}

	return http.StatusOK, coords, nil
}
