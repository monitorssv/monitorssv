package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strings"
)

const DefaultValidatorIndex = int64(-1)

const (
	PendingInitialized = "pending_initialized" // Pending
	PendingQueued      = "pending_queued"      // Pending
	ActiveOngoing      = "active_ongoing"      // Active
	ActiveExiting      = "active_exiting"      // Active
	ActiveSlashed      = "active_slashed"      // Slashed
	ExitedUnSlashed    = "exited_unslashed"    // Exited
	ExitedSlashed      = "exited_slashed"      // Slashed
	WithdrawalPossible = "withdrawal_possible" // Withdrawal
	WithdrawalDone     = "withdrawal_done"     // Withdrawal
)

const (
	ValidatorPending = "Pending"
	ValidatorActive  = "Active"
	ValidatorSlashed = "Slashed"
	ValidatorExited  = "Exited"
	ValidatorUnknown = "Unknown"
)

func GetStatusDescription(status string) string {
	switch status {
	case PendingInitialized, PendingQueued:
		return ValidatorPending
	case ActiveOngoing, ActiveExiting:
		return ValidatorActive
	case ActiveSlashed, ExitedSlashed:
		return ValidatorSlashed
	case ExitedUnSlashed, WithdrawalPossible, WithdrawalDone:
		return ValidatorExited
	default:
		return ValidatorUnknown
	}
}

type ValidatorInfo struct {
	gorm.Model
	ClusterID         string `gorm:"type:VARCHAR(64); uniqueIndex:clusterid_pubkey_registration_block; index" json:"cluster_id"`
	Owner             string `gorm:"type:VARCHAR(64); index" json:"owner"`
	OperatorIds       string `json:"operator_ids"`
	PublicKey         string `gorm:"type:VARCHAR(255); uniqueIndex:clusterid_pubkey_registration_block; index" json:"public_key"`
	ValidatorIndex    int64  `gorm:"index" json:"validator_index"`
	RegistrationBlock int64  `gorm:"uniqueIndex:clusterid_pubkey_registration_block; index" json:"registration_block"`
	RemoveBlock       int64  `gorm:"index" json:"remove_block"`
	ExitedBlock       int64  `gorm:"index" json:"exited_block"`
	IsSlashed         bool   `gorm:"default:false" json:"is_slashed"`
	IsOnline          bool   `gorm:"default:true" json:"is_online"`
	Status            string `json:"status"`
}

func (s *ValidatorInfo) TableName() string {
	return "validator_infos"
}

func (s *Store) getValidatorChanges() ([]BlockChange, error) {
	var changes []BlockChange

	err := s.db.Model(&ValidatorInfo{}).
		Select("registration_block as block, count(*) as `change`").
		Group("registration_block").
		Find(&changes).Error
	if err != nil {
		return nil, err
	}
	var removals []BlockChange
	err = s.db.Model(&ValidatorInfo{}).
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

func (s *Store) GetLatestValidators() ([]ValidatorInfo, error) {
	var validators []ValidatorInfo
	err := s.db.Model(&ValidatorInfo{}).Order("id DESC").Where("remove_block = 0").Limit(10).Find(&validators).Error
	if err != nil {
		return nil, err
	}
	return validators, nil
}

func (s *Store) GetActiveValidatorCount() (int64, error) {
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where("remove_block = 0").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetActiveButExitedValidatorCount(clusterId string) (int64, error) {
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ClusterID: clusterId}).Where("remove_block = 0 AND exited_block != 0").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetValidatorByPublicKey(publicKey string) (*ValidatorInfo, error) {
	var validator ValidatorInfo
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey}).First(&validator).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &validator, nil
}

func (s *Store) GetValidatorByValidatorIndex(validatorIndex int64) (*ValidatorInfo, error) {
	var validator ValidatorInfo
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ValidatorIndex: validatorIndex}).First(&validator).Error
	if err != nil {
		return nil, err
	}
	return &validator, nil
}

func (s *Store) GetValidatorByPubKeyAndBlock(publicKey string, block uint64) (*ValidatorInfo, error) {
	var validators []ValidatorInfo
	query := s.db.Model(&ValidatorInfo{}).Where("public_key = ? AND registration_block < ?", publicKey, block)
	query = query.Where("(remove_block = 0 OR remove_block > ?) AND (exited_block = 0 OR exited_block > ?)", block, block)

	err := query.Find(&validators).Error
	if err != nil {
		return nil, err
	}

	if len(validators) == 0 {
		return nil, nil
	}

	if len(validators) > 1 {
		return nil, fmt.Errorf("validator exist in multiple active clusters, pubkey: %s, Block: %d", publicKey, block)
	}
	return &validators[0], nil
}

