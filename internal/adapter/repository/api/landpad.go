package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
)

type LandpadRepository struct {
	client *resty.Client
}

const landpadActiveStatus = "active"

func NewLandpadRepository(apiBaseURL string) *LandpadRepository {
	return &LandpadRepository{
		client: resty.New().
			SetBaseURL(apiBaseURL).
			SetHeader("Content-Type", "application/json"),
	}
}

func (r *LandpadRepository) IsExists(id string) (bool, error) {
	resp, err := r.client.R().
		SetPathParam("landpadId", id).
		Get("/v4/landpads/{landpadId}")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode())
	}
	return resp.StatusCode() == http.StatusOK, nil
}

func (r *LandpadRepository) IsActive(id string) (bool, error) {
	resp, err := r.client.R().
		SetPathParam("landpadId", id).
		Get("/v4/landpads/{landpadId}")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode())
	}

	respData := gjson.ParseBytes(resp.Body())
	return respData.Get("Status").String() == landpadActiveStatus, nil
}
