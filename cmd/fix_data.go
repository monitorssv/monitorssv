package main

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/eth1/ssv"
	"github.com/monitorssv/monitorssv/store"
	"github.com/urfave/cli/v2"
	"math/big"
	"strconv"
	"strings"
)

// The fix command cannot have side effects, even if it is executed multiple times.
var fixDBOperatorValidatorCountCmd = &cli.Command{
	Name:  "fix-db-operator-validator-count",
	Usage: "fix db operator validator count",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "conf-path",
			Usage: "config.yaml path",
			Value: "",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.InitConfig(ctx.String("conf-path"))
		if err != nil {
			log.Errorw("InitConfig", "err", err)
			return err
		}

		store, err := store.NewStore(cfg)
		if err != nil {
			log.Errorw("NewStore", "err", err)
			return err
		}

		eth1Client, err := client.NewEth1Client(cfg)
		if err != nil {
			log.Errorw("NewEth1Client", "err", err)
			return err
		}

		clusters, err := store.GetAllClusters()
		if err != nil {
			log.Errorw("GetAllClusters", "err", err)
			return err
		}

		for _, cluster := range clusters {
			if cluster.Active || cluster.ValidatorCount == 0 {
				continue
			}

			operatorIds, err := getOperatorIds(cluster.OperatorIds)
			if err != nil {
				log.Errorw("GetOperatorIds", "err", err)
				return err
			}

			var opIds = make([]uint64, 0)
			for _, operatorId := range operatorIds {
				operator, err := store.GetOperatorByOperatorId(operatorId)
				if err != nil {
					log.Errorw("GetOperatorByOperatorId", "err", err)
					return err
				}
				if strings.Contains(operator.ClusterIds, cluster.ClusterID) {
					opIds = append(opIds, operatorId)
				}
			}

			log.Infow("Will BatchUpdateOperatorValidatorCounts", "opIds", opIds, "cluster.OperatorIds", operatorIds, "cluster.ValidatorCount", cluster.ValidatorCount)

			if len(opIds) != 0 {
				err = store.BatchUpdateOperatorValidatorCounts(opIds, cluster.ValidatorCount, false)
				if err != nil {
					log.Errorw("BatchUpdateOperatorValidatorCounts", "err", err)
					return err
				}

				for _, operatorId := range opIds {
					log.Infow("Will UpdateOperatorClusterIds", "operatorId", operatorId, "clusterID", cluster.ClusterID)
					if err = store.UpdateOperatorClusterIds(operatorId, cluster.ClusterID, false); err != nil {
						return err
					}
				}
			}
		}

		maxOperatorId, err := store.GetMaxOperatorId()
		if err != nil {
			log.Warn("failed to get max operator id: %v", err)
			return err
		}
		contractInfo, err := ssv.GetSSVContractInfo(cfg.Network)
		if err != nil {
			log.Errorw("GetSSVContractInfo", "err", err)
			return err
		}
		ssvView, err := ssv.NewSsv(contractInfo.SSVNetworkView, eth1Client.GetClient())
		for i := uint64(1); i <= maxOperatorId; i++ {
			operator, err := store.GetOperatorByOperatorId(i)
			if err != nil {
				log.Errorw("GetOperatorByOperatorId", "err", err)
				return err
			}
			if operator.RemoveBlock != 0 {
				continue
			}

			_, _, chainValidatorCount, _, _, _, err := ssvView.GetOperatorById(nil, i)
			if err != nil {
				log.Errorw("GetOperatorById", "err", err)
				return err
			}
			if chainValidatorCount != operator.ValidatorCount {
				log.Warnw("chainValidatorCount != operator.ValidatorCount", "operatorId", i, "chainValidatorCount", chainValidatorCount, "operator.ValidatorCount", operator.ValidatorCount)
				{
					operator, err := store.GetOperatorByOperatorId(i)
					if err != nil {
						log.Errorw("GetOperatorByOperatorId", "err", err)
						return err
					}
					log.Warnw("abnormal operator", "Id", operator.OperatorId, "ValidatorCount", operator.ValidatorCount)
					clusterIds := strings.Split(operator.ClusterIds, ",")

					var totalCount uint32
					for _, clusterId := range clusterIds {
						cluster, err := store.GetClusterByClusterId(clusterId)
						if err != nil {
							log.Warnw("GetClusterByClusterId", "err", err)
							return err
						}
						if cluster.Active {
							totalCount += cluster.ValidatorCount
						}
					}
					log.Warnw("abnormal operator", "totalCount", totalCount)

					if totalCount == chainValidatorCount {
						log.Infow("update abnormal data", "Id", operator.OperatorId, "ValidatorCount", totalCount)
						err = store.UpdateOperatorValidatorCount(i, totalCount)
						if err != nil {
							log.Errorw("UpdateOperatorValidatorCount", "err", err)
							return err
						}
					} else {
						log.Warnw("abnormal operator:total cluster validatorCount does not match chainValidatorCount", "clusterTotalCount", totalCount, "chainValidatorCount", chainValidatorCount)
					}
				}
				continue
			}
			log.Infow("chainValidatorCount == operator.ValidatorCount", "operatorId", i)
		}

		return nil
	},
}

