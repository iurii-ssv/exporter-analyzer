package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_testSelectFromClosedChannel(t *testing.T) {
	fmt.Println("test started")
	go func() {
		for {
			fmt.Println("ping")
			time.Sleep(1 * time.Second)
		}
	}()
	fmt.Println("test ended")
}
