package main

import (
	"context"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"code-is-fun/client"
)

// proposal accepts slot only (since there is only 1 proposer per slot, there is no need to specify a cluster).
func proposal(targetSlot phase0.Slot) {
	const baseURL = "http://mainnet-ssv-node-exporter-1.production.vnet.ops.ssvlabsinternal.com"
	c := client.NewClient(baseURL)

	resp, err := c.GetValidatorTraces(context.Background(), uint64(targetSlot), uint64(targetSlot), []string{"PROPOSER"})
	if err != nil {
		panic(fmt.Sprintf("error getting validator traces: %v", err))
	}

	prefix := "  "

	for _, entry := range resp.Data {
		fmt.Println("Pre:")
		for _, pre := range entry.Pre {
			tis, err := timeIntoSlot(targetSlot, parseTime(pre.Time))
			if err != nil {
				panic(fmt.Sprintf("timeIntoSlot: %v", err))
			}

			fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
		}

		fmt.Println()

		fmt.Println("Consensus:")
		for i, consensus := range entry.Consensus {
			round := i + 1

			fmt.Println("----------" + fmt.Sprintf("[round=%d]", round) + "----------")

			fmt.Println(prefix + "proposal:")
			{
				prop := consensus.Proposal
				if prop != nil {
					prefix := prefix + "  "

					tis, err := timeIntoSlot(targetSlot, parseTime(prop.Time))
					if err != nil {
						panic(fmt.Sprintf("timeIntoSlot: %v", err))
					}

					fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
				}
			}

			fmt.Println(prefix + "prepare:")
			for _, prep := range consensus.Prepares {
				prefix := prefix + "  "

				tis, err := timeIntoSlot(targetSlot, parseTime(prep.Time))
				if err != nil {
					panic(fmt.Sprintf("timeIntoSlot: %v", err))
				}

				fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
			}

			fmt.Println(prefix + "commit:")
			for _, commit := range consensus.Commits {
				prefix := prefix + "  "

				tis, err := timeIntoSlot(targetSlot, parseTime(commit.Time))
				if err != nil {
					panic(fmt.Sprintf("timeIntoSlot: %v", err))
				}

				fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
			}

			fmt.Println(prefix + "round-change:")
			for _, rc := range consensus.RoundChanges {
				prefix := prefix + "  "

				tis, err := timeIntoSlot(targetSlot, parseTime(rc.Time))
				if err != nil {
					panic(fmt.Sprintf("timeIntoSlot: %v", err))
				}

				fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
			}
		}
		fmt.Println("----------------------------")

		fmt.Println()

		fmt.Println("Post:")
		for _, post := range entry.Post {
			tis, err := timeIntoSlot(targetSlot, parseTime(post.Time))
			if err != nil {
				panic(fmt.Sprintf("timeIntoSlot: %v", err))
			}

			fmt.Println(prefix + fmt.Sprintf("%d ms", tis.Milliseconds()))
		}
	}
}
