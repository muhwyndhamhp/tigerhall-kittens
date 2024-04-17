package tiger

import (
	"context"

	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/scopes"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

// Create implements entities.TigerRepository.
func (r *repo) Create(ctx context.Context, tiger *entities.Tiger) error {
	err := r.db.WithContext(ctx).Create(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements entities.TigerRepository.
func (r *repo) FindAll(ctx context.Context, page, pageSize int) ([]entities.Tiger, int, error) {
	var res []entities.Tiger
	var count int64

	q := r.db.
		WithContext(ctx).
		Model(&entities.Tiger{})

	err := q.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = q.
		Scopes(scopes.Paginate(page, pageSize)).
		Order("last_seen DESC").
		Find(&res).
		Error
	if err != nil {
		return nil, 0, err
	}

	return res, int(count), nil
}

// FindByID implements entities.TigerRepository.
func (r *repo) FindByID(ctx context.Context, id uint) (*entities.Tiger, error) {
	var res entities.Tiger
	err := r.db.WithContext(ctx).First(&res, id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Update implements entities.TigerRepository.
func (r *repo) Update(ctx context.Context, tiger *entities.Tiger, id uint) error {
	if tiger.ID == 0 {
		tiger.ID = id
	}

	err := r.db.WithContext(ctx).Save(tiger).Error
	if err != nil {
		return err
	}

	return nil
}

func NewTigerRepository(db *gorm.DB) entities.TigerRepository {
	return &repo{db}
}
