package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRemoteURL(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect error
	}{
		{
			"Valid URL",
			"http://localhost:3001",
			nil,
		},
		{
			"Valid URL ipv6",
			"http://[::1]:3001",
			nil,
		},
		{
			"Invalid URL",
			"localhost:3001",
			errors.New("bindplane_remote_url must be an http or https URL"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bindplane_remote_url = tc.input
			defer func() {
				bindplane_remote_url = ""
			}()
			err := validateRemoteURL()
			require.Equal(t, tc.expect, err)
		})
	}
}

func TestValidateTargetBranch(t *testing.T) {
	require.Error(t, validateTargetBranch(), "target_branch is required")
	target_branch = "main"
	defer func() {
		target_branch = ""
	}()
	require.NoError(t, validateTargetBranch())
}

func TestValidateAuth(t *testing.T) {
	cases := []struct {
		name string
		key  string
		user string
		pass string
		err  error
	}{
		{
			"Missing key and user",
			"",
			"",
			"",
			errors.New("either bindplane_api_key or bindplane_username is required"),
		},
		{
			"Valid key",
			"key",
			"",
			"",
			nil,
		},
		{
			"Valid user",
			"",
			"user",
			"pass",
			nil,
		},
		{
			"User without password",
			"",
			"user",
			"",
			errors.New("bindplane_password is required when using bindplane_username"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bindplane_api_key = tc.key
			bindplane_username = tc.user
			bindplane_password = tc.pass
			defer func() {
				bindplane_api_key = ""
				bindplane_username = ""
				bindplane_password = ""
			}()

			err := validateAuth()
			require.Equal(t, tc.err, err)
		})
	}
}

func ValidateWriteBack(t *testing.T) {
	cases := []struct {
		name            string
		enableWriteBack bool
		outputDir       string
		outputBranch    string
		token           string
		githubURL       string
		err             error
	}{
		{
			"Valid write back with token",
			true,
			"./output",
			"main",
			"token",
			"",
			nil,
		},
		{
			"Valid write back with github url",
			true,
			"./output",
			"main",
			"",
			"https://token@git.corp.net:8433",
			nil,
		},
		{
			"Missing token and github url",
			true,
			"./output",
			"main",
			"",
			"",
			errors.New("either token or github_url is required when enable_otel_config_write_back is true"),
		},
		{
			"Missing output dir",
			true,
			"",
			"main",
			"token",
			"",
			errors.New("configuration_output_dir is required when enable_otel_config_write_back is true"),
		},
		{
			"Valid output branch",
			true,
			"./output",
			"",
			"token",
			"",
			errors.New("configuration_output_branch is required when enable_otel_config_write_back is true"),
		},
		{
			"Invalid github url scheme",
			true,
			"./output",
			"main",
			"",
			"git@git.corp.net:8433",
			errors.New("github_url must be an https URL"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enable_otel_config_write_back = tc.enableWriteBack
			configuration_output_dir = tc.outputDir
			configuration_output_branch = tc.outputBranch
			token = tc.token
			github_url = tc.githubURL
			defer func() {
				enable_otel_config_write_back = false
				configuration_output_dir = ""
				configuration_output_branch = ""
				token = ""
				github_url = ""
			}()

			require.Equal(t, tc.err, validateWriteBack())

		})
	}
}

func TestValidateActionsEnvironment(t *testing.T) {
	require.Error(t, validateActionsEnvironment())

	os.Setenv("GITHUB_ACTOR", "actor")
	defer os.Unsetenv("GITHUB_ACTOR")

	require.Error(t, validateActionsEnvironment())

	os.Setenv("GITHUB_REPOSITORY", "repo")
	defer os.Unsetenv("GITHUB_REPOSITORY")

	require.NoError(t, validateActionsEnvironment())
}
