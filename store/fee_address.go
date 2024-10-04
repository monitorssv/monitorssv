package store

import (
	"errors"
	"gorm.io/gorm"
)

type FeeAddressInfo struct {
	gorm.Model
	Owner      string `gorm:"type:VARCHAR(64); uniqueIndex" json:"owner"`
	FeeAddress string `json:"fee_address"`
}

func (s *FeeAddressInfo) TableName() string {
	return "fee_address_infos"
}

func (s *Store) CreateOrUpdateClusterFeeAddress(info *FeeAddressInfo) error {
	var feeAddress FeeAddressInfo
	err := s.db.Where(&FeeAddressInfo{Owner: info.Owner}).First(&feeAddress).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(info).Error
	}
	if err != nil {
		return err
	}

	feeAddress.FeeAddress = info.FeeAddress
	return s.db.Save(&feeAddress).Error
}

func (s *Store) GetClusterFeeAddress(owner string) (FeeAddressInfo, error) {
	var feeAddress FeeAddressInfo
	err := s.db.Model(&FeeAddressInfo{}).Where(&FeeAddressInfo{Owner: owner}).First(&feeAddress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return FeeAddressInfo{}, nil
	}

	if err != nil {
		return FeeAddressInfo{}, err
	}

	return feeAddress, nil
}
