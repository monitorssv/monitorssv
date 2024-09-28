package ssv

import (
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"testing"
)

func TestOperatorMultipleWhitelistRemoved(t *testing.T) {
	operatorInfoWhitelistedAddress := "0x6A6C79d8dA4d3B1C8963073529CD026b36817eB6,0x5E33db0b37622F7E6b2f0654aA7B985D854EA9Cb,0x6B7468504757a4918a96078F352572d172115263,0x87393BE8ac323F2E63520A6184e5A8A9CC9fC051"
	whitelistAddresses := []common.Address{common.HexToAddress("0x5E33db0b37622F7E6b2f0654aA7B985D854EA9Cb"), common.HexToAddress("0x87393BE8ac323F2E63520A6184e5A8A9CC9fC051")}
	for _, whitelistAddr := range whitelistAddresses {
		if !strings.Contains(operatorInfoWhitelistedAddress, whitelistAddr.String()) {
			ssvLog.Warnw("Event: OperatorMultipleWhitelistRemoved", "whitelistedAddress", operatorInfoWhitelistedAddress, "whitelistAddr", whitelistAddr)
		}
		operatorInfoWhitelistedAddress = strings.Replace(operatorInfoWhitelistedAddress, whitelistAddr.String(), "", -1)
		operatorInfoWhitelistedAddress = strings.Replace(operatorInfoWhitelistedAddress, ",,", ",", -1)
		operatorInfoWhitelistedAddress = strings.Trim(operatorInfoWhitelistedAddress, ",")
	}
	t.Log("------", operatorInfoWhitelistedAddress)
}
