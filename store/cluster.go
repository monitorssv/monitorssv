package store

import (
	"errors"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"gorm.io/gorm"
	"math/big"
)

type ClusterInfo struct {
	gorm.Model
	ClusterID            string `gorm:"type:VARCHAR(64); uniqueIndex" json:"cluster_id"`
	Owner                string `gorm:"type:VARCHAR(64); index" json:"owner"`
	EoaOwner             string `gorm:"type:VARCHAR(64); index" json:"eoa_owner"`
	OperatorIds          string `json:"operator_ids"`
	ValidatorCount       uint32 `gorm:"index" json:"validator_count"`
	NetworkFeeIndex      uint64 `json:"network_fee_index"`
	Index                uint64 `json:"cluster_index"`
	Active               bool   `json:"active"`
	Balance              string `json:"balance"`
	BurnFee              uint64 `json:"burn_fee"`
	OnChainBalance       string `json:"on_chain_balance"`
	LiquidationBlock     uint64 `json:"liquidation_block"`
	CalcLiquidationBlock uint64 `json:"calc_liquidation_block"`
}

func CalcClusterOnChainBalance(curBlock uint64, clusterInfo *ClusterInfo) string {
	if clusterInfo.OnChainBalance == "" || !clusterInfo.Active {
		return "0"
	}

	preOnChainBalance, _ := big.NewInt(0).SetString(clusterInfo.OnChainBalance, 10)
	burnFeeRate := big.NewInt(0).SetUint64(clusterInfo.BurnFee)
	oneBlockBurnAmount := big.NewInt(0).Mul(burnFeeRate, big.NewInt(int64(clusterInfo.ValidatorCount)))
	burnAmount := big.NewInt(0).Mul(oneBlockBurnAmount, big.NewInt(int64(curBlock-clusterInfo.CalcLiquidationBlock)))
	curBalance := big.NewInt(0)
	if preOnChainBalance.Cmp(big.NewInt(0)) > 0 {
		curBalance = big.NewInt(0).Sub(preOnChainBalance, burnAmount)
	}

	return utils.ToSSV(curBalance, "%.2f")
}

func (s *ClusterInfo) TableName() string {
	return "cluster_infos"
}

func (s *Store) GetClusters(page int, itemsPerPage int) ([]ClusterInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64
	err := s.db.Model(&ClusterInfo{}).Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}
	var clusters []ClusterInfo
	err = s.db.Model(&ClusterInfo{}).Order("validator_count DESC, id ASC").Offset(offset).Limit(perPage).Find(&clusters).Error
	if err != nil {
		return nil, 0, err
	}
	return clusters, totalCount, nil
}

func (s *Store) GetAllClusters() ([]ClusterInfo, error) {
	var clusters []ClusterInfo
	err := s.db.Model(&ClusterInfo{}).Find(&clusters).Error
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (s *Store) GetActiveClusterCount() (int64, error) {
	var totalCount int64
	err := s.db.Model(&ClusterInfo{}).Where("active = 1").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *Store) GetActiveClusters() ([]ClusterInfo, error) {
	var clusters []ClusterInfo
	err := s.db.Model(&ClusterInfo{}).Where("validator_count != 0").Find(&clusters).Error
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (s *Store) GetAllClusterByEoaOwner(owner string) ([]ClusterInfo, error) {
	var clusters []ClusterInfo
	err := s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{EoaOwner: owner}).Find(&clusters).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (s *Store) GetClusterByOwner(page, itemsPerPage int, owner string) ([]ClusterInfo, int64, error) {
	perPage, offset := pagingCheck(page, itemsPerPage)
	var totalCount int64

	err := s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{Owner: owner}).Count(&totalCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	var clusters []ClusterInfo
	err = s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{Owner: owner}).Offset(offset).Limit(perPage).Find(&clusters).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	return clusters, totalCount, nil
}

func (s *Store) GetClusterByClusterId(clusterId string) (*ClusterInfo, error) {
	var cluster ClusterInfo
	err := s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{ClusterID: clusterId}).First(&cluster).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (s *Store) GetNoUpdatedClustersOwner() ([]string, error) {
	var owners []string
	err := s.db.Model(&ClusterInfo{}).Where("eoa_owner = '0x'").Distinct().Pluck("owner", &owners).Error
	if err != nil {
		return nil, err
	}

	return owners, nil
}

func (s *Store) CreateOrUpdateCluster(info *ClusterInfo) error {
	var cluster ClusterInfo
	err := s.db.Where(&ClusterInfo{ClusterID: info.ClusterID}).First(&cluster).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(info).Error
	}
	if err != nil {
		return err
	}
	if cluster.Owner != info.Owner || cluster.OperatorIds != info.OperatorIds {
		return errors.New("cluster owner or operatorIds is abnormal")
	}

	cluster.ValidatorCount = info.ValidatorCount
	cluster.NetworkFeeIndex = info.NetworkFeeIndex
	cluster.Index = info.Index
	cluster.Active = info.Active
	cluster.Balance = info.Balance
	return s.db.Save(&cluster).Error
}

func (s *Store) ClusterLiquidation(clusterID string, liquidationBlock uint64) error {
	return s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{ClusterID: clusterID}).Updates(map[string]interface{}{"liquidation_block": liquidationBlock, "calc_liquidation_block": liquidationBlock, "on_chain_balance": "0"}).Error
}

func (s *Store) UpdateClusterLiquidationInfo(clusterID string, liquidationBlock uint64, calculateLiquidationBlock uint64, burnFee uint64, onChainBalance string) error {
	return s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{ClusterID: clusterID}).Updates(map[string]interface{}{"liquidation_block": liquidationBlock, "burn_fee": burnFee, "calc_liquidation_block": calculateLiquidationBlock, "on_chain_balance": onChainBalance}).Error
}

func (s *Store) UpdateClusterStatus(clusterID string, status string) error {
	return s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{ClusterID: clusterID}).Update("status", status).Error
}

func (s *Store) UpdateClusterEoaOwner(owner string, eoaOwner string) error {
	return s.db.Model(&ClusterInfo{}).Where(&ClusterInfo{Owner: owner}).Update("eoa_owner", eoaOwner).Error
}
