package controller

import (
	"encoding/json"
	"fmt"
	"geoservice/internal/metrics"
	service "geoservice/internal/service"
	models "geoservice/models"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type ControllerOption func(*Controller)

type Controller struct {
	Servicer service.GeoServicer
}

func NewController(token *jwtauth.JWTAuth, options ...ControllerOption) *Controller {
	service, err := service.NewGeoService(service.WithToken(token))
	if err != nil {
		fmt.Println(err)
	}

	controller := &Controller{Servicer: service}
	return controller
}

type Responder interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	SearchAnswer(w http.ResponseWriter, r *http.Request)
	GeocodeAnswer(w http.ResponseWriter, r *http.Request)
	NotFoundAnswer(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

// registerUser handle POST-requests on api/register
// @Summary Register
// @Tags Login
// @Description Register user
// @Accept  json
// @Produce  json
// @Param  input  body  models.User true  "username and password"
// @Success 200 {object} string
// @Failure 400 {object} models.ErrorResponce
// @Failure 500 {object} models.ErrorResponce
// @Router /api/register [post]
func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		end := time.Since(start).Seconds()
		metrics.RegisterCounter.Inc()
		metrics.RegisterHistogram.Observe(end)
	}()
	var newUser models.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		newErrorResponce(w, err, http.StatusInternalServerError)
		return
	}

	status, err := c.Servicer.RegisterUser(newUser)
	if err != nil {
		newErrorResponce(w, err, status)
		return
	}

	message := fmt.Sprintf("User %s sucessfully created", newUser.Username)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// loginUser handle POST-requests on api/login
// @Summary Login
// @Tags Login
// @Description Login user
// @Accept  json
// @Produce  json
// @Param  input  body  models.User true  "username and password"
// @Success 200 {object} string
// @Failure 400 {object} models.ErrorResponce
// @Failure 500 {object} models.ErrorResponce
// @Router /api/login [post]
func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		end := time.Since(start).Seconds()
		metrics.LoginCounter.Inc()
		metrics.LoginHistogram.Observe(end)
	}()
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		newErrorResponce(w, err, http.StatusInternalServerError)
		return
	}

	status, tokenString, err := c.Servicer.LoginUser(user)
	if err != nil {
		newErrorResponce(w, err, status)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte(tokenString))
}

// searchAnswer handle POST-requests on api/address/search
// @Summary SearchCity
// @Tags Search
// @Description Search city Name by coords
// @Accept  json
// @Produce  json
// @Param  coordinates  body  models.RequestAddressSearch true  "Lattitude and Longitude"
// @Success 200 {object} string
// @Failure 400 {object} models.ErrorResponce
// @Failure 500 {object} models.ErrorResponce
// @Router /api/address/search [post]
func (c *Controller) SearchAnswer(w http.ResponseWriter, r *http.Request) {
	var coordinates models.RequestAddressSearch
	json.NewDecoder(r.Body).Decode(&coordinates)

	status, address, err := c.Servicer.SearchAnswer(coordinates)
	if err != nil {
		newErrorResponce(w, err, status)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte("you are in " + address.Address.City))
}

// geocodeAnswer handle POST-requests on api/address/geocode
// @Summary SearchCoords
// @Tags Search
// @Description Search coords by address
// @Accept  json
// @Produce  json
// @Param  coordinates  body  models.Address true  "House number, road, suburb, city, state, country"
// @Success 200 {object} string
// @Failure 400 {object} models.ErrorResponce
// @Failure 500 {object} models.ErrorResponce
// @Router /api/address/search [post]
func (c *Controller) GeocodeAnswer(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	json.NewDecoder(r.Body).Decode(&address)

	status, coords, err := c.Servicer.GeocodeAnswer(address)
	if err != nil {
		newErrorResponce(w, err, status)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte("Your lattitude = " + coords[0].Lat + "; Your longitude = " + coords[0].Lon))
}

// GetUserByID handle GET-requests on api/users/{id}
// @Summary Get username by ID
// @Tags Login
// @Description Search username by id
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 400 {object} models.ErrorResponce
// @Failure 500 {object} models.ErrorResponce
// @Router /api/users/{id} [get]
func (c *Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := c.Servicer.GetByID(id)

	if err != nil {
		newErrorResponce(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Username of user with ID " + id + " is " + user.Username))
}

func (c *Controller) NotFoundAnswer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func newErrorResponce(w http.ResponseWriter, err error, responce int) {
	errResponce := models.ErrorResponce{Message: err.Error()}
	http.Error(w, errResponce.Message, responce)
}
