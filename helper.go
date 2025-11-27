package main

import (
	"fmt"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func timeIntoSlot(targetSlot phase0.Slot, t time.Time) (time.Duration, error) {
	targetSlotStartTime, err := BlockchainMainnet.SlotStartTime(targetSlot)
	if err != nil {
		return 0, fmt.Errorf("get target slot start time: %w", err)
	}
	return t.Sub(targetSlotStartTime), nil
}

func parseTime(ts string) time.Time {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		panic(fmt.Sprintf("time.Parse: %v", err))
	}

	return t
}
