package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	geoalert "github.com/micro-gis/geo-alert-api/domain/geoalert"
	"github.com/micro-gis/geo-alert-api/domain/queries"
	"github.com/micro-gis/geo-alert-api/services"
	"github.com/micro-gis/geo-alert-api/utils/http_utils"
	"github.com/micro-gis/oauth-go/oauth"
	"github.com/micro-gis/utils/rest_errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	GeoAlertController geoalertsControllerInterface = &geoalertsController{}
)

type geoalertsControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	GetUserGeoAlerts(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Upsert(w http.ResponseWriter, r *http.Request)
}

type geoalertsController struct {
}

func (c *geoalertsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.AuthenticateRequest(r); err != nil {
		fmt.Println(err)
		http_utils.ResponseError(w, err)
		return
	}

	userId := oauth.GetCallerId(r)

	if userId == 0 {
		restErr := rest_errors.NewUnauthorizedError("user not authenticated")
		http_utils.ResponseError(w, restErr)
		return
	}
	var geoRequest geoalert.GeoAlert
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, respErr)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(requestBody, &geoRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, respErr)
		return
	}

	geoRequest.UserId = userId
	result, createErr := services.GeoAlertService.Create(geoRequest)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
		return
	}
	http_utils.ResponseJson(w, http.StatusCreated, result)
	return

}

func (c *geoalertsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	geoalertId := strings.TrimSpace(vars["id"])
	geoa, err := services.GeoAlertService.Get(geoalertId)

	if err != nil {
		http_utils.ResponseError(w, err)
		return
	}
	http_utils.ResponseJson(w, http.StatusOK, geoa)
}
func (c *geoalertsController) GetUserGeoAlerts(w http.ResponseWriter, r *http.Request) {
	authErr := http_utils.AuthenticateRequest(r, false)
	if authErr != nil {
		http_utils.ResponseError(w, authErr)
	}
	vars := mux.Vars(r)
	UserId, converr := strconv.Atoi(strings.TrimSpace(vars["user_id"]))
	if converr != nil {
		http_utils.ResponseError(w, rest_errors.NewBadRequestError("user_id params must be an integer"))
		return
	}
	geoa, err := services.GeoAlertService.GetUserGeoAlerts(int64(UserId), oauth.IsPublic(r))

	if err != nil {
		http_utils.ResponseError(w, err)
		return
	}
	http_utils.ResponseJson(w, http.StatusOK, geoa)
}

func (c *geoalertsController) Search(w http.ResponseWriter, r *http.Request) {
	authErr := http_utils.AuthenticateRequest(r, false)
	if authErr != nil {
		http_utils.ResponseError(w, authErr)
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiErr := rest_errors.NewBadRequestError("Invalid json body")
		http_utils.ResponseError(w, apiErr)
		return
	}
	defer r.Body.Close()

	var query queries.EsQuery
	if err := json.Unmarshal(bytes, &query); err != nil {
		apiErr := rest_errors.NewBadRequestError("invalid json query body")
		http_utils.ResponseError(w, apiErr)
		return
	}

	geoalerts, restErr := services.GeoAlertService.Search(query, oauth.IsPublic(r))
	if restErr != nil {
		http_utils.ResponseError(w, restErr)
		return
	}
	http_utils.ResponseJson(w, http.StatusOK, geoalerts)
}

func (cont *geoalertsController) Delete(w http.ResponseWriter, r *http.Request) {
	authErr := http_utils.AuthenticateRequest(r, true)
	if authErr != nil {
		http_utils.ResponseError(w, authErr)
		return
	}
	vars := mux.Vars(r)
	geoalertId := strings.TrimSpace(vars["id"])

	geoa, err := services.GeoAlertService.Get(geoalertId)

	if err != nil {
		http_utils.ResponseError(w, err)
		return
	}

	forceOwnerErr := http_utils.RestrictRequestToResourceOwner(r, geoa.UserId)
	if forceOwnerErr != nil {
		http_utils.ResponseError(w, forceOwnerErr)
		return
	}

	_, delErr := services.GeoAlertService.Delete(geoalertId)
	if delErr != nil {
		http_utils.ResponseError(w, delErr)
		return
	}
	http_utils.ResponseJson(w, http.StatusNoContent, nil)
}

func (cont *geoalertsController) Upsert(w http.ResponseWriter, r *http.Request) {
	authErr := http_utils.AuthenticateRequest(r, true)
	if authErr != nil {
		http_utils.ResponseError(w, authErr)
		return
	}

	vars := mux.Vars(r)
	geoalertId := strings.TrimSpace(vars["id"])

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, respErr)
		return
	}
	defer r.Body.Close()

	var geoalertRequest geoalert.GeoAlert
	if err := json.Unmarshal(requestBody, &geoalertRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid geoalert json body")
		http_utils.ResponseError(w, respErr)
		return
	}
	forceOwnerErr := http_utils.RestrictRequestToResourceOwner(r, geoalertRequest.UserId)
	if forceOwnerErr != nil {
		http_utils.ResponseError(w, forceOwnerErr)
		return
	}
	result, createErr := services.GeoAlertService.Upsert(geoalertRequest, geoalertId)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
		return
	}
	http_utils.ResponseJson(w, http.StatusCreated, result)
}
