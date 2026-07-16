package entity

type Store struct {
	ID        string `gorm:"column:id;primaryKey"`
	UserID    string `gorm:"column:user_id"`
	StoreName string `gorm:"column:store_name"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (s *Store) TableName() string {
	return "stores"
}
