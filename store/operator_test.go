package store

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStore_BatchUpdateOperatorValidatorCount(t *testing.T) {
	db := initDB(t)

	err := db.BatchUpdateOperatorValidatorCount([]uint64{1, 2, 3, 4}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStore_GetActiveOperatorCount(t *testing.T) {
	db := initDB(t)

	count, err := db.GetActiveOperatorCount()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestStore_GetOperatorByOwner(t *testing.T) {
	db := initDB(t)

	operators, totalCount, err := db.GetOperatorByOwner(1, 10, "0x24F34a87a28088cf58808D03C7f6017C6aD2150e")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(operators)
}

func TestStore_UpdateOperatorClusterIds(t *testing.T) {
	db := initDB(t)
	db.UpdateOperatorClusterIds(1, "e9a5b0fb62ec5a4a0b4614bf7a6313a65e3713d487657749c816123b07222fba", true)
	db.UpdateOperatorClusterIds(1, "4244d5b8c14b4e515366b691d13e457c00e034b26361390e229af1ee1999a59c", true)
	db.UpdateOperatorClusterIds(1, "87ab407cbdf02aba3f01e2a57667cf3524d219a5aac62b59782f25acae486f05", true)
	db.UpdateOperatorClusterIds(1, "e9a5b0fb62ec5a4a0b4614bf7a6313a65e3713d487657749c816123b07222fba", false)
}

func TestStore_UpdateAllOperatorClusterIds(t *testing.T) {
	db := initDB(t)
	clusters, err := db.GetAllClusters()
	if err != nil {
		t.Fatal(err)
	}

	for _, cluster := range clusters {
		if cluster.ValidatorCount == 0 {
			continue
		}
		operatorIds := strings.Split(cluster.OperatorIds, ",")
		for _, operatorIdStr := range operatorIds {
			operatorId, err := strconv.Atoi(operatorIdStr)
			if err != nil {
				t.Fatal(err)
			}

			err = db.UpdateOperatorClusterIds(uint64(operatorId), cluster.ClusterID, true)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestGetOperators(t *testing.T) {
	db := initDB(t)
	start := time.Now()
	operators, totalCount, err := db.GetOperators(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(operators)
	t.Log(time.Since(start))
}

func TestStore_getOperatorChanges(t *testing.T) {
	db := initDB(t)
	operatorChanges, err := db.getOperatorChanges()
	if err != nil {
		t.Fatal(err)
	}
	totalOperator := 0
	for _, operatorChange := range operatorChanges {
		totalOperator += operatorChange.Change
	}
	t.Log(totalOperator)
}
