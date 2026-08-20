package txsim

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

var abiLine = regexp.MustCompile(`^\s*ABI:\s*(.*),$`)

var contractErrors map[string]abi.Error

func PopulateErrorMap(bindingsPath string) error {
	contractErrors = make(map[string]abi.Error)

	entries, err := os.ReadDir(bindingsPath)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("reading bindings directory %s: %w", bindingsPath, err))
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		// The name of the contract name must match the name of its directory.
		contractName := e.Name()
		dirPath := filepath.Join(bindingsPath, contractName)

		files, err := os.ReadDir(dirPath)
		if err != nil {
			return withstack.Wrap(fmt.Errorf("reading contract directory %s: %w", dirPath, err))
		}

		// Expect exactly one file in the subdir; pick the first regular file.
		if len(files) != 1 {
			return fmt.Errorf("expected exactly one file in directory, but found %d", len(files))
		}

		currentContractErrors, err := parseErrorsFromBinding(filepath.Join(dirPath, files[0].Name()))
		if err != nil {
			return fmt.Errorf("parsing errors from binding %s: %w", files[0].Name(), err)
		}

		for _, errorABI := range currentContractErrors {
			contractErrors[errorABI.ID.Hex()[2:10]] = errorABI
		}
	}

	return nil
}

func parseErrorsFromBinding(bindingPath string) (map[string]abi.Error, error) {
	f, err := os.Open(filepath.Clean(bindingPath))
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("opening binding file %s: %w", bindingPath, err))
	}
	defer func() { _ = f.Close() }()

	var abiString string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if m := abiLine.FindStringSubmatch(line); m != nil {
			abiString, _ = strconv.Unquote(m[1])
			break

		}
	}
	if scanErr := sc.Err(); scanErr != nil {
		return nil, withstack.Wrap(fmt.Errorf("scanning binding file %s: %w", bindingPath, scanErr))
	}
	if abiString == "" {
		return nil, fmt.Errorf("no ABI line found in %s", bindingPath)
	}

	parsed, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return parsed.Errors, nil
}
