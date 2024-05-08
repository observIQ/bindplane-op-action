package action

import (
	"context"
	"fmt"

	"github.com/observiq/bindplane-op-action/client"
	"github.com/observiq/bindplane-op-action/client/config"
	"github.com/observiq/bindplane-op-action/client/version"

	"go.uber.org/zap"
)

// Option is a function that configures an Action option
type Option func(*Action)

// WithBindPlaneRemoteURL sets the remote URL for the BindPlane client
func WithBindPlaneRemoteURL(u string) Option {
	return func(a *Action) {
		a.config.Network.RemoteURL = u
	}
}

// WithBindPlaneAPIKey sets the API key for the BindPlane client
func WithBindPlaneAPIKey(k string) Option {
	return func(a *Action) {
		a.config.Auth.APIKey = k
	}
}

// WithBindPlaneUsername sets the username for the BindPlane client
func WithBindPlaneUsername(u string) Option {
	return func(a *Action) {
		a.config.Auth.Username = u
	}
}

// WithBindPlanePassword sets the password for the BindPlane client
func WithBindPlanePassword(p string) Option {
	return func(a *Action) {
		a.config.Auth.Password = p
	}
}

// WithTLSCACert sets the certificate authority for the BindPlane client
func WithTLSCACert(c string) Option {
	return func(a *Action) {
		if c == "" {
			return
		}
		a.config.Network.CertificateAuthority = []string{c}
	}
}

// WithTargetBranch sets the branch to read resources from
func WithTargetBranch(b string) Option {
	return func(a *Action) {
		a.targetBranch = b
	}
}

// WithDestinationPath sets the path to write resources to
func WithDestinationPath(p string) Option {
	return func(a *Action) {
		a.destinationPath = p
	}
}

// WithSourcePath sets the path to read resources from
func WithSourcePath(p string) Option {
	return func(a *Action) {
		a.sourcePath = p
	}
}

// WithProcessorPath sets the path to read processors from
func WithProcessorPath(p string) Option {
	return func(a *Action) {
		a.processorPath = p
	}
}

// WithConfigurationPath sets the path to read configuration from
func WithConfigurationPath(p string) Option {
	return func(a *Action) {
		a.configurationPath = p
	}
}

// WithOTELConfigWriteBack sets the flag to enable writing back configuration
func WithOTELConfigWriteBack(b bool) Option {
	return func(a *Action) {
		a.enableWriteBack = b
	}
}

// WithConfigurationOutputDir sets the directory to write back configuration to
func WithConfigurationOutputDir(d string) Option {
	return func(a *Action) {
		a.configurationOutputDir = d
	}
}

// WithGithubToken sets the token to authenticate with GitHub
func WithGithubToken(t string) Option {
	return func(a *Action) {
		a.githubToken = t
	}
}

// WithGithubURL sets the URL to the GitHub repository
func WithGithubURL(u string) Option {
	return func(a *Action) {
		a.githubURL = u
	}
}

// WithAutoRollout sets the flag to enable auto rollout
func WithAutoRollout(b bool) Option {
	return func(a *Action) {
		a.autoRollout = b
	}
}

// WithConfigurationOutputBranch sets the branch to write back the configuration to
func WithConfigurationOutputBranch(b string) Option {
	return func(a *Action) {
		a.configurationOutputBranch = b
	}
}

// New creates a new Action with a configured bindPlane client
func New(logger *zap.Logger, opts ...Option) (*Action, error) {
	action := &Action{}
	for _, opt := range opts {
		opt(action)
	}

	c, err := client.NewBindPlane(&action.config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create BindPlane client: %w", err)
	}
	action.client = c

	action.Logger = logger

	return action, nil
}

// Action is a struct that contains the BindPlane client
// and user defined configuration options
type Action struct {
	Logger *zap.Logger

	// Branch name and paths to read
	// resources from
	targetBranch      string
	destinationPath   string
	sourcePath        string
	processorPath     string
	configurationPath string

	// Auto rollout options
	autoRollout bool

	// Write back options
	enableWriteBack           bool
	configurationOutputDir    string
	configurationOutputBranch string
	githubToken               string
	githubURL                 string

	// Config holds the following options:
	// - Remote URL
	// - API Key
	// - Username
	// - Password
	// - Certificate Authority
	config config.Config

	client *client.BindPlane
}

// TestConnection wraps the BindPlane client's Version method
func (a *Action) TestConnection() (version.Version, error) {
	v, err := a.client.Version(context.Background())
	if err != nil {
		return version.Version{}, fmt.Errorf("failed to test connection: %w", err)
	}
	return v, err
}

func (a *Action) Apply() error {
	// TODO: Apply destinations, sources, processors, and configuration
	return nil
}

func (a *Action) AutoRollout() error {
	return nil
}

func (a *Action) WriteBack() error {
	return nil
}

// Run executes the action
func (a *Action) Run() error {
	if err := a.Apply(); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	if a.autoRollout {
		if err := a.AutoRollout(); err != nil {
			return fmt.Errorf("failed to auto rollout configuration: %s", err)
		}
	}

	if a.enableWriteBack {
		if err := a.WriteBack(); err != nil {
			return fmt.Errorf("failed to write back configuration: %s", err)
		}
	}

	return nil
}
