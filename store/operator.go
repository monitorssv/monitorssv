package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strings"
)

type OperatorInfo struct {
	gorm.Model
	Owner                string `json:"owner"`
	OperatorId           uint64 `gorm:"uniqueIndex" json:"operator_id"`
	PubKey               string `json:"pub_key"`
	OperatorName         string `json:"operator_name"`
	ValidatorCount       uint32 `json:"validator_count"`
	ClusterIds           string `json:"cluster_ids"`
	OperatorFee          string `json:"operator_fee"`
	OperatorEarnings     string `gorm:"default:-" json:"operator_earnings"`
	PrivacyStatus        bool   `json:"privacy_status"`
	WhitelistedAddress   string `json:"whitelisted_address"`
	WhitelistingContract string `json:"whitelisting_contract"`
	RegistrationBlock    int64  `gorm:"index" json:"registration_block"`
	RemoveBlock          int64  `gorm:"index" json:"remove_block"`
	PendingOperatorFee   string `gorm:"default:0;index" json:"pending_operator_fee"`
	ApprovalBeginTime    uint64 `gorm:"default:0" json:"approval_begin_time"`
	ApprovalEndTime      uint64 `gorm:"default:0" json:"approval_end_time"`
}

func (s *OperatorInfo) TableName() string {
	return "operator_infos"
}

func (s *Store) getOperatorChanges() ([]BlockChange, error) {
	var changes []BlockChange
	err := s.db.Model(&OperatorInfo{}).
		Select("registration_block as block, count(*) as `change`").
		Group("registration_block").
		Find(&changes).Error
	if err != nil {
		return nil, err
	}

	var removals []BlockChange
	err = s.db.Model(&OperatorInfo{}).
		Select("remove_block as block, count(*) as `change`").
		Where("remove_block != 0").
		Group("remove_block").
		Find(&removals).Error
	if err != nil {
		return nil, err
	}

	for i := range removals {
		removals[i].Change = -removals[i].Change
	}

	changes = append(changes, removals...)

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Block < changes[j].Block
	})
	return changes, nil
}

func (s *Store) GetActiveOperatorCount() (int64, error) {
	var totalCount int64
	err := s.db.Model(&OperatorInfo{}).Where("remove_block = 0").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetMaxOperatorId() (uint64, error) {
	var operator OperatorInfo
	err := s.db.Model(&OperatorInfo{}).Last(&operator).Error
	if err != nil {
		return 0, err
	}
	return operator.OperatorId, nil
}

func (s *Store) GetOperators(page int, itemsPerPage int) ([]OperatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&OperatorInfo{}).Where("remove_block = 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	var operators []OperatorInfo
	err = s.db.Model(&OperatorInfo{}).Where("remove_block = 0").Order("validator_count DESC, operator_id ASC").Offset(offset).Limit(perPage).Find(&operators).Error
	if err != nil {
		return nil, 0, err
	}
	return operators, totalCount, nil
}

func (s *Store) GetOperatorByOperatorId(operatorId uint64) (*OperatorInfo, error) {
	var operator OperatorInfo
	err := s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).First(&operator).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &operator, nil
}

func (s *Store) GetOperatorByOperatorName(page, itemsPerPage int, operatorName string) ([]OperatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&OperatorInfo{}).Where("operator_name LIKE ?", fmt.Sprintf("%%%s%%", operatorName)).Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	var operators []OperatorInfo
	err = s.db.Model(&OperatorInfo{}).Where("operator_name LIKE ?", fmt.Sprintf("%%%s%%", operatorName)).Offset(offset).Limit(perPage).Find(&operators).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	return operators, totalCount, nil
}

func (s *Store) GetOperatorByOwner(page, itemsPerPage int, owner string) ([]OperatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{Owner: owner}).Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	var operators []OperatorInfo
	err = s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{Owner: owner}).Order("operator_id DESC").Offset(offset).Limit(perPage).Find(&operators).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	return operators, totalCount, nil
}

func (s *Store) CreateOperator(info *OperatorInfo) error {
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

func (s *Store) UpdateOperatorName(operatorId uint64, name string) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("operator_name", name).Error
}

func (s *Store) UpdateOperatorFee(operatorId uint64, fee string) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("operator_fee", fee).Error
}

