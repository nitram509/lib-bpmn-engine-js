package main

import (
	_ "embed"
	"fmt"

	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"syscall/js"
)

//go:embed "simple_task.bpmn"
var simpleTaskBpmn []byte

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("runEngine", js.FuncOf(runEngine))
	<-done
}

func runEngine(this js.Value, args []js.Value) interface{} {
	// create a new named engine
	bpmnEngine := bpmn_engine.New("a name")
	// basic example loading a BPMN from file,
	process, err := bpmnEngine.LoadFromBytes(simpleTaskBpmn)
	if err != nil {
		panic("file \"simple_task.bpmn\" can't be read.")
	}
	// register a handler for a service task by defined task type
	bpmnEngine.NewTaskHandler().Id("hello-world").Handler(printContextHandler)
	// setup some variables
	variables := map[string]interface{}{}
	variables["foo"] = "bar"
	// and execute the process
	bpmnEngine.CreateAndRunInstance(process.ProcessKey, variables)
	return true
}

func printContextHandler(job bpmn_engine.ActivatedJob) {
	println("< Hello World >")
	println(fmt.Sprintf("ElementId                = %s", job.GetElementId()))
	println(fmt.Sprintf("BpmnProcessId            = %s", job.GetBpmnProcessId()))
	println(fmt.Sprintf("ProcessDefinitionKey     = %d", job.GetProcessDefinitionKey()))
	println(fmt.Sprintf("ProcessDefinitionVersion = %d", job.GetProcessDefinitionVersion()))
	println(fmt.Sprintf("CreatedAt                = %s", job.GetCreatedAt()))
	println(fmt.Sprintf("Variable 'foo'           = %s", job.GetVariable("foo")))
	job.Complete() // don't forget this one, or job.Fail("foobar")
}
