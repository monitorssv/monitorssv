package store

import (
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"math/big"
)

type SSVReward struct {
	gorm.Model
	MerkleRoot string          `json:"merkle_root"`
	Account    string          `gorm:"type:VARCHAR(64); uniqueIndex" json:"account"`
	Amount     decimal.Decimal `json:"amount"`
	Claimed    decimal.Decimal `json:"claimed"`
	Proofs     string          `json:"proofs"`
}

func (s *SSVReward) TableName() string {
	return "ssv_rewards"
}

func (s *Store) GetMerkleRoot() (string, error) {
	var reward SSVReward

	result := s.db.First(&reward)
	if result.Error != nil {
		return "", result.Error
	}

	return reward.MerkleRoot, nil
}

func (s *Store) GetAllSSVRewardAccount() ([]string, error) {
	var accounts []string
	result := s.db.Model(&SSVReward{}).Pluck("Account", &accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (s *Store) GetSSVReward(account string) SSVReward {
	var reward SSVReward
	// no reward return empty SSVReward
	s.db.Where(&SSVReward{Account: account}).First(&reward)
	reward.Account = account
	return reward
}

func (s *Store) CreateOrUpdateSSVReward(merkleRoot, account string, amountInt *big.Int, proofs string) error {
	amount := decimal.NewFromBigInt(amountInt, 0)
	var reward SSVReward
	err := s.db.Where(&SSVReward{Account: account}).First(&reward).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		reward = SSVReward{
			MerkleRoot: merkleRoot,
			Account:    account,
			Amount:     amount,
			Proofs:     proofs,
		}
		return s.db.Create(&reward).Error
	}
	if err != nil {
		return err
	}

	if !amount.GreaterThanOrEqual(reward.Claimed) {
		return errors.New("the reward received is greater than the total reward")
	}

	reward.MerkleRoot = merkleRoot
	reward.Amount = amount
	reward.Proofs = proofs
	return s.db.Save(&reward).Error
}

func (s *Store) AddClaimed(account string, amount *big.Int) error {
	var reward SSVReward
	err := s.db.Model(&SSVReward{}).Where(&SSVReward{Account: account}).First(&reward).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reward.Account = account
			reward.Claimed = decimal.NewFromBigInt(amount, 0)
			return s.db.Create(&reward).Error
		}
	}

	reward.Claimed = reward.Claimed.Add(decimal.NewFromBigInt(amount, 0))
	return s.db.Save(&reward).Error
}

func (s *Store) UpdateClaimed(account string, amount *big.Int) error {
	var reward SSVReward
	err := s.db.Model(&SSVReward{}).Where(&SSVReward{Account: account}).First(&reward).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reward.Account = account
			reward.Claimed = decimal.NewFromBigInt(amount, 0)
			return s.db.Create(&reward).Error
		}
	}

	reward.Claimed = decimal.NewFromBigInt(amount, 0)
	return s.db.Save(&reward).Error
}
