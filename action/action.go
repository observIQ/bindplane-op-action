package action

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/observiq/bindplane-op-action/action/state"
	"github.com/observiq/bindplane-op-action/internal/client"
	"github.com/observiq/bindplane-op-action/internal/client/config"
	"github.com/observiq/bindplane-op-action/internal/client/model"
	"github.com/observiq/bindplane-op-action/internal/client/version"
	"github.com/observiq/bindplane-op-action/internal/repo"
	"gopkg.in/yaml.v3"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

const (
	BugError = "This is a bug with the action, please reach out to support or file an issue on Github https://github.com/observIQ/bindplane-op-action/issues"
)

type rType string

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
	action.state = state.NewMemory()

	return action, nil
}

// Action is a struct that contains the BindPlane client
// and user defined configuration options
type Action struct {
	Logger *zap.Logger

	// Branch name and paths to read
	// resources from
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

	// State holds the current state of the action
	state state.State
}

// TestConnection wraps the BindPlane client's Version method
func (a *Action) TestConnection() (version.Version, error) {
	v, err := a.client.Version(context.Background())
	if err != nil {
		return version.Version{}, fmt.Errorf("failed to test connection: %w", err)
	}
	return v, err
}

// Run executes the action
func (a *Action) Run() error {
	if err := a.Apply(); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	if a.autoRollout {
		if err := a.AutoRollout(); err != nil {
			return fmt.Errorf("failed to rollout configuration: %s", err)
		}
	}

	if a.enableWriteBack {
		if err := a.WriteBack(); err != nil {
			return fmt.Errorf("failed to write back configuration: %s", err)
		}
	}

	return nil
}

// RunRollout progresses a rollout for a configuration
func (a *Action) RunRollout(config string) error {
	if err := a.client.StartRollout(config); err != nil {
		return fmt.Errorf("start rollout: %w", err)
	}

	return nil
}

