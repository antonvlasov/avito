package handler

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run(15)
	timer := time.NewTimer(120 * time.Second)
	<-timer.C
	Stop()
}
