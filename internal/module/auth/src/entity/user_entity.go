package entity

type User struct {
	ID        string `gorm:"column:id;primaryKey;type:varchar(36)"`
	Name      string `gorm:"column:name;type:varchar(100);not null;unique"`
	Email     string `gorm:"column:email;type:varchar(255);not null;unique"`
	Password  string `gorm:"column:password;type:varchar(255);not null"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (User) TableName() string {
	return "users"
}