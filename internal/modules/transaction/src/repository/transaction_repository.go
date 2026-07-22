package repository

import (
	"github.com/PetaFlops-web/backend-shop-smbk/internal/modules/transaction/src/entity"
	"github.com/PetaFlops-web/backend-shop-smbk/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	repository.Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{Log: log}
}

// ListByStoreAndDate mengambil daftar transaksi milik sebuah toko pada tanggal tertentu
func (r *TransactionRepository) ListByStoreAndDate(db *gorm.DB, storeId string, date string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := db.Where("store_id = ? AND transaction_date = ?", storeId, date).Find(&transactions).Error
	return transactions, err
}