func getOperatorIds(operatorIdsStr string) ([]uint64, error) {
	operatorIds := make([]uint64, 0)
	for _, operatorIdStr := range strings.Split(operatorIdsStr, ",") {
		operatorId, err := strconv.ParseUint(operatorIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
		operatorIds = append(operatorIds, operatorId)
	}
	return operatorIds, nil
}

var feeRecipientAddressUpdatedTopic = common.HexToHash("0x259235c230d57def1521657e7c7951d3b385e76193378bc87ef6b56bc2ec3548")

// The fix command cannot have side effects, even if it is executed multiple times.
var fixDBClusterFeeAddressCmd = &cli.Command{
	Name:  "fix-db-cluster-fee-address",
	Usage: "fix db cluster fee address",
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "fromBlock",
			Usage: "fromBlock",
			Value: 0,
		},
		&cli.StringFlag{
			Name:  "conf-path",
			Usage: "config.yaml path",
			Value: "",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.InitConfig(ctx.String("conf-path"))
		if err != nil {
			log.Errorw("InitConfig", "err", err)
			return err
		}

		db, err := store.NewStore(cfg)
		if err != nil {
			log.Errorw("NewStore", "err", err)
			return err
		}

		eth1Client, err := client.NewEth1Client(cfg)
		if err != nil {
			log.Errorw("NewEth1Client", "err", err)
			return err
		}

		curBlock, err := eth1Client.BlockNumber()
		if err != nil {
			return err
		}

		contractInfo, err := ssv.GetSSVContractInfo(cfg.Network)
		if err != nil {
			log.Errorw("GetSSVContractInfo", "err", err)
			return err
		}

		events := ssv.GetAllSSVEvent()
		event := events[feeRecipientAddressUpdatedTopic]
		fromBlock := contractInfo.DeployBlock
		if ctx.Uint64("fromBlock") != 0 {
			fromBlock = ctx.Uint64("fromBlock")
		}
		for fromBlock < curBlock {
			nextBlock := fromBlock + 20000
			if nextBlock >= curBlock {
				nextBlock = curBlock
			}

			filter := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(fromBlock)),
				ToBlock:   big.NewInt(int64(nextBlock)),
				Addresses: []common.Address{contractInfo.SSVNetwork},
				Topics:    [][]common.Hash{{feeRecipientAddressUpdatedTopic}},
			}
			log.Infow("scan block", "fromBlock", fromBlock, "nextBlock", nextBlock)

			fromBlock = nextBlock
			logs, err := eth1Client.FilterLogs(filter)
			if err != nil {
				return err
			}

			for _, log := range logs {
				var owner common.Address
				copy(owner[:], log.Topics[1][12:])

				data, err := event.Inputs.Unpack(log.Data)
				if err != nil {
					return err
				}
				recipientAddress := data[0].(common.Address)
				if err = db.CreateOrUpdateClusterFeeAddress(&store.FeeAddressInfo{
					Owner:      owner.String(),
					FeeAddress: recipientAddress.String(),
				}); err != nil {
					return err
				}
			}
		}

		return nil
	},
}