func (s *Store) UpdatePendingOperator(operatorId uint64, pendingFee string, beginTime, endTime uint64) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Updates(map[string]interface{}{"pending_operator_fee": pendingFee, "approval_begin_time": beginTime, "approval_end_time": endTime}).Error
}

func (s *Store) CancelUpdateOperatorFee(operatorId uint64) error {
	return s.db.Model(&OperatorInfo{}).
		Where("operator_id = ?", operatorId).
		Updates(map[string]interface{}{
			"pending_operator_fee": "0",
			"approval_begin_time":  0,
			"approval_end_time":    0,
		}).Error
}

func (s *Store) UpdateOperatorEarning(operatorId uint64, earning string) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("operator_earnings", earning).Error
}

func (s *Store) UpdateOperatorPrivacyStatus(operatorId uint64, status bool) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("privacy_status", status).Error
}

func (s *Store) UpdateWhitelistedAddress(operatorId uint64, whitelistedAddress string) error {
	// multiple address, separated by commas
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("whitelisted_address", whitelistedAddress).Error
}

func (s *Store) UpdateWhitelistingContract(operatorId uint64, whitelistingContract string) error {
	// multiple address, separated by commas
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("whitelisting_contract", whitelistingContract).Error
}

func (s *Store) RemoveOperator(operatorId uint64, removeBlock int64) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("remove_block", removeBlock).Error
}

func (s *Store) UpdateOperatorValidatorCount(operatorId uint64, count uint32) error {
	return s.db.Model(&OperatorInfo{}).Where(&OperatorInfo{OperatorId: operatorId}).Update("validator_count", count).Error
}

func (s *Store) BatchUpdateOperatorValidatorCount(operatorIds []uint64, increment bool) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var updateSQL string
		if increment {
			updateSQL = "UPDATE operator_infos SET validator_count = validator_count + 1 WHERE operator_id IN ?"
		} else {
			updateSQL = "UPDATE operator_infos SET validator_count = CASE WHEN validator_count > 0 THEN validator_count - 1 ELSE 0 END WHERE operator_id IN ?"
		}

		result := tx.Exec(updateSQL, operatorIds)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (s *Store) GetOperatorValidatorCount(clusterIdsStr string) (uint32, error) {
	if clusterIdsStr == "" {
		return 0, nil
	}

	clusterIds := strings.Split(clusterIdsStr, ",")

	var totalCount uint32
	for _, clusterId := range clusterIds {
		cluster, err := s.GetClusterByClusterId(clusterId)
		if err != nil {
			return 0, err
		}
		if cluster.Active {
			totalCount += cluster.ValidatorCount
		}
	}
	return totalCount, nil
}

func (s *Store) BatchUpdateOperatorValidatorCounts(operatorIds []uint64, count uint32, increment bool) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var updateSQL string
		if increment {
			updateSQL = "UPDATE operator_infos SET validator_count = validator_count + ? WHERE operator_id IN ?"
		} else {
			updateSQL = "UPDATE operator_infos SET validator_count = CASE WHEN validator_count >= ? THEN validator_count - ? ELSE 0 END WHERE operator_id IN ?"
		}

		var result *gorm.DB
		if increment {
			result = tx.Exec(updateSQL, count, operatorIds)
		} else {
			result = tx.Exec(updateSQL, count, count, operatorIds)
		}

		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (s *Store) UpdateOperatorClusterIds(operatorId uint64, clusterId string, increment bool) error {
	var operator OperatorInfo
	if err := s.db.Where("operator_id = ?", operatorId).First(&operator).Error; err != nil {
		return err
	}

	if increment {
		if operator.ClusterIds == "" {
			operator.ClusterIds = clusterId
		} else if !strings.Contains(operator.ClusterIds, clusterId) {
			operator.ClusterIds += "," + clusterId
		} else {
			return nil
		}
	} else {
		operator.ClusterIds = strings.Replace(operator.ClusterIds, clusterId, "", -1)
		operator.ClusterIds = strings.Replace(operator.ClusterIds, ",,", ",", -1)
		operator.ClusterIds = strings.Trim(operator.ClusterIds, ",")
	}

	return s.db.Save(&operator).Error
}
