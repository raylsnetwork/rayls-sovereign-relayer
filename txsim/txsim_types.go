package txsim

import (
	"fmt"
	"strings"
)

type ContractError struct {
	Sig  string
	Args []interface{}
}

func (e ContractError) IsEmpty() bool {
	return e.Sig == "" && len(e.Args) == 0
}

// Name returns just the error identifier without its argument list, e.g.
// "ProgramData__Reverted(uint256,bytes32,bytes4,bytes)" → "ProgramData__Reverted". Useful for a
// concise `revert reason` log field.
func (e ContractError) Name() string {
	if i := strings.IndexByte(e.Sig, '('); i >= 0 {
		return e.Sig[:i]
	}
	return e.Sig
}

// Reason returns the most human-readable form of the revert for a log field.
// For a Solidity `require(false, "msg")` the decoder yields Sig "Error(string)"
// with the message in Args[0] — Name() would flatten that to "Error", so return
// the string arg instead. For custom errors (e.g. ProgramData__Reverted(...)),
// there is no embedded message, so fall back to the error identifier.
func (e ContractError) Reason() string {
	if e.Sig == "Error(string)" && len(e.Args) > 0 {
		if s, ok := e.Args[0].(string); ok {
			return s
		}
	}
	return e.Name()
}

func (e ContractError) String() string {
	if e.IsEmpty() {
		return ""
	}
	return fmt.Sprintf("%s: %v", e.Sig, e.Args)
}
