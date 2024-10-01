package ssv

import (
	"testing"
)

func TestGetContractCreator(t *testing.T) {
	ssv := initSSV(t)
	creations, err := GetContractCreator("0x9e9b391344917d88D3eeb9144089dc8f72f42583", ssv.cfg.EtherScan.Endpoint, ssv.cfg.EtherScan.ApiKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(creations)
}

func TestClusterOwner(t *testing.T) {
	ssv := initSSV(t)
	clusters, err := ssv.store.GetAllClusters()
	if err != nil {
		t.Fatal(err)
	}

	for _, cluster := range clusters {
		code, err := ssv.client.CodeAt(cluster.Owner)
		if err != nil {
			ssvLog.Warnf("failed to CodeAt: %v", err)
			continue
		}

		if len(code) == 0 {
			ssvLog.Infow("eoa owner", "owner", cluster.Owner, "EoaOwner", cluster.EoaOwner)
			if cluster.Owner != cluster.EoaOwner {
				ssvLog.Warnw("owner not match")
			}
			continue
		}

		creations, err := GetContractCreator(cluster.Owner, ssv.cfg.EtherScan.Endpoint, ssv.cfg.EtherScan.ApiKey)
		if err != nil {
			t.Fatal(err)
		}

		if creations.ContractCreator != cluster.EoaOwner {
			ssvLog.Warnw("eoa owner not match")
		} else {
			ssvLog.Infow("contract owner", "owner", cluster.Owner, "EoaOwner", cluster.EoaOwner)
		}
	}
}
