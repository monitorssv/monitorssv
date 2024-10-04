package store

import (
	"encoding/json"
	"github.com/monitorssv/monitorssv/config"
	"testing"
)

func initDB(t *testing.T) *Store {
	cfg, err := config.InitConfig("../deploy/monitorssv/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	db, err := NewStore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestStore_CalculateChartData(t *testing.T) {
	db := initDB(t)

	chart, err := db.CalculateChartData()
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(chart)
	t.Log(string(data))
}
