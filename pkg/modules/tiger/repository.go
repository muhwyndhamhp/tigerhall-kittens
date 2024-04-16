package tiger

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// Create implements entities.TigerRepository.
func (r *repo) Create(tiger *entities.Tiger) error {
	err := r.db.Create(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements entities.TigerRepository.
func (r *repo) FindAll(page, pageSize int) ([]entities.Tiger, error) {
	var res []entities.Tiger
	err := r.db.Scopes(scopes.Paginate(page, pageSize)).Order("last_seen DESC").Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID implements entities.TigerRepository.
func (r *repo) FindByID(id uint) (*entities.Tiger, error) {
	var res entities.Tiger
	err := r.db.First(&res, id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Update implements entities.TigerRepository.
func (r *repo) Update(tiger *entities.Tiger, id uint) error {
	if tiger.ID == 0 {
		tiger.ID = id
	}

	err := r.db.Save(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

func NewTigerRepository(db *gorm.DB) entities.TigerRepository {
	return &repo{db}
}
