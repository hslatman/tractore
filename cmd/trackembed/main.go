package main

import (
	"encore.app/pkg/trackembed"
	"github.com/extism/go-pdk"
)

type Input struct {
	Raw         []byte `json:"raw"`
	TrackingURL string `json:"trackingURL"`
}

type Output struct {
	Raw []byte `json:"raw"`
}

//export embed
func embed() int32 {
	params := Input{}
	err := pdk.InputJSON(&params)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	raw, err := trackembed.Pixel(params.Raw, params.TrackingURL)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	if err = pdk.OutputJSON(Output{Raw: raw}); err != nil {
		pdk.SetError(err)
		return 1
	}

	return 0
}

// run embed, because Wasm modules compiled using the Go stdlib don't support
// calling exported methods directly. The plugin code is called using "_start".
func main() {
	embed()
}
