package store

import (
	"testing"
)

func TestDeleteAlarmByEoaOwner(t *testing.T) {
	db := initDB(t)

	alarmInfo, err := db.GetAlarmByEoaOwner("0x52EC98881E3a62452E8f6bFb74290B51a442975b")
	if err != nil {
		t.Fatal(err)
	}

	if alarmInfo == nil {
		t.Fatal("test data not found")
	}

	// delete test data
	err = db.DeleteAlarmByEoaOwner("0x52EC98881E3a62452E8f6bFb74290B51a442975b")
	if err != nil {
		t.Fatal(err)
	}

	// re create test data
	err = db.CreateOrUpdateAlarmInfo(&AlarmInfo{
		EoaOwner:                   alarmInfo.EoaOwner,
		AlarmType:                  alarmInfo.AlarmType,
		AlarmChannel:               alarmInfo.AlarmChannel,
		AlarmChannelHash:           alarmInfo.AlarmChannelHash,
		ReportLiquidationThreshold: alarmInfo.ReportLiquidationThreshold,
		ReportOperatorFeeChange:    alarmInfo.ReportOperatorFeeChange,
		ReportNetworkFeeChange:     alarmInfo.ReportNetworkFeeChange,
		ReportProposeBlock:         alarmInfo.ReportProposeBlock,
		ReportMissedBlock:          alarmInfo.ReportMissedBlock,
		ReportBalanceDecrease:      alarmInfo.ReportBalanceDecrease,
		ReportExitedButNotRemoved:  alarmInfo.ReportExitedButNotRemoved,
		ReportWeekly:               alarmInfo.ReportWeekly,
	})
	if err != nil {
		t.Fatal(err)
	}
}
