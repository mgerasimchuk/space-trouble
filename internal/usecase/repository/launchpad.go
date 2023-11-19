package repository

import "time"

//go:generate mockgen -source=launchpad.go -destination=../../adapter/repository/mock/launchpad.go -package=mock

type LaunchpadRepository interface {
	IsExists(id string) (bool, error)
	IsActive(id string) (bool, error)
	IsDateAvailableForLaunch(launchpadID string, launchDate time.Time) (bool, error)
}
