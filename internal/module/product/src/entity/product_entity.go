package entity

type Product struct {
	ID           string `gorm:"column:id;primaryKey;type:varchar(36)"`
	StoreID      string `gorm:"column:store_id;type:varchar(36);not null"`
	ProductName  string `gorm:"column:product_name;type:varchar(255);not null"`
	CostPrice    int64  `gorm:"column:cost_price;type:bigint;not null;default:0"`
	SellingPrice int64  `gorm:"column:selling_price;type:bigint;not null;default:0"`
	Stock        int    `gorm:"column:stock;type:int;not null;default:0"`
	Unit         string `gorm:"column:unit;type:varchar(50);not null"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Product) TableName() string {
	return "products"
}