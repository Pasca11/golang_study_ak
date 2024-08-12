package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"rpc_service/models"
	"strings"
)

type GeoService struct{}

func (c *GeoService) SearchAnswer(coordinates models.RequestAddressSearch, res *models.ResponseAddress) error {
	var address models.ResponseAddress
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", coordinates.Lat, coordinates.Lng)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &address)
	if err != nil {
		return err
	}

	*res = address
	return nil
}

func (c *GeoService) GeocodeAnswer(address models.Address, res *[]models.GetCoords) error {
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
		return err
	}

	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(answer, &coords)
	if err != nil {
		return err
	}

	*res = coords
	return nil
}

func main() {
	geoService := new(GeoService)
	rpc.Register(geoService)

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rpc listening on port 1234")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(conn)
	}
}
