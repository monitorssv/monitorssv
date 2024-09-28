package ssv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiOperatorData struct {
	ID                   int    `json:"id"`
	IDStr                string `json:"id_str"`
	DeclaredFee          string `json:"declared_fee"`
	PreviousFee          string `json:"previous_fee"`
	Fee                  string `json:"fee"`
	PublicKey            string `json:"public_key"`
	OwnerAddress         string `json:"owner_address"`
	AddressWhitelist     string `json:"address_whitelist"`
	IsPrivate            bool   `json:"is_private"`
	WhitelistingContract string `json:"whitelisting_contract"`
	Location             string `json:"location"`
	SetupProvider        string `json:"setup_provider"`
	Eth1NodeClient       string `json:"eth1_node_client"`
	Eth2NodeClient       string `json:"eth2_node_client"`
	MevRelays            string `json:"mev_relays"`
	Description          string `json:"description"`
	WebsiteURL           string `json:"website_url"`
	TwitterURL           string `json:"twitter_url"`
	LinkedinURL          string `json:"linkedin_url"`
	DkgAddress           string `json:"dkg_address"`
	Logo                 string `json:"logo"`
	Type                 string `json:"type"`
	Name                 string `json:"name"`
	Performance          struct {
		Day   float64 `json:"24h"`
		Month float64 `json:"30d"`
	} `json:"performance"`
	IsValid         bool   `json:"is_valid"`
	IsDeleted       bool   `json:"is_deleted"`
	IsActive        int    `json:"is_active"`
	Status          string `json:"status"`
	ValidatorsCount int    `json:"validators_count"`
	Version         string `json:"version"`
	Network         string `json:"network"`
}

var ssvOperatorApi = "https://api.ssv.network/api/v4/%s/operators/%d"

func GetOperatorName(network string, operatorId uint64) (string, error) {
	url := fmt.Sprintf(ssvOperatorApi, network, operatorId)
	b, err := httpGet(url)
	if err != nil {
		return "", err
	}

	var operator ApiOperatorData
	if err = json.Unmarshal(b, &operator); err != nil {
		return "", err
	}

	return operator.Name, nil
}

func httpGet(url string) ([]byte, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func toStrOperatorIds(operatorIds []uint64) string {
	var operatorIdsStr string
	for _, operatorId := range operatorIds {
		if operatorIdsStr == "" {
			operatorIdsStr = fmt.Sprintf("%d", operatorId)
			continue
		}

		operatorIdsStr = fmt.Sprintf("%s,%d", operatorIdsStr, operatorId)
	}

	return operatorIdsStr
}
