package link

import (
	"go/adv-demo/pkg/db"
	"gorm.io/gorm/clause"
	"log"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}
func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.Database.DB.Where("hash = ?", hash).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil

}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) FindById(id uint) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.Table("links").
		Where("deleted_at is null").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset uint) ([]Link, error) {
	var links []Link
	err := repo.Database.Table("links").
		Find(&links, limit, offset).
		Where("deleted_at is null").
		Order("id").
		Scan(&links)
	if err != nil {
		log.Println(err)
	}
	return links, nil
}
