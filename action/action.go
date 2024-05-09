package action

import (
	"context"
	"fmt"
	"os"

	"github.com/observiq/bindplane-op-action/client"
	"github.com/observiq/bindplane-op-action/client/config"
	"github.com/observiq/bindplane-op-action/client/model"
	"github.com/observiq/bindplane-op-action/client/version"
	"gopkg.in/yaml.v3"

	"go.uber.org/zap"
)

const (
	BugError = "This is a bug with the action, please reach out to support or file an issue on Github https://github.com/observIQ/bindplane-op-action/issues"
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
	if a.destinationPath != "" {
		a.Logger.Info("Applying resources", zap.String("type", string(destination)), zap.String("file", a.destinationPath))
		err := a.apply(a.destinationPath)
		if err != nil {
			return fmt.Errorf("destinations: %w", err)
		}
	} else {
		a.Logger.Info("No destination path provided, skipping destinations")
	}

	if a.sourcePath != "" {
		a.Logger.Info("Applying resources", zap.String("type", string(source)), zap.String("file", a.destinationPath))
		err := a.apply(a.sourcePath)
		if err != nil {
			return fmt.Errorf("sources: %w", err)
		}
	} else {
		a.Logger.Info("No source path provided, skipping sources")
	}

	if a.processorPath != "" {
		a.Logger.Info("Applying resources", zap.String("type", string(processor)), zap.String("file", a.destinationPath))
		err := a.apply(a.processorPath)
		if err != nil {
			return fmt.Errorf("processors: %w", err)
		}
	} else {
		a.Logger.Info("No processor path provided, skipping processors")
	}

	if a.configurationPath != "" {
		a.Logger.Info("Applying resources", zap.String("type", string(configuration)), zap.String("file", a.destinationPath))
		err := a.apply(a.configurationPath)
		if err != nil {
			return fmt.Errorf("configuration: %w", err)
		}
	} else {
		a.Logger.Info("No configuration path provided, skipping configuration")
	}

	return nil
}

// apply takes a file path and applies it to the BindPlane API. If an
// error is found in the response status, it will be returned
func (a *Action) apply(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to read file at path %s: %w", path, err)
	}

	resources := []*model.AnyResource{}
	decoder := yaml.NewDecoder(f)
	for {
		var resource model.AnyResource
		if err := decoder.Decode(&resource); err != nil {
			if err.Error() == "EOF" {
				break
			}
			// TODO(jsirianni): Should we continue and report the error after?
			return fmt.Errorf("resource file %s is malformed, failed to unmarshal yaml: %w", path, err)
		}
		resources = append(resources, &resource)
	}

	if len(resources) == 0 {
		a.Logger.Warn("No resources found in file", zap.String("file", path))
		return nil
	}

	resp, err := a.client.Apply(context.Background(), resources)
	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	if resp == nil {
		// return error indicating that this is a bug with bindplane or the action
		return fmt.Errorf("nil response from client: %s", BugError)
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
			a.Logger.Info("Applied resource", zap.String("name", name), zap.String("status", string(status)))
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

// AutoRollout TODO
func (a *Action) AutoRollout() error {
	// TODO(jsirianni): This should only be configurations managed
	// by a.configurationsPath
	configs, err := a.client.Configurations(context.Background())
	if err != nil {
		return fmt.Errorf("get configurations: %w", err)
	}

	for _, c := range configs {
		status, err := a.client.RolloutStatus(c.Metadata.Name)
		if err != nil {
			return fmt.Errorf("rollout status: %w", err)
		}

		if status.Status.Rollout.Status == model.RolloutStatusPending {
			a.Logger.Info("Pending rollout", zap.String("name", c.Metadata.Name))
		} else {
			a.Logger.Info("No pending rollout", zap.String("name", c.Metadata.Name))
			continue
		}

		a.Logger.Info("Starting rollout", zap.String("name", c.Metadata.Name))

		if err := a.client.StartRollout(c.Metadata.Name); err != nil {
			return fmt.Errorf("start rollout: %w", err)
		}
	}

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
