package main

import (
	_ "embed"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/norunners/vert"
	"syscall/js"
)

//go:embed "simple_task.bpmn"
var simpleTaskBpmn []byte

type jsBinding struct {
	Name string `js:"name"`
}

var engineCounter = 0

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("runEngine", js.FuncOf(runEngine))
	js.Global().Set("NewBpmnEngine", js.FuncOf(newBpmnEngine))
	<-done
}

func newBpmnEngine(this js.Value, args []js.Value) interface{} {
	engine := bpmn_engine.New(fmt.Sprintf("engine-%d", engineCounter))
	process, _ := engine.LoadFromBytes(simpleTaskBpmn)
	engineCounter++
	obj := vert.ValueOf(jsBinding{Name: engine.GetName()})
	obj.Set("GetName", js.FuncOf(func(this js.Value, args []js.Value) any {
		return engine.GetName()
	}))
	obj.Set("CreateAndRunInstance", js.FuncOf(func(this js.Value, args []js.Value) any {
		engine.CreateAndRunInstance(process.ProcessKey, nil)
		return js.Undefined()
	}))
	obj.Set("NewTaskHandlerForId", js.FuncOf(func(this js.Value, args []js.Value) any {
		id := args[0].String()
		jsHandler := args[1]
		engine.NewTaskHandler().Id(id).Handler(func(job bpmn_engine.ActivatedJob) {
			type retObj struct{}
			obj := vert.ValueOf(retObj{})
			obj.Set("GetKey", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetKey()
			}))
			obj.Set("GetProcessInstanceKey", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetProcessInstanceKey()
			}))
			obj.Set("GetBpmnProcessId", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetBpmnProcessId()
			}))
			obj.Set("GetProcessDefinitionVersion", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetProcessDefinitionVersion()
			}))
			obj.Set("GetProcessDefinitionKey", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetProcessDefinitionKey()
			}))
			obj.Set("GetElementId", js.FuncOf(func(this js.Value, args []js.Value) any {
				return job.GetElementId()
			}))
			obj.Set("GetCreatedAt", js.FuncOf(func(this js.Value, args []js.Value) any {
				return ValueOf(job.GetCreatedAt()).JSValue()
			}))
			obj.Set("Fail", js.FuncOf(func(this js.Value, args []js.Value) any {
				reason := args[0].String()
				job.Fail(reason)
				return js.Undefined()
			}))
			obj.Set("Complete", js.FuncOf(func(this js.Value, args []js.Value) any {
				job.Complete()
				return js.Undefined()
			}))
			jsHandler.Invoke(obj.JSValue())
		})
		return js.Undefined()
	}))
	return obj.JSValue()
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
