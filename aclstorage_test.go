package main

import (
	"testing"
	"time"
)

func TestACLStorage(t *testing.T) {
	cs := newACLStorage()
	cs.gc()
	if cs.exist("n") != nil {
		t.Errorf("must not exist")
	}
	cs.add("n", "acl", "name", 1)
	cs.add("m", "acl", "name", 5)
	if cs.exist("n") == nil {
		t.Errorf("exist")
	}
	time.Sleep(time.Second * 2)
	cs.gc()
	if cs.exist("n") != nil {
		t.Errorf("must gc")
	}
}
