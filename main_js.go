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

type EngineWrapper struct {
	engine *bpmn_engine.BpmnEngineState
}

type ActivatedJobHandler struct {
	handler js.Value
}

var engineCounter = 0

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("NewBpmnEngine", js.FuncOf(newBpmnEngine))
	<-done
}

func newBpmnEngine(this js.Value, args []js.Value) interface{} {
	engine := bpmn_engine.New(fmt.Sprintf("engine-%d", engineCounter))
	ew := EngineWrapper{
		engine: &engine,
	}
	engineCounter++
	obj := vert.ValueOf(jsBinding{Name: engine.GetName()})
	obj.Set("GetName", js.FuncOf(func(this js.Value, args []js.Value) any {
		return engine.GetName()
	}))
	obj.Set("CreateAndRunInstance", js.FuncOf(ew.CreateAndRunInstance))
	obj.Set("LoadFromString", js.FuncOf(ew.JsLoadFromString))
	obj.Set("NewTaskHandlerForId", js.FuncOf(ew.JsNewTaskHandlerForId))
	return obj.JSValue()
}

func (ew EngineWrapper) JsLoadFromString(this js.Value, args []js.Value) any {
	xmlString := args[0].String()
	process, err := ew.engine.LoadFromBytes([]byte(xmlString))
	if err != nil {
		return js.ValueOf(false)
	}
	return js.ValueOf(process.ProcessKey)
}

func (ew EngineWrapper) JsNewTaskHandlerForId(this js.Value, args []js.Value) any {
	id := args[0].String()
	jsHandler := args[1]
	ajh := ActivatedJobHandler{
		handler: jsHandler,
	}
	ew.engine.NewTaskHandler().Id(id).Handler(ajh.JsActivatedJobHandler)
	return js.Undefined()
}

func (ew EngineWrapper) CreateAndRunInstance(this js.Value, args []js.Value) any {
	processKey := int64(args[0].Float())
	ew.engine.CreateAndRunInstance(processKey, nil)
	return js.Undefined()
}

func (ajh ActivatedJobHandler) JsActivatedJobHandler(job bpmn_engine.ActivatedJob) {
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
		return vert.ValueOf(job.GetCreatedAt()).JSValue()
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
	ajh.handler.Invoke(obj.JSValue())
}
