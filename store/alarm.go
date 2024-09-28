package store

import (
	"errors"
	"gorm.io/gorm"
)

type AlarmInfo struct {
	gorm.Model
	EoaOwner                   string `gorm:"type:VARCHAR(64); uniqueIndex" json:"eoa_owner"`
	AlarmType                  int    `json:"alarm_type"`
	AlarmChannel               string `json:"alarm_channel"`
	AlarmChannelHash           string `json:"alarm_channel_hash"`
	ReportLiquidationThreshold uint64 `json:"report_liquidation_threshold"`
	ReportOperatorFeeChange    bool   `json:"report_operator_fee_change"`
	ReportNetworkFeeChange     bool   `json:"report_network_fee_change"`
	ReportProposeBlock         bool   `json:"report_propose_block"`
	ReportMissedBlock          bool   `json:"report_missed_block"`
	ReportBalanceDecrease      bool   `json:"report_balance_decrease"`
	ReportExitedButNotRemoved  bool   `json:"report_exited_but_not_removed"`
	ReportWeekly               bool   `json:"report_weekly"`
}

func (s *AlarmInfo) TableName() string {
	return "alarm_infos"
}

func (s *Store) GetAllAlarmInfos() ([]AlarmInfo, error) {
	var alarmInfos []AlarmInfo
	err := s.db.Model(&AlarmInfo{}).Find(&alarmInfos).Error
	if err != nil {
		return nil, err
	}
	return alarmInfos, nil
}

// DeleteAlarmByEoaOwner Unscoped delete
func (s *Store) DeleteAlarmByEoaOwner(eoaOwner string) error {
	return s.db.Model(&AlarmInfo{}).Unscoped().Where("eoa_owner = ?", eoaOwner).Delete(&AlarmInfo{}).Error
}

func (s *Store) GetAlarmByEoaOwner(eoaOwner string) (*AlarmInfo, error) {
	var alarmInfo AlarmInfo
	err := s.db.Model(&AlarmInfo{}).Where(&AlarmInfo{EoaOwner: eoaOwner}).First(&alarmInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &alarmInfo, nil
}

func (s *Store) CreateOrUpdateAlarmInfo(info *AlarmInfo) error {
	var alarmInfo AlarmInfo
	err := s.db.Model(&AlarmInfo{}).Where(&AlarmInfo{EoaOwner: info.EoaOwner}).First(&alarmInfo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(info).Error
	}
	if err != nil {
		return err
	}

	alarmInfo.AlarmType = info.AlarmType
	alarmInfo.AlarmChannel = info.AlarmChannel         // Encrypted
	alarmInfo.AlarmChannelHash = info.AlarmChannelHash // Unencrypted AlarmChannel's hash
	alarmInfo.ReportLiquidationThreshold = info.ReportLiquidationThreshold
	alarmInfo.ReportOperatorFeeChange = info.ReportOperatorFeeChange
	alarmInfo.ReportNetworkFeeChange = info.ReportNetworkFeeChange
	alarmInfo.ReportProposeBlock = info.ReportProposeBlock
	alarmInfo.ReportMissedBlock = info.ReportMissedBlock
	alarmInfo.ReportBalanceDecrease = info.ReportBalanceDecrease
	alarmInfo.ReportExitedButNotRemoved = info.ReportExitedButNotRemoved
	alarmInfo.ReportWeekly = info.ReportWeekly
	return s.db.Save(alarmInfo).Error
}
