// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js,wasm

package asset

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall/js"
)

func openAsset(name string) (File, error) {
	loc := js.Global().Get("location")

	if loc.Truthy() {
		// we're probably in a browser.
		resp, err := http.Get(name)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("asset: get %q: %s", name, resp.Status)
		}

		// we need a ReadSeeker, so we need to buffer the data
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &jsFile{bytes.NewReader(b)}, nil
	}

	// for NodeJS, use the same logic as desktop.
	if !filepath.IsAbs(name) {
		name = filepath.Join("assets", name)
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type jsFile struct {
	*bytes.Reader
}

func (*jsFile) Close() error { return nil }
