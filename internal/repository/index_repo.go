// Package repository contains the repository layer for the Moneybots API
package repository

import (
	"fmt"

	"github.com/nsvirk/moneybotsapi/internal/models"
	"gorm.io/gorm"
)

// Repository is the database repository for indices
type IndexRepository struct {
	DB *gorm.DB
}

// NewIndexRepository creates a new index repository
func NewIndexRepository(db *gorm.DB) *IndexRepository {
	return &IndexRepository{DB: db}
}

// TruncateIndices truncates the indices table
func (r *IndexRepository) TruncateIndices() error {
	return r.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", models.IndexTableName)).Error
}

// InsertIndices inserts a batch of indices into the database
func (r *IndexRepository) InsertIndices(indexInstruments []models.IndexModel) (int64, error) {
	// insert the records into the database
	result := r.DB.Create(indexInstruments)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to insert batch into %s: %v", models.IndexTableName, result.Error)
	}
	return result.RowsAffected, nil
}

// GetIndicesRecordCount returns the number of records in the indices table
func (r *IndexRepository) GetIndicesRecordCount() (int64, error) {
	var count int64
	err := r.DB.Table(models.IndexTableName).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get indices record count: %v", err)
	}
	return count, nil
}

// GetNSEIndexInstruments fetches the instruments for a given NSE index
func (r *IndexRepository) GetNSEIndexInstruments(indexName string) ([]models.IndexModel, error) {
	var indexInstruments []models.IndexModel
	err := r.DB.Where("index = ?", indexName).Find(&indexInstruments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch index instruments: %v", err)
	}

	return indexInstruments, nil
}

// GetNSEIndexNames fetches the names of all NSE indices
func (r *IndexRepository) GetNSEIndexNames() ([]string, error) {
	var indices []string
	err := r.DB.Table(models.IndexTableName).Select("DISTINCT index").Find(&indices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nse index names: %v", err)
	}
	return indices, nil
}