package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

type LaunchpadRepository struct {
	client *resty.Client
}

const launchpadActiveStatus = "active"

func NewLaunchpadRepository(apiBaseURL string) *LaunchpadRepository {
	return &LaunchpadRepository{
		client: resty.New().
			SetBaseURL(apiBaseURL).
			SetHeader("Content-Type", "application/json"),
	}
}

func (r *LaunchpadRepository) IsExists(id string) (bool, error) {
	resp, err := r.client.R().
		SetPathParam("launchpadId", id).
		Get("/v4/launchpads/{launchpadId}")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode())
	}
	return resp.StatusCode() == http.StatusOK, nil
}

func (r *LaunchpadRepository) IsActive(id string) (bool, error) {
	resp, err := r.client.R().
		SetPathParam("launchpadId", id).
		Get("/v4/launchpads/{launchpadId}")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode())
	}

	respData := gjson.ParseBytes(resp.Body())
	return respData.Get("Status").String() == launchpadActiveStatus, nil
}

func (r *LaunchpadRepository) IsDateAvailableForLaunch(launchpadID string, launchDate time.Time) (bool, error) {
	resp, err := r.client.R().SetBody(map[string]any{
		"query": map[string]any{
			"launchpad": launchpadID,
			"$and": []map[string]any{
				{
					"date_utc": map[string]any{"$gte": launchDate.Format(time.RFC3339)},
				},
				{
					"date_utc": map[string]any{"$lt": launchDate.Add(24 * time.Hour).Format(time.RFC3339)},
				},
			},
		},
	}).Post("/v4/launches/query")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode())
	}

	respData := gjson.ParseBytes(resp.Body())
	return respData.Get("TotalDocs").Int() == 0, nil
}
