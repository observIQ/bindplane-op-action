package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/observiq/bindplane-op-action/action"
	"github.com/observiq/bindplane-op-action/internal/repo"
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

	// Retrieve the commit message from the head commit on the branch
	message, err := commitMessage(github_url, branch, token)
	if err != nil {
		logger.Error("error getting commit message", zap.Error(err))
		os.Exit(exitClientError)
	}

	// If the commit message contains `progress rollout <name>`, progress the rollout
	// for the configuration instead of running the full workflow.
	if name, ok := extractConfigName(message); ok {
		err := action.RunRollout(name)
		if err != nil {
			logger.Error("error progressing rollout", zap.Error(err))
			os.Exit(exitClientError)
		}
		return
	}

	// Run the full workflow
	if err := action.Run(); err != nil {
		action.Logger.Error("error running action", zap.Error(err))
		os.Exit(exitClientError)
	}

	os.Exit(0)
}

// commitMessage clones the repository and returns the commit message of the
// head commit on the provided branch.
func commitMessage(cloneURL, branch, token string) (string, error) {
	repo, err := repo.CloneRepo(cloneURL, branch, token)
	if err != nil {
		return "", fmt.Errorf("clone repository branch %s: %w", branch, err)
	}

	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("get head commit: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", fmt.Errorf("get commit object: %w", err)
	}

	return commit.Message, nil
}

// extractName extracts a configuration name from a commit message.
// The commit message should contain the suffix "progress rollout <name>"
//
// Examples:
// - progress rollout test
// - this is a commit message progress rollout test
func extractConfigName(input string) (string, bool) {
	input = strings.TrimSpace(input)
	pattern := `progress rollout (\S+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}
