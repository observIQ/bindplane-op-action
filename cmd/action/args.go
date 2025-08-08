package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/observiq/bindplane-op-action/action"
)

// parseArgs parses the arguments passed to the action. The action will always
// pass the correct number of arguments regardless of the user input. If the
// number of args passed is not expected, an error will be returned indicating
// a bug with the action.
func parseArgs() error {
	args := os.Args

	// Add one to account for arg 0 being the binary name
	count := argCount + 1
	if len(args) != count {
		return fmt.Errorf("Not enough arguments, expected 18, got %d. %s.", len(args), action.BugError)
	}

	// First arg is always the binary name, so we skip it. We could
	// also use args[1:] to get all args after the binary name but
	// that could introduce confusion as the first arg would be at index 0.
	_ = args[0]

	// Order matters. The action will pass in the exact order defined
	// in action.yml in the root of the action repository.
	bindplane_remote_url = args[1]
	bindplane_api_key = args[2]
	bindplane_username = args[3]
	bindplane_password = args[4]
	target_branch = args[5]
	destination_path = args[6]
	configuration_path = args[7]

	b, err := strconv.ParseBool(args[8])
	if err != nil {
		return fmt.Errorf("enable_otel_config_write_back must be a boolean value")
	}
	enable_otel_config_write_back = b

	configuration_output_dir = args[9]
	token = args[10]

	b, err = strconv.ParseBool(args[11])
	if err != nil {
		return fmt.Errorf("enable_auto_rollout must be a boolean value")
	}
	enable_auto_rollout = b

	configuration_output_branch = args[12]
	if configuration_output_branch == "" {
		configuration_output_branch = target_branch
	}

	tls_ca_cert = args[13]
	if err := writeTLSFile("ca.crt", tls_ca_cert); err != nil {
		return fmt.Errorf("failed to write TLS file: %w", err)
	}

	source_path = args[14]
	processor_path = args[15]
	github_url = args[16]
	user_agent = args[17]

	return nil
}

// writeTLSFile takes a file path and writes the given contents to it
func writeTLSFile(path string, contents string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600) // #nosec G304 user defined filepath
	if err != nil {
		return fmt.Errorf("failed to create ca.crt file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		return fmt.Errorf("failed to write to ca.crt file: %w", err)
	}

	return nil
}
