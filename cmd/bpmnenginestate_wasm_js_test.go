package main

import (
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"testing"
)

func Test(t *testing.T) {
	wrap(bpmn_engine.New("test"))
}
