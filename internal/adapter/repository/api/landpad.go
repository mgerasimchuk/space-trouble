package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LandpadRepository struct {
	apiBaseUri string
}

const landpadActiveStatus = "active"

func NewLandpadRepository(apiBaseUri string) *LandpadRepository {
	return &LandpadRepository{
		apiBaseUri: apiBaseUri,
	}
}

func (r *LandpadRepository) IsExists(id string) (bool, error) {
	_, statusCode, err := r.getLandpad(id)
	if err != nil {
		return false, err
	}

	return statusCode == http.StatusOK, nil
}

func (r *LandpadRepository) IsActive(id string) (bool, error) {
	lp, _, err := r.getLandpad(id)
	if err != nil {
		return false, err
	}

	return lp.Status == landpadActiveStatus, nil
}

func (r *LandpadRepository) getLandpad(id string) (lp *landpad, statusCode int, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/v4/landpads/%s", r.apiBaseUri, id))

	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return nil, resp.StatusCode, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	lp = new(landpad)
	err = json.Unmarshal(body, lp)
	if err != nil {
		return lp, resp.StatusCode, err
	}

	return lp, resp.StatusCode, nil
}

type landpad struct {
	Status string
}
