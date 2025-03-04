package store

import (
	"errors"
	"gorm.io/gorm"
)

type NetworkInfo struct {
	gorm.Model
	UpcomingNetworkFee uint64 `json:"upcoming_network_fee"`
	CurrentNetworkFee  uint64 `json:"current_network_fee"`
}

func (s *NetworkInfo) TableName() string {
	return "network_infos"
}

func (s *Store) GetNetworkInfo() (*NetworkInfo, error) {
	var info NetworkInfo
	err := s.db.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *Store) UpdateCurrentNetworkFee(fee uint64) error {
	var info NetworkInfo
	err := s.db.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(&NetworkInfo{CurrentNetworkFee: fee}).Error
	}
	if err != nil {
		return err
	}
	info.CurrentNetworkFee = fee
	return s.db.Save(&info).Error
}

func (s *Store) UpdateUpcomingNetworkFee(fee uint64) error {
	var info NetworkInfo
	err := s.db.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(&NetworkInfo{UpcomingNetworkFee: fee}).Error
	}
	if err != nil {
		return err
	}
	info.UpcomingNetworkFee = fee
	return s.db.Save(&info).Error
}
