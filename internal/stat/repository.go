package stat

import (
	"go/adv-demo/pkg/db"
	"gorm.io/datatypes"
	"time"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{db}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.Find(&stat, "link_id = ? and date = ? ", linkId, currentDate)
	if stat.ID == 0 {
		repo.Create(&Stat{
			LinkId: linkId,
			Date:   currentDate,
			Clicks: 1,
		})

	} else {
		stat.Clicks++
		repo.Save(&stat)
	}
}
