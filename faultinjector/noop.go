//go:build !faultinjection

// When built WITHOUT the "faultinjection" tag, all fault injection functions
// are compiled as no-ops. The release binary contains none of the fault
// injection logic, HTTP handler, persistence, or os.Exit() paths.
package faultinjector

import "net/http"

// Check is a no-op in release builds — always returns nil.
func Check(string) error { return nil }

// SetPersistPath is a no-op in release builds.
func SetPersistPath(string) {}

// Enable is a no-op in release builds.
func Enable() {}

// NewHTTPServer returns nil in release builds — no HTTP server is created.
func NewHTTPServer(string) *http.Server { return nil }

// CodeOf is a no-op in release builds — always returns "". Mirrors the
// faultinjection-tagged CodeOf so production cutpoint code using the documented
// switch faultinjector.CodeOf(err) pattern compiles in both build modes.
func CodeOf(error) string { return "" }
