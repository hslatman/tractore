# trackembed

An Extism plugin 

## Compilation

```console
GOOS=wasip1 GOARCH=wasm go build -o ./pkg/wasm/modules/embed.wasm cmd/trackembed/main.go
```

Compilation using TinyGo will fail, because some of the email parsing logic depends 
on the `net/smtp` package.

## Usage

See the [Extism docs}(https://extism.org/docs/quickstart/host-quickstart) on how to run the plugin within an Extism runtime.