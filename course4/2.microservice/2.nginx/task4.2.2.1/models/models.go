package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseAddress struct {
	Address Address `json:"address"`
}

type Address struct {
	House_number string `json:"house_number"`
	Road         string `json:"road"`
	Suburb       string `json:"suburb"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type RequestAddressSearch struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GetCoords struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type ErrorResponce struct {
	Message string `json:"message"`
}