// Apply applies destinations, sources, processors, and configurations
// in that order. It is important to apply destinations first, followed
// by resource library sources and processors. Configurations should be
// applied last because they will reference other resources.
func (a *Action) Apply() error {
	if a.destinationPath != "" {
		a.Logger.Info("Applying resources", zap.String("Kind", string(model.KindDestination)), zap.String("file", a.destinationPath))
		err := a.apply(a.destinationPath)
		if err != nil {
			return fmt.Errorf("destinations: %w", err)
		}
	} else {
		a.Logger.Info("No destination path provided, skipping destinations")
	}

	if a.sourcePath != "" {
		a.Logger.Info("Applying resources", zap.String("Kind", string(model.KindSource)), zap.String("file", a.sourcePath))
		err := a.apply(a.sourcePath)
		if err != nil {
			return fmt.Errorf("sources: %w", err)
		}
	} else {
		a.Logger.Info("No source path provided, skipping sources")
	}

	if a.processorPath != "" {
		a.Logger.Info("Applying resources", zap.String("Kind", string(model.KindProcessor)), zap.String("file", a.processorPath))
		err := a.apply(a.processorPath)
		if err != nil {
			return fmt.Errorf("processors: %w", err)
		}
	} else {
		a.Logger.Info("No processor path provided, skipping processors")
	}

	if a.configurationPath != "" {
		a.Logger.Info("Applying resources", zap.String("Kind", string(model.KindConfiguration)), zap.String("file", a.configurationPath))
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
	resources, err := decodeAnyResourceFile(path)
	if err != nil {
		return fmt.Errorf("decode resources: %w", err)
	}

	resp, err := a.client.Apply(context.Background(), resources)
	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	if resp == nil {
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

		// Attach the configuration resource to the state
		// so we can use it for auto rollout
		if kind == string(model.KindConfiguration) {
			a.state.SetConfiguration(s.Resource.Metadata.Name, s.Resource)
			a.Logger.Debug("Configuration resource added to state", zap.String("name", name))
		}

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
	configurations := []model.Configuration{}
	for _, name := range a.state.ConfigurationNames() {
		configuration, err := a.client.Configuration(context.Background(), name)
		if err != nil {
			return fmt.Errorf("get configuration %s: %w", name, err)
		}
		if configuration == nil {
			return fmt.Errorf("configuration '%s' is nil: %s", name, BugError)
		}
		configurations = append(configurations, *configuration)
	}

	for _, c := range configurations {
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
	a.Logger.Info(
		"Cloning repository", zap.String("branch", a.configurationOutputBranch),
	)

	repo, err := repo.CloneRepo(a.githubURL, a.configurationOutputBranch, a.githubToken)
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}

	tree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree: %w", err)
	}

	rawConfigs := make(map[string]string)
	for _, name := range a.state.ConfigurationNames() {
		rawConfig, err := a.client.RawConfiguration(context.Background(), name)
		if err != nil {
			return fmt.Errorf("get configuration %s: %w", name, err)
		}
		if rawConfig == "" {
			return fmt.Errorf("configuration '%s' is empty: %s", name, BugError)
		}
		rawConfigs[name] = rawConfig
	}

	for name, rawConfig := range rawConfigs {
		path := fmt.Sprintf("./out_repo/%s/%s.yaml", a.configurationOutputDir, name)

		// Create the directory if it doesn't exist. MkdirAll will return
		// nil if the directory already exists. Returns an error if something
		// goes wrong.
		dir, _ := filepath.Split(path)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600) // #nosec G304 user defined filepath
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		defer f.Close()

		if _, err := f.WriteString(rawConfig); err != nil {
			return fmt.Errorf("write file %s: %w", path, err)
		}

		a.Logger.Info("Raw configuration written to file", zap.String("name", name), zap.String("path", path))
	}

	status, err := tree.Status()
	if err != nil {
		return fmt.Errorf("get work tree status: %w", err)
	}

	if status.IsClean() {
		a.Logger.Info("No changes to write back")
		return nil
	}

	a.Logger.Info("Detected changes, writing back to repository")
	for path := range status {
		a.Logger.Info("file changed", zap.String("path", path))
		_, err := tree.Add(path)
		if err != nil {
			return fmt.Errorf("git add file %s: %w", path, err)
		}
	}

	commitMessage := "BindPlane OP Action: Update OTEL Configs"
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bindplane-op-action",
			Email: "bindplane-op-action",
			When:  time.Now(),
		},
	}
	_, err = tree.Commit(commitMessage, commitOptions)
	if err != nil {
		return fmt.Errorf("commit changes: %w", err)
	}

	pushOpts := &git.PushOptions{
		RemoteName: "origin",
	}

	if err = repo.Push(pushOpts); err != nil {
		return fmt.Errorf("push changes: %w", err)
	}

	a.Logger.Info("Changes written back to repository")

	return nil
}

// decodeAnyResourceFile takes a file path and decodes it into a slice of
// model.AnyResource. If the file is empty, it will return an error.
// This function supports globbing, but does not gaurantee ordering. This
// function should not be passed multiple files with differing resource
// types such as Destinations and Configurations.
func decodeAnyResourceFile(path string) ([]*model.AnyResource, error) {
	// Glob will return nil matches if there are IO errors. Glob only returns
	// an error if an invalid pattern is given.
	matches, err := filepath.Glob(path) // #nosec G304 user defined filepath
	if err != nil {
		return nil, fmt.Errorf("glob path %s: %w", path, err)
	}
	if matches == nil {
		return nil, fmt.Errorf("no matching files found when globbing %s", path)
	}

	resources := []*model.AnyResource{}

	for _, match := range matches {
		f, err := os.Open(match) // #nosec G304 user defined filepath
		if err != nil {
			return nil, fmt.Errorf("unable to read file at path %s: %w", path, err)
		}
		defer f.Close()

		decoder := yaml.NewDecoder(f)
		for {
			resource := &model.AnyResource{}
			if err := decoder.Decode(resource); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				// TODO(jsirianni): Should we continue and report the error after?
				return nil, fmt.Errorf("resource file %s is malformed, failed to unmarshal yaml: %w", path, err)
			}
			resources = append(resources, resource)
		}
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no resources found in file: %s", path)
	}

	return resources, nil
}
