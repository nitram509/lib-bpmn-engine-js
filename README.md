# lib-bpmn-engine-js

## status

**experimental** playground is online
https://nitram509.github.io/lib-bpmn-engine-js/

## build

```shell
GOOS=js GOARCH=wasm go build -o static/main.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static/wasm_exec.js
```
