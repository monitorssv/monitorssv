package store

import (
	"errors"
	"fmt"
	"github.com/monitorssv/monitorssv/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sort"
)

type Store struct {
	db *gorm.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Store.User,
		cfg.Store.Pass,
		cfg.Store.Host,
		cfg.Store.DB,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if cfg.Store.LogMode == "silent" {
		db.Logger = logger.Default.LogMode(logger.Silent)
	} else if cfg.Store.LogMode == "error" || cfg.Store.LogMode == "err" {
		db.Logger = logger.Default.LogMode(logger.Error)
	} else if cfg.Store.LogMode == "warn" {
		db.Logger = logger.Default.LogMode(logger.Warn)
	} else if cfg.Store.LogMode == "info" {
		db.Logger = logger.Default.LogMode(logger.Info)
	}

	err = db.AutoMigrate(&AlarmInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&BlockInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&ClusterInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&FeeAddressInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&EventInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&OperatorInfo{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&ScanPoint{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&SSVReward{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&ValidatorInfo{})
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

type BlockChange struct {
	Block  int64
	Change int
}

type ChartData struct {
	Name       string `json:"name"`
	Validators int    `json:"validators"`
	Operators  int    `json:"operators"`
}

func (s *Store) CalculateChartData() ([]ChartData, error) {
	validatorChanges, err := s.getValidatorChanges()
	if err != nil {
		return nil, err
	}
	operatorChanges, err := s.getOperatorChanges()
	if err != nil {
		return nil, err
	}

	if len(validatorChanges) == 0 && len(operatorChanges) == 0 {
		return nil, errors.New("no validators or operators found")
	}

	allChanges := append(validatorChanges, operatorChanges...)
	sort.Slice(allChanges, func(i, j int) bool {
		return allChanges[i].Block < allChanges[j].Block
	})

	minBlock := allChanges[0].Block
	maxBlock := allChanges[len(allChanges)-1].Block

	interval := (maxBlock - minBlock) / 11
	var chartData []ChartData
	totalValidators, totalOperators := 0, 0
	validatorIndex, operatorIndex := 0, 0

	for i := int64(0); i <= 11; i++ {
		targetBlock := minBlock + i*interval
		if i == 11 {
			targetBlock = maxBlock
		}

		for validatorIndex < len(validatorChanges) && validatorChanges[validatorIndex].Block <= targetBlock {
			totalValidators += validatorChanges[validatorIndex].Change
			validatorIndex++
		}

		for operatorIndex < len(operatorChanges) && operatorChanges[operatorIndex].Block <= targetBlock {
			totalOperators += operatorChanges[operatorIndex].Change
			operatorIndex++
		}

		chartData = append(chartData, ChartData{
			Name:       fmt.Sprintf("%d", targetBlock),
			Validators: totalValidators,
			Operators:  totalOperators,
		})
	}

	return chartData, nil
}

func pagingCheck(page int, itemsPerPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if itemsPerPage < 1 || itemsPerPage > 100 {
		itemsPerPage = 10
	}
	offset := (page - 1) * itemsPerPage
	return itemsPerPage, offset
}

func adminPagingCheck(page int, itemsPerPage int) (int, int) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * itemsPerPage
	return itemsPerPage, offset
}
