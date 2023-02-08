package test_gui

import (
	"app/src/state"
	"app/src/utils"

	"github.com/robfig/cron/v3"
)

/*
	Test GUI package generates random data to test various GUI elements
*/

func MakeFakeSyncingData() {
	limit := 4

	// on call
	for i := 0; i < limit; i++ {
		state.SyncingData = append(state.SyncingData, state.SyncingDataStruct{Label: utils.RandStringRunes(20), Progress: utils.RandFloat64(0.0, 1.0), Status: utils.RandStringRunes(10)})
	}

	// on cron
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every 1s", func() {
		state.SyncingData = nil
		for i := 0; i < limit; i++ {
			state.SyncingData = append(state.SyncingData, state.SyncingDataStruct{Label: utils.RandStringRunes(20), Progress: utils.RandFloat64(0.0, 1.0), Status: utils.RandStringRunes(10)})
		}
	})
	go c.Run()
}

func MakeFakeHistoryData() {
	limit := 30

	// on call
	for i := 0; i < limit; i++ {
		state.HistoryData = append(state.HistoryData, state.HistoryDataStruct{Label: utils.RandStringRunes(20), Status: utils.RandStringRunes(10)})
	}

	// on cron
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every 1s", func() {
		state.HistoryData = nil
		for i := 0; i < limit; i++ {
			state.HistoryData = append(state.HistoryData, state.HistoryDataStruct{Label: utils.RandStringRunes(20), Status: utils.RandStringRunes(10)})
		}
	})
	go c.Run()
}
