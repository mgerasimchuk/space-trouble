package repository

//go:generate mockgen -source=landpad.go -destination=mock/landpad.go -package=mock

type LandpadRepository interface {
	IsExists(id string) (bool, error)
	IsActive(id string) (bool, error)
}
