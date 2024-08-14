package service

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	geov1 "grpcService/protos/gen"
	"io"
	"net/http"
	"strings"
)

type Service struct {
	geov1.UnimplementedGeoServiceServer
}

func Register(grpc *grpc.Server) {
	geov1.RegisterGeoServiceServer(grpc, &Service{})
}

func (s *Service) SearchAnswer(ctx context.Context, coordinates *geov1.RequestAddressSearch) (*geov1.ResponseAddress, error) {
	var address *geov1.ResponseAddress
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", coordinates.Lat, coordinates.Lng)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *Service) GeoCodeAnswer(ctx context.Context, address *geov1.Address) (*geov1.GetCoords, error) {
	parts := []string{}
	parts = append(parts, strings.Split(address.HouseNumber, " ")...)
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
	coords := &geov1.GetCoords{
		Coords: make([]*geov1.Cooords, 0, 10),
	}

	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}

	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(answer, &coords.Coords)
	if err != nil {
		return nil, err
	}

	return coords, nil
}
