package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LaunchpadRepository struct {
	apiBaseUri string
}

const launchpadActiveStatus = "active"

func NewLaunchpadRepository(apiBaseUri string) *LaunchpadRepository {
	return &LaunchpadRepository{apiBaseUri: apiBaseUri}
}

func (r *LaunchpadRepository) IsExists(id string) (bool, error) {
	_, statusCode, err := r.getLaunchpad(id)
	if err != nil {
		return false, err
	}

	return statusCode == http.StatusOK, nil
}

func (r *LaunchpadRepository) IsActive(id string) (bool, error) {
	lp, _, err := r.getLaunchpad(id)
	if err != nil {
		return false, err
	}

	return lp.Status == launchpadActiveStatus, nil
}

func (r *LaunchpadRepository) IsDateAvailableForLaunch(launchpadID string, launchDate time.Time) (bool, error) {
	requestBody := fmt.Sprintf(`
{
  "query": {
      "launchpad": "%s",
      "$and": [{"date_utc": {"$gte": "%s"}}, {"date_utc": {"$lt": "%s"}}]
  }
}
`, launchpadID, launchDate.Format(time.RFC3339), launchDate.Add(24*time.Hour).Format(time.RFC3339))

	resp, err := http.Post(fmt.Sprintf("%s/v4/launches/query", r.apiBaseUri), "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("got unexpected http status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	lqr := launchesQueryResult{}
	err = json.Unmarshal(body, &lqr)
	if err != nil {
		return false, err
	}

	return lqr.TotalDocs == 0, nil
}

func (r *LaunchpadRepository) getLaunchpad(id string) (lp *launchpad, statusCode int, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/v4/launchpads/%s", r.apiBaseUri, id))

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

	lp = new(launchpad)
	err = json.Unmarshal(body, lp)
	if err != nil {
		return lp, resp.StatusCode, err
	}

	return lp, resp.StatusCode, nil
}

type launchpad struct {
	Status string
}

type launchesQueryResult struct {
	TotalDocs int
}
