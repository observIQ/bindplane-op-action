package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/observiq/bindplane-op-action/action"
)

func validate() error {
	if err := validateRemoteURL(); err != nil {
		return err
	}

	if err := validateTargetBranch(); err != nil {
		return err
	}

	if err := validateAuth(); err != nil {
		return err
	}

	if err := validateWriteBack(); err != nil {
		return err
	}

	if err := validateActionsEnvironment(); err != nil {
		return err
	}

	if err := validateFilePaths(); err != nil {
		return err
	}

	return nil
}

func validateRemoteURL() error {
	if bindplane_remote_url == "" {
		return fmt.Errorf("bindplane_remote_url is required")
	}

	u, err := url.Parse(bindplane_remote_url)
	if err != nil {
		return fmt.Errorf("bindplane_remote_url is not a valid URL: %s", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("bindplane_remote_url must be an http or https URL")
	}

	return nil
}

func validateTargetBranch() error {
	if target_branch == "" {
		return fmt.Errorf("target_branch is required")
	}
	return nil
}

func validateAuth() error {
	if bindplane_api_key == "" && bindplane_username == "" {
		return fmt.Errorf("either bindplane_api_key or bindplane_username is required")
	}

	if bindplane_username != "" && bindplane_password == "" {
		return fmt.Errorf("bindplane_password is required when using bindplane_username")
	}

	return nil
}

func validateWriteBack() error {
	if !enable_otel_config_write_back {
		return nil
	}

	if configuration_output_dir == "" {
		return fmt.Errorf("configuration_output_dir is required when enable_otel_config_write_back is true")
	}

	// configuration_output_branch is optional and should be set to
	// target_branch if not provided by the user. If it is empty here,
	// it means we failed to set it in parseArgs or failed to validate
	// that target_branch was set.
	if configuration_output_branch == "" {
		return fmt.Errorf("configuration_output_branch is not set. %s", action.BugError)
	}

	// If a token is not set, github_url is required because it can contain
	// the token.
	if token == "" && github_url == "" {
		return fmt.Errorf("either token or github_url is required when enable_otel_config_write_back is true")
	}

	if github_url != "" {
		p, err := url.Parse(github_url)
		if err != nil {
			return fmt.Errorf("github_url is not a valid URL: %s", err)
		}

		// Action does not support ssh URLs
		if p.Scheme != "http" && p.Scheme != "https" {
			return fmt.Errorf("github_url must be an http or https URL")
		}
	}

	return nil
}

func validateActionsEnvironment() error {
	if os.Getenv("GITHUB_ACTOR") == "" {
		return fmt.Errorf("GITHUB_ACTOR is not set, is the action running in a GitHub runner environment?")
	}

	if os.Getenv("GITHUB_REPOSITORY") == "" {
		return fmt.Errorf("GITHUB_REPOSITORY is not set, is the action running in a GitHub runner environment?")
	}

	return nil
}

func validateFilePaths() error {
	if destination_path != "" {
		if _, err := os.Stat(destination_path); os.IsNotExist(err) {
			return fmt.Errorf("destination_path does not exist: %s", destination_path)
		}
	}

	if source_path != "" {
		if _, err := os.Stat(source_path); os.IsNotExist(err) {
			return fmt.Errorf("source_path does not exist: %s", source_path)
		}
	}

	if processor_path != "" {
		if _, err := os.Stat(processor_path); os.IsNotExist(err) {
			return fmt.Errorf("processor_path does not exist: %s", processor_path)
		}
	}

	if configuration_path != "" {
		if _, err := os.Stat(configuration_path); os.IsNotExist(err) {
			return fmt.Errorf("configuration_path does not exist: %s", configuration_path)
		}
	}

	return nil
}
