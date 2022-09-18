# lib-bpmn-engine-js

## status

**experimental**

## build

```shell
GOOS=js GOARCH=wasm go build -o static/main.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static/wasm_exec.js
```
