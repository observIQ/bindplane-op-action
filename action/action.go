package action

import (
	"context"
	"fmt"

	"github.com/observiq/bindplane-op-action/client"
	"github.com/observiq/bindplane-op-action/client/config"
	"github.com/observiq/bindplane-op-action/client/model"
	"github.com/observiq/bindplane-op-action/client/version"

	"go.uber.org/zap"
)

type rType string

const (
	destination   rType = "destination"
	source        rType = "source"
	processor     rType = "processor"
	configuration rType = "configuration"
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

// Apply applies destinations, sources, processors, and configurations
// in that order. It is important to apply destinations first, followed
// by resource library sources and processors. Configurations should be
// applied last because they will reference other resources.
func (a *Action) Apply() error {
	if err := a.apply(destination, a.destinationPath); err != nil {
		return fmt.Errorf("destinations: %w", err)
	}

	if err := a.apply(source, a.sourcePath); err != nil {
		return fmt.Errorf("sources: %w", err)
	}

	if err := a.apply(processor, a.processorPath); err != nil {
		return fmt.Errorf("processors: %w", err)
	}

	if err := a.apply(configuration, a.configurationPath); err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	return nil
}

// apply takes a resource type and a file path and applies it to the BindPlane API
// If an error is found in the response status, it will be returned
func (a *Action) apply(resourceType rType, path string) error {
	a.Logger.Info("Applying resource", zap.String("type", string(resourceType)), zap.String("file", path))
	resp, err := a.client.ApplyFile(context.Background(), path)
	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	for _, s := range resp {
		name := s.Resource.Metadata.Name
		id := s.Resource.Metadata.ID
		kind := s.Resource.Kind
		status := s.Status
		a.Logger.Info(
			"Resource applied",
			zap.String("name", name),
			zap.String("id", id),
			zap.String("kind", kind),
			zap.String("status", string(status)),
		)

		switch status {
		case model.StatusUnchanged, model.StatusConfigured, model.StatusCreated:
			continue
		case model.StatusInvalid:
			return fmt.Errorf("invalid resource: %s: %s", name, s.Reason)
		case model.StatusError:
			return fmt.Errorf("error: %s: %s", name, s.Reason)
		case model.StatusForbidden:
			return fmt.Errorf("forbidden: %s: %s", name, s.Reason)
		default:
			return fmt.Errorf("unexpected status: %s", status)
		}
	}

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
		return fmt.Errorf("failed to apply resources: %w", err)
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
