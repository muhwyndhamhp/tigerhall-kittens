package tiger

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/models"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// Create implements models.TigerRepository.
func (r *repo) Create(tiger *models.Tiger) error {
	err := r.db.Create(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements models.TigerRepository.
func (r *repo) FindAll(page, pageSize int) ([]models.Tiger, error) {
	var res []models.Tiger
	err := r.db.Scopes(scopes.Paginate(page, pageSize)).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID implements models.TigerRepository.
func (r *repo) FindByID(id uint) (*models.Tiger, error) {
	var res models.Tiger
	err := r.db.First(&res, id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Update implements models.TigerRepository.
func (r *repo) Update(tiger *models.Tiger, id uint) error {
	if tiger.ID == 0 {
		tiger.ID = id
	}

	err := r.db.Save(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

func NewTigerRepository(db *gorm.DB) models.TigerRepository {
	return &repo{db}
}
