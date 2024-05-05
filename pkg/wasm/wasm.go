package wasm

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	extism "github.com/extism/go-sdk"
)

//go:embed modules/embed.wasm
var embedPixelWasm []byte

type embeddedModule struct {
	data []byte
}

func (m embeddedModule) ToWasmData(ctx context.Context) (extism.WasmData, error) {
	return extism.WasmData{
		Data: m.data,
	}, nil
}

var _ extism.Wasm = (embeddedModule)(embeddedModule{})

// Plugin wraps an [extism.Plugin], and provides a facade for the
// Extism plugin function call conventions.
type Plugin struct {
	plugin *extism.Plugin
}

// New returns a new [Plugin], wrapping an [extism.Plugin]. The plugin
// code is loaded from an embedded Wasm module provided from the modules
// subdirectory, and embedded at compilation time.
func New() (*Plugin, error) {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			embeddedModule{
				data: embedPixelWasm,
			},
		},
	}
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(context.Background(), manifest, config, []extism.HostFunction{})
	if err != nil {
		return nil, fmt.Errorf("failed creating plugin: %w", err)
	}

	return &Plugin{
		plugin: plugin,
	}, nil
}

// EmbedPixel calls the "_start" function of the "trackembed" Wasm module
// provided through the embedded Wasm file. It expects the raw email contents,
// and will return the new raw contents if embedding the tracking pixel through
// the Wasm module succeeds.
func (p *Plugin) EmbedPixel(raw []byte, trackingURL string) ([]byte, error) {
	inputData, err := json.Marshal(struct {
		Raw         []byte `json:"raw"`
		TrackingURL string `json:"trackingURL"`
	}{
		Raw:         raw,
		TrackingURL: trackingURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed marshaling JSON: %w", err)
	}

	exit, out, err := p.plugin.Call("_start", inputData)
	if err != nil {
		return nil, fmt.Errorf("failed Wasm embedding pixel: %w", err)
	}

	if exit != 0 {
		return nil, fmt.Errorf("executing plugin failed: %w", err)
	}

	v := struct {
		Raw []byte `json:"raw"`
	}{
		Raw: raw,
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return nil, fmt.Errorf("failed unmarshaling plugin response: %w", err)
	}

	return v.Raw, nil
}
