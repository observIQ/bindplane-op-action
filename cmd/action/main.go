package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/observiq/bindplane-op-action/action"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// argCount is the number of arguments passed to the action, and does not
// include the binary name itself (which is returned by os.Args[0]).
// When adding new arguments to the action, this number should be updated
// and new global variables should be declared and handled in parseArgs().
const argCount = 16

// Global variables will be used when creating the action configuration. These
// are the options set by the user. Their order in parseArgs() is important.
var (
	bindplane_remote_url          string
	bindplane_api_key             string
	bindplane_username            string
	bindplane_password            string
	target_branch                 string
	destination_path              string
	configuration_path            string
	enable_otel_config_write_back bool
	configuration_output_dir      string
	token                         string
	enable_auto_rollout           bool
	configuration_output_branch   string
	tls_ca_cert                   string
	source_path                   string
	processor_path                string
	github_url                    string
)

const (
	exitParseArgsError            = 100
	exitValidationError           = 101
	exitClientInitError           = 102
	exitClientTestConnectionError = 103
	exitLoggerInitError           = 104
	exitClientError               = 1
)

func main() {
	if err := parseArgs(); err != nil {
		fmt.Printf("Error parsing arguments: %s\n", err)
		os.Exit(exitParseArgsError)
	}

	if err := validate(); err != nil {
		fmt.Printf("Error validating arguments: %s\n", err)
		os.Exit(exitValidationError)
	}

	zapConf := zap.NewProductionConfig()
	zapConf.Level.SetLevel(zap.DebugLevel) // TODO(jsirianni): Expose this as an option
	zapConf.OutputPaths = []string{"stdout"}
	zapConf.DisableStacktrace = true
	zapConf.DisableCaller = true
	zapConf.EncoderConfig.TimeKey = "time"
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapConf.Build()
	if err != nil {
		fmt.Printf("failed to create logger: %s\n", err)
		os.Exit(exitLoggerInitError)
	}

	branch := strings.Split(os.Getenv("GITHUB_REF"), "/")[2]
	if branch != target_branch {
		logger.Info(
			"Skipping action, branch does not match target branch",
			zap.String("branch", branch),
			zap.String("target_branch", target_branch),
		)
		os.Exit(0)
	}

	action, err := action.New(
		logger,

		// Client options
		action.WithBindPlaneRemoteURL(bindplane_remote_url),
		action.WithBindPlaneAPIKey(bindplane_api_key),
		action.WithBindPlaneUsername(bindplane_username),
		action.WithBindPlanePassword(bindplane_password),
		action.WithTLSCACert(tls_ca_cert),

		// Base action options for reading resources
		// from the repo, to apply to bindplane
		action.WithDestinationPath(destination_path),
		action.WithSourcePath(source_path),
		action.WithProcessorPath(processor_path),
		action.WithConfigurationPath(configuration_path),

		// Auto rollout option(s)
		action.WithAutoRollout(enable_auto_rollout),

		// Write back option(s)
		action.WithOTELConfigWriteBack(enable_otel_config_write_back),
		action.WithConfigurationOutputDir(configuration_output_dir),
		action.WithConfigurationOutputBranch(configuration_output_branch),
		action.WithGithubToken(token),
		action.WithGithubURL(github_url),
	)
	if err != nil {
		fmt.Printf("Error creating action: %s\n", err)
		os.Exit(exitClientInitError)
	}

	logger.Info("Testing connection to BindPlane API")
	version, err := action.TestConnection()
	if err != nil {
		fmt.Printf("Error testing connection: %s\n", err)
		os.Exit(exitClientTestConnectionError)
	}
	logger.Info(
		"Connection to BindPlane API successful",
		zap.Any("bindplane_version", version.Tag),
	)

	if err := action.Run(); err != nil {
		action.Logger.Error("error running action", zap.Error(err))
		os.Exit(exitClientError)
	}

	os.Exit(0)
}
