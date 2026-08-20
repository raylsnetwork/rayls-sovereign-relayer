// Package testdata provides shared fixture builders for the DvP v2 unit-test
// suite. Tests in dvp/handler, dvp/service, contractclient,
// key-operation-service-app/service, and private-relayer/dest/... import
// from here instead of defining their own ad-hoc helpers.
//
// The package mirrors the functional-options idiom used by
// private-relayer/source/logparser/testdata: every constructor returns a
// sane default value and accepts a variadic of Option closures that mutate
// individual fields. This avoids the combinatorial explosion of "one
// helper per scenario" and keeps load-bearing invariants (like the v2
// salt-mirror rule on DvpSwap — see dvp/handler/receiver.go:441-475) in a
// single place.
//
// Sibling *_test.go files in this package validate that defaults are
// well-formed and that each Option mutates the expected field — they are
// the regression net for the fixture builders themselves.
package testdata