func (s *Store) GetValidators(page int, itemsPerPage int) ([]ValidatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where("remove_block = 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var validators []ValidatorInfo
	err = s.db.Model(&ValidatorInfo{}).Where("remove_block = 0").Order("id DESC").Offset(offset).Limit(perPage).Find(&validators).Error
	if err != nil {
		return nil, 0, err
	}
	return validators, totalCount, nil
}

func (s *Store) AdminGetValidators(page int, itemsPerPage int) ([]ValidatorInfo, int64, error) {
	perPage, offset := adminPagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where("remove_block = 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var validators []ValidatorInfo
	err = s.db.Model(&ValidatorInfo{}).Where("remove_block = 0").Order("id ASC").Offset(offset).Limit(perPage).Find(&validators).Error
	if err != nil {
		return nil, 0, err
	}
	return validators, totalCount, nil
}

func (s *Store) GetActiveButExitedValidators(page int, itemsPerPage int) ([]ValidatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64

	err := s.db.Model(&ValidatorInfo{}).Where("remove_block = 0 AND exited_block != 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	var validators []ValidatorInfo
	err = s.db.Model(&ValidatorInfo{}).Where("remove_block = 0 AND exited_block != 0").Offset(offset).Limit(perPage).Find(&validators).Error
	if err != nil {
		return nil, 0, err
	}
	return validators, totalCount, nil
}

func (s *Store) GetActiveButExitedValidatorsByClusterId(clusterId string) ([]ValidatorInfo, error) {
	var validators []ValidatorInfo
	err := s.db.Model(&ValidatorInfo{}).Where("remove_block = 0 AND exited_block != 0 AND cluster_id = ?", clusterId).Find(&validators).Error
	if err != nil {
		return nil, err
	}
	return validators, nil
}

func (s *Store) GetClusterOfflineValidatorCount(clusterId string) (int64, error) {
	var totalOfflineCount int64
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ClusterID: clusterId}).Where("remove_block = 0 AND is_online = 0").Count(&totalOfflineCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return totalOfflineCount, nil
}

func (s *Store) GetValidatorByClusterId(page int, itemsPerPage int, clusterId string) ([]ValidatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ClusterID: clusterId}).Where("remove_block = 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var validators []ValidatorInfo
	err = s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ClusterID: clusterId}).Where("remove_block = 0").Order("exited_block DESC, id DESC").Offset(offset).Limit(perPage).Find(&validators).Error
	if err != nil {
		return nil, 0, err
	}
	return validators, totalCount, nil
}

// GetAllValidatorsByClusterId should check  is_slashed == false
func (s *Store) GetAllValidatorsByClusterId(clusterId string) ([]ValidatorInfo, error) {
	var validators []ValidatorInfo
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ClusterID: clusterId}).Where("remove_block = 0").Find(&validators).Error
	if err != nil {
		return nil, err
	}
	return validators, nil
}

func (s *Store) GetValidatorByOwner(page int, itemsPerPage int, owner string) ([]ValidatorInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{Owner: owner}).Where("remove_block = 0").Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	var validators []ValidatorInfo
	err = s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{Owner: owner}).Where("remove_block = 0").Order("id DESC").Offset(offset).Limit(perPage).Find(&validators).Error
	if err != nil {
		return nil, 0, err
	}
	return validators, totalCount, nil
}

func (s *Store) CreateValidator(info *ValidatorInfo) error {
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

func (s *Store) UpdateValidatorStatus(publicKey string, status string) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey}).Where("remove_block = 0").Update("status", status).Error
}

func (s *Store) RemoveValidator(publicKey, clusterId string, removeBlock int64) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey, ClusterID: clusterId}).Update("remove_block", removeBlock).Error
}

func (s *Store) ExitValidator(publicKey string, exitBlock int64) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey}).Where("remove_block = 0").Update("exited_block", exitBlock).Error
}

func (s *Store) ValidatorSlash(publicKey string) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey}).Where("remove_block = 0").Update("is_slashed", true).Error
}

func (s *Store) UpdateValidatorOnlineStatus(validatorIndex int64, status bool) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{ValidatorIndex: validatorIndex}).Update("is_online", status).Error
}

func (s *Store) UpdateValidatorOnlineStatusByPubKey(pubKey string, status bool) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: pubKey}).Update("is_online", status).Error
}

func (s *Store) UpdateValidatorIndex(publicKey string, index int64) error {
	return s.db.Model(&ValidatorInfo{}).Where(&ValidatorInfo{PublicKey: publicKey}).Where("remove_block = 0 AND validator_index = -1").Update("validator_index", index).Error
}
