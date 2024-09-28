package store

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

type BlockInfo struct {
	gorm.Model
	ClusterID   string `gorm:"type:VARCHAR(64); index" json:"cluster_id"`
	BlockNumber uint64 `json:"block_number"`
	Epoch       uint64 `json:"epoch"`
	Slot        uint64 `gorm:"uniqueIndex" json:"slot"`
	Proposer    uint64 `json:"proposer"`
	PublicKey   string `json:"public_key"`
	IsMissed    bool   `gorm:"index" json:"is_missed"`
}

func (s *BlockInfo) TableName() string {
	return "block_infos"
}

func (s *Store) CreateBlock(info *BlockInfo) error {
	err := s.db.Create(info).Error
	if err == nil {
		return nil
	}

	// maybe rescan event
	if strings.Contains(err.Error(), "Duplicate entry") {
		return nil
	}
	return err
}

func (s *Store) GetLatestBlocks() ([]BlockInfo, error) {
	var blocks []BlockInfo
	err := s.db.Model(&BlockInfo{}).Order("id DESC").Where("is_missed = 0").Limit(10).Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (s *Store) GetTotalBlockCount() (int64, error) {
	var totalCount int64
	err := s.db.Model(&BlockInfo{}).Where("is_missed = 0").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetBlockInfoByClusterId(clusterId string) (int64, int64, error) {
	var totalCount int64
	err := s.db.Model(&BlockInfo{}).Where("cluster_id = ? AND is_missed = 0", clusterId).Count(&totalCount).Error
	if err != nil {
		return 0, 0, err
	}
	var totalMissedCount int64
	err = s.db.Model(&BlockInfo{}).Where("cluster_id = ? AND is_missed = 1", clusterId).Count(&totalMissedCount).Error
	if err != nil {
		return 0, 0, err
	}

	return totalCount, totalMissedCount, nil
}

func (s *Store) GetValidatorTotalBlockCount(pubKey string) (int64, error) {
	var totalCount int64
	err := s.db.Model(&BlockInfo{}).Where(&BlockInfo{PublicKey: pubKey}).Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetBlockByClusterId(page int, itemsPerPage int, clusterID string) ([]BlockInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&BlockInfo{}).Where(&BlockInfo{ClusterID: clusterID}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	var blocks []BlockInfo
	err = s.db.Model(&BlockInfo{}).Where(&BlockInfo{ClusterID: clusterID}).Offset(offset).Limit(perPage).Find(&blocks).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	return blocks, totalCount, nil
}

func (s *Store) GetAllBlocks() ([]BlockInfo, error) {
	var blocks []BlockInfo
	err := s.db.Model(&BlockInfo{}).Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}
