package repository

type LandpadRepository interface {
	IsExists(id string) (bool, error)
	IsActive(id string) (bool, error)
}
