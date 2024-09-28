package store

import (
	"errors"
	"gorm.io/gorm"
)

type ScanPoint struct {
	gorm.Model
	Eth1Block uint64 `json:"eth1_block"`
	Eth2Slot  uint64 `json:"eth2_slot"`
}

func (s *ScanPoint) TableName() string {
	return "scan_points"
}

func (s *Store) GetScanPoint() (uint64, uint64, error) {
	var scanPoint ScanPoint
	err := s.db.First(&scanPoint).Error
	if err != nil {
		return 0, 0, err
	}
	return scanPoint.Eth1Block, scanPoint.Eth2Slot, nil
}

func (s *Store) UpdateScanEth1Block(block uint64) error {
	var scanPoint ScanPoint
	err := s.db.First(&scanPoint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(&ScanPoint{Eth1Block: block}).Error
	}
	if err != nil {
		return err
	}
	scanPoint.Eth1Block = block
	return s.db.Save(&scanPoint).Error
}

func (s *Store) UpdateScanEth2Slot(slot uint64) error {
	var scanPoint ScanPoint
	err := s.db.First(&scanPoint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(&ScanPoint{Eth2Slot: slot}).Error
	}
	if err != nil {
		return err
	}
	scanPoint.Eth2Slot = slot
	return s.db.Save(&scanPoint).Error
}
