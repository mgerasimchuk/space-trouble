package repository

import "time"

type LaunchpadRepository interface {
	IsExists(id string) (bool, error)
	IsActive(id string) (bool, error)
	IsDateAvailableForLaunch(launchpadID string, launchDate time.Time) (bool, error)
}
