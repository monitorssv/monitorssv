package store

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

type EventInfo struct {
	gorm.Model
	BlockNumber uint64 `gorm:"block_number"`
	Owner       string `gorm:"type:VARCHAR(64); index" json:"owner"`
	TxHash      string `gorm:"type:VARCHAR(70); uniqueIndex:txhash_logindex" json:"tx_hash"`
	LogIndex    uint   `gorm:"uniqueIndex:txhash_logindex" json:"log_index"`
	Action      string `json:"action"`
	ClusterID   string `gorm:"index" json:"cluster_id"`
}

func (s *EventInfo) TableName() string {
	return "event_infos"
}

func (s *Store) TxHashLogIndexIsExist(txHash string, logIndex uint) bool {
	var eventInfo EventInfo
	result := s.db.Select("id").Where("tx_hash = ? AND log_index = ?", txHash, logIndex).First(&eventInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (s *Store) GetLatestEvents() ([]EventInfo, error) {
	var events []EventInfo
	err := s.db.Model(&EventInfo{}).Order("id DESC").Limit(10).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Store) GetEventByClusterId(page int, itemsPerPage int, clusterID string) ([]EventInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&EventInfo{}).Where(&EventInfo{ClusterID: clusterID}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	var events []EventInfo
	err = s.db.Model(&EventInfo{}).Where(&EventInfo{ClusterID: clusterID}).Order("id DESC").Offset(offset).Limit(perPage).Find(&events).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	return events, totalCount, nil
}

func (s *Store) GetEventByAccount(page int, itemsPerPage int, owner string) ([]EventInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&EventInfo{}).Where(&EventInfo{Owner: owner}).Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	var events []EventInfo
	err = s.db.Model(&EventInfo{}).Where(&EventInfo{Owner: owner}).Order("id DESC").Offset(offset).Limit(perPage).Find(&events).Error
	if err != nil {
		return nil, 0, err
	}
	return events, totalCount, nil
}

func (s *Store) CreateEvent(events []*EventInfo) error {
	err := s.db.Create(events).Error
	if err == nil {
		return nil
	}

	// maybe rescan event
	if strings.Contains(err.Error(), "Duplicate entry") {
		return nil
	}
	return err
}
