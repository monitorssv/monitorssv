package client

import (
	"errors"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/config"
	"testing"
)

func initClient(t *testing.T) *Client {
	_ = logging.SetLogLevel("*", "INFO")
	cfg, err := config.InitConfig("../../deploy/monitorssv/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(cfg.Eth2Rpc)
	return client
}

func TestClient(t *testing.T) {
	client := initClient(t)

	proposers, err := client.GetEpochProposer(310973)
	if err != nil {
		t.Fatal(err)
	}
	for _, proposer := range proposers.Data {
		t.Log(proposer)
		header, err := client.GetSlotHeader(uint64(proposer.Slot))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(proposer.Slot, "----", header.Data.Header.Message.ProposerIndex)
	}
}

func TestGetHeader(t *testing.T) {
	client := initClient(t)
	_, err := client.GetSlotHeader(9848730)
	if !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
	_, err = client.GetBlockBySlot(9848730)
	if !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
}

func TestGetSlotValidatorsByPubKey(t *testing.T) {
	var pubkeys = []string{
		"0x981ff2b59a8d4ea1aeeed19f9a9e2d6895f86508833d3e27b92f081920d7015ae68877ab24728d0d30e5f63021a250fc",
	}
	client := initClient(t)
	validators, err := client.GetSlotValidatorsByPubKey(313332*32, pubkeys)
	if err != nil {
		t.Fatal(err)
	}
	for _, validator := range validators {
		t.Log(validator)
	}
}

func TestGetSlotValidatorsByIndex(t *testing.T) {
	var indexs = []uint64{
		1538889,
	}

	client := initClient(t)
	validators, err := client.GetSlotValidatorsByIndex(9977247, indexs)
	if err != nil {
		t.Fatal(err)
	}
	for _, validator := range validators {
		t.Log(validator)
	}
}

func TestFinalizedEpoch(t *testing.T) {
	client := initClient(t)
	slot, err := client.GetLatestSlot()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("---", slot)

	finalizedEpoch, err := client.GetFinalizedEpoch()
	if err != nil {
		t.Fatal(err)
	}
	finalizedSlot := finalizedEpoch * 32
	t.Log(finalizedEpoch, finalizedSlot)
}

func TestGetBlockBySlot(t *testing.T) {
	client := initClient(t)
	block, err := client.GetBlockBySlot(9206550)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block)
}
