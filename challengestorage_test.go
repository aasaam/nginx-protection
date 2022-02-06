package main

import (
	"testing"
	"time"
)

func TestChallengeStorage(t *testing.T) {
	cs := newChallengeStorage()
	cs.gc()
	if cs.exist("n") {
		t.Errorf("must not exist")
	}
	cs.set("n", 0)
	cs.set("m", 5)
	if !cs.exist("n") {
		t.Errorf("exist")
	}
	time.Sleep(time.Second * 1)
	cs.gc()
	if cs.exist("n") {
		t.Errorf("must gc")
	}
}
