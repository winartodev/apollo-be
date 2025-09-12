package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winartodev/apollo-be/helper"
	"github.com/winartodev/apollo-be/infrastructure/http/response"
	"github.com/winartodev/apollo-be/modules/country/delivery/dto"
)

const (
	countriesApiURL = "https://www.apicountries.com/countries"
)

type CountryHandler struct {
}

func NewCountryHandler() *CountryHandler {
	return &CountryHandler{}
}

func (h *CountryHandler) Countries(e echo.Context) error {
	req, err := http.NewRequest("GET", countriesApiURL, nil)
	if err != nil {
		return response.FailedResponse(e, http.StatusBadRequest, err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response.FailedResponse(e, http.StatusBadRequest, err)
	}

	defer resp.Body.Close()

	var apiResponse []struct {
		Name         string   `json:"name"`
		Alpha3Code   string   `json:"alpha3Code"`
		Flags        dto.Flag `json:"flags"`
		CallingCodes []string `json:"callingCodes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return response.FailedResponse(e, http.StatusBadRequest, err)
	}

	data := make([]dto.CountryResponse, len(apiResponse))
	for i, country := range apiResponse {
		data[i] = dto.CountryResponse{
			Name:         country.Name,
			Code:         country.Alpha3Code,
			Flags:        country.Flags,
			CallingCodes: helper.GetFirstElement(country.CallingCodes),
		}
	}

	return response.SuccessResponse(e, http.StatusOK, "", data, nil)
}

func (h *CountryHandler) RegisterRoutes(api *echo.Group) error {
	api.GET("/countries", h.Countries)

	return nil
}
