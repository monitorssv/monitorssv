package eth2

import (
	"fmt"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/eth2/client"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"sort"
)

const EffectiveBalance = 32000000000

type Balance struct {
	Epoch  uint64
	Amount uint64
}

func (bm *BeaconMonitor) validatorMonitor(epoch uint64) error {
	slot := (epoch+1)*32 - 1
	itemsPerPage := 1000
	page := 1

	validatorInfoMap := make(map[string]*client.StandardValidatorEntry)
	clusterBalanceAlarms := make(map[string][]uint64)
	clusterSlashAlarms := make(map[string][]uint64)

	for {
		validators, totalCount, err := bm.store.AdminGetValidators(page, itemsPerPage)
		if err != nil {
			log.Errorw("AdminGetValidators", "err", err)
			return err
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(itemsPerPage)))
		log.Infow("validatorMonitor", "epoch", epoch, "page", page, "totalPages", totalPages, "itemsPerPage", itemsPerPage)

		var pubKeys []string
		var indexs []uint64
		validatorMap := make(map[string]*store.ValidatorInfo)

		for i := range validators {
			v := &validators[i]
			if v.IsSlashed || (v.Status != store.ValidatorActive && v.ExitedBlock != 0) {
				continue
			}

			pubKey := fmt.Sprintf("0x%s", v.PublicKey)
			validatorMap[pubKey] = v

			if v.ValidatorIndex != store.DefaultValidatorIndex {
				indexs = append(indexs, uint64(v.ValidatorIndex))
			} else {
				pubKeys = append(pubKeys, pubKey)
			}
		}

		if len(pubKeys) > 0 {
			log.Infow("GetSlotValidatorsByPubKey", "pubKeys", len(pubKeys))
			validatorInfoMap1, err := bm.client.GetSlotValidatorsByPubKey(slot, pubKeys)
			if err != nil {
				log.Warnw("GetSlotValidatorsByPubKey", "err", err)
				return err
			}

			for _, pubKey := range pubKeys {
				if info, ok := validatorInfoMap1[pubKey]; ok {
					log.Infow("UpdateValidatorIndex", "pubKey", pubKey, "index", info.Index)
					err = bm.store.UpdateValidatorIndex(removePubKeyPrefix(pubKey), int64(info.Index))
					if err != nil {
						log.Errorw("UpdateValidatorIndex", "err", err)
						return err
					}
				} else {
					err = bm.store.UpdateValidatorStatus(removePubKeyPrefix(pubKey), store.ValidatorUnknown)
					if err != nil {
						log.Errorw("UpdateValidatorStatus", "err", err)
						return err
					}
					err = bm.store.UpdateValidatorOnlineStatusByPubKey(removePubKeyPrefix(pubKey), false)
					if err != nil {
						log.Errorw("UpdateValidatorOnlineStatusByPubKey", "err", err, "pubKey", pubKey, "online", false)
						return err
					}
				}
			}

			for _, validatorInfo := range validatorInfoMap1 {
				log.Infow("UpdateValidatorIndex", "pubKey", validatorInfo.Validator.Pubkey, "index", validatorInfo.Index)
				err = bm.store.UpdateValidatorIndex(removePubKeyPrefix(validatorInfo.Validator.Pubkey), int64(validatorInfo.Index))
				if err != nil {
					log.Errorw("UpdateValidatorIndex", "err", err)
					return err
				}
			}
			mergeMaps(validatorInfoMap, validatorInfoMap1)
		}

		if len(indexs) > 0 {
			log.Infow("GetSlotValidatorsByIndex", "indexs", len(indexs))
			validatorInfoMap2, err := bm.client.GetSlotValidatorsByIndex(slot, indexs)
			if err != nil {
				log.Warnw("GetSlotValidatorsByIndex", "err", err)
				return err
			}
			mergeMaps(validatorInfoMap, validatorInfoMap2)
		}

		for _, validatorInfo := range validatorInfoMap {
			v := validatorMap[validatorInfo.Validator.Pubkey]
			if v == nil {
				continue
			}

			bm.updateBalanceHistory(uint64(validatorInfo.Index), epoch, uint64(validatorInfo.Balance))
			if isAlarm := bm.checkBalanceChange(uint64(validatorInfo.Index), v.IsOnline); isAlarm {
				clusterBalanceAlarms[v.ClusterID] = append(clusterBalanceAlarms[v.ClusterID], uint64(validatorInfo.Index))
			}

			if v.Status != store.GetStatusDescription(validatorInfo.Status) {
				log.Infow("UpdateValidatorStatus", "pubKey", validatorInfo.Validator.Pubkey, "status", validatorInfo.Status)
				err = bm.store.UpdateValidatorStatus(removePubKeyPrefix(validatorInfo.Validator.Pubkey), store.GetStatusDescription(validatorInfo.Status))
				if err != nil {
					log.Errorw("UpdateValidatorStatus", "err", err)
					return err
				}
			}

			if store.GetStatusDescription(validatorInfo.Status) == store.ValidatorExited {
				err = bm.store.UpdateValidatorOnlineStatus(int64(validatorInfo.Index), false)
				if err != nil {
					log.Errorw("UpdateValidatorOnlineStatus", "err", err, "validatorIndex", validatorInfo.Index, "online", false)
					return err
				}

				exitEpoch := epoch
				if validatorInfo.Validator.ExitEpoch != 0 {
					exitEpoch = uint64(validatorInfo.Validator.ExitEpoch)
				} else {
					log.Warnw("ExitValidator", "validatorIndex", validatorInfo.Index, "status", validatorInfo.Status, "validatorInfo.Validator.ExitEpoch", 0)
				}

				blockNumber, err := bm.getEth1ExBlock(exitEpoch * 32)
				if err != nil {
					log.Errorw("GetEth1ExBlock", "err", err)
					return err
				}
				err = bm.store.ExitValidator(removePubKeyPrefix(validatorInfo.Validator.Pubkey), int64(blockNumber))
				if err != nil {
					log.Errorw("ExitValidator", "err", err)
					return err
				}
			}

			if validatorInfo.Validator.Slashed {
				clusterSlashAlarms[v.ClusterID] = append(clusterSlashAlarms[v.ClusterID], uint64(validatorInfo.Index))
				log.Infow("ValidatorSlash", "pubKey", validatorInfo.Validator.Pubkey)

				err = bm.store.ValidatorSlash(removePubKeyPrefix(validatorInfo.Validator.Pubkey))
				if err != nil {
					log.Errorw("ValidatorSlash", "err", err)
					return err
				}
			}
		}

		if page*itemsPerPage >= int(totalCount) {
			break
		}
		page++
	}

	for clusterId, balanceAlarms := range clusterBalanceAlarms {
		if len(balanceAlarms) > 0 {
			bm.validatorBalanceDeltaAlarmChan <- alert.ValidatorBalanceDeltaNotify{
				Epoch:     epoch,
				ClusterId: clusterId,
				Index:     balanceAlarms,
			}
		}
	}

	for clusterId, slashAlarms := range clusterSlashAlarms {
		if len(slashAlarms) > 0 {
			bm.validatorSlashAlarmChan <- alert.ValidatorSlashNotify{
				Epoch:     epoch,
				ClusterId: clusterId,
				Index:     slashAlarms,
			}
		}
	}

	return nil
}

