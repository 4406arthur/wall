package developer

import "wall/pkg/entity"

//Repository repository interface
type Repository interface {
	Find(name string) (*entity.Developer, error)
}
