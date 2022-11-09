# lib-bpmn-engine-js

## status

**experimental** playground is online
https://nitram509.github.io/lib-bpmn-engine-js/

## link

The actual BPMN engine's sources are available here: https://github.com/nitram509/lib-bpmn-engine

## build

```shell
GOOS=js GOARCH=wasm go build -o static/main.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static/wasm_exec.js
```

## test

```shell
export PATH="$PATH:$(go env GOROOT)/misc/wasm"
GOOS=js GOARCH=wasm go test ./...
```

## WiP: test code generator

```shell
cd cmd
go run . -package=github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine -type=BpmnEngineState
```