func mergeMaps[K comparable, V any](map1, map2 map[K]V) {
	for k, v := range map2 {
		map1[k] = v
	}
}

func (bm *BeaconMonitor) updateBalanceHistory(validatorIndex, epoch, balance uint64) {
	history := bm.validatorBalanceHistory[validatorIndex]
	newBalance := Balance{Epoch: epoch, Amount: balance}

	history[2] = newBalance

	sort.Slice(history[:], func(i, j int) bool {
		return history[i].Epoch > history[j].Epoch
	})

	bm.validatorBalanceHistory[validatorIndex] = history
}

func (bm *BeaconMonitor) checkBalanceChange(validatorIndex uint64, isOnline bool) bool {
	history := bm.validatorBalanceHistory[validatorIndex]
	if history[0].Epoch == 0 || history[1].Epoch == 0 || history[2].Epoch == 0 {
		return false
	}

	if history[0].Amount > history[1].Amount {
		if !isOnline {
			log.Infow("checkBalanceChange: UpdateValidatorOnlineStatus", "validatorIndex", validatorIndex, "online", true)
			err := bm.store.UpdateValidatorOnlineStatus(int64(validatorIndex), true)
			if err != nil {
				log.Errorw("UpdateValidatorOnlineStatus", "err", err, "validatorIndex", validatorIndex, "online", true)
			}
		}
	}

	if history[0].Amount != EffectiveBalance && history[0].Amount < history[1].Amount && history[1].Amount < history[2].Amount {
		if isOnline {
			log.Infow("checkBalanceChange: UpdateValidatorOnlineStatus", "validatorIndex", validatorIndex, "online", false)
			err := bm.store.UpdateValidatorOnlineStatus(int64(validatorIndex), false)
			if err != nil {
				log.Errorw("UpdateValidatorOnlineStatus", "err", err, "validatorIndex", validatorIndex, "online", false)
			}
		}

		log.Infow("Validator balance decrease", "validatorIndex", validatorIndex, "curBalance", history[0].Amount, "preBalance", history[1].Amount, "epoch", history[0].Epoch)
		return true
	}

	return false
}
