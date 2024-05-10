package action

import (
	"testing"

	"github.com/observiq/bindplane-op-action/client/config"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

func TestWithBindPlaneRemoteURL(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set remote URL",
			"http://localhost:3001",
			&Action{
				config: config.Config{
					Network: config.Network{
						RemoteURL: "http://localhost:3001",
					},
				},
			},
		},
		{
			"Set remote URL ipv6",
			"http://[::1]:3001",
			&Action{
				config: config.Config{
					Network: config.Network{
						RemoteURL: "http://[::1]:3001",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithBindPlaneRemoteURL(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}

}

func TestWithBindPlaneAPIKey(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set API key",
			"123456",
			&Action{
				config: config.Config{
					Auth: config.Auth{
						APIKey: "123456",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithBindPlaneAPIKey(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithBindPlaneUsername(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set username",
			"admin",
			&Action{
				config: config.Config{
					Auth: config.Auth{
						Username: "admin",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithBindPlaneUsername(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithBindPlanePassword(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set password",
			"password",
			&Action{
				config: config.Config{
					Auth: config.Auth{
						Password: "password",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithBindPlanePassword(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithTLSCACert(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set certificate authority",
			"ca.pem",
			&Action{
				config: config.Config{
					Network: config.Network{
						CertificateAuthority: []string{"ca.pem"},
					},
				},
			},
		},
		{
			"Empty certificate authority",
			"",
			&Action{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithTLSCACert(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithDestinationPath(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set destination path",
			"/tmp",
			&Action{
				destinationPath: "/tmp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithDestinationPath(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithSourcePath(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set source path",
			"/tmp",
			&Action{
				sourcePath: "/tmp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithSourcePath(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithProcessorPath(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set processor path",
			"/tmp",
			&Action{
				processorPath: "/tmp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithProcessorPath(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithConfigurationPath(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set configuration path",
			"/tmp",
			&Action{
				configurationPath: "/tmp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithConfigurationPath(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithEnableAutoRollout(t *testing.T) {
	cases := []struct {
		name   string
		intput bool
		expect *Action
	}{
		{
			"Enable auto rollout",
			true,
			&Action{
				autoRollout: true,
			},
		},
		{
			"Disable auto rollout",
			false,
			&Action{
				autoRollout: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithAutoRollout(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithOTELConfigWriteBack(t *testing.T) {
	cases := []struct {
		name   string
		intput bool
		expect *Action
	}{
		{
			"Enable write back",
			true,
			&Action{
				enableWriteBack: true,
			},
		},
		{
			"Disable write back",
			false,
			&Action{
				enableWriteBack: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithOTELConfigWriteBack(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithConfigurationOutputDir(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set output directory",
			"/tmp",
			&Action{
				configurationOutputDir: "/tmp",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithConfigurationOutputDir(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithConfigurationOutputBranch(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set output branch",
			"main",
			&Action{
				configurationOutputBranch: "main",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithConfigurationOutputBranch(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithGithubToken(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set token",
			"123456",
			&Action{
				githubToken: "123456",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithGithubToken(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestWithGithubURL(t *testing.T) {
	cases := []struct {
		name   string
		intput string
		expect *Action
	}{
		{
			"Set URL",
			"git@github.com:org/repo.git",
			&Action{
				githubURL: "git@github.com:org/repo.git",
			},
		},
		{
			"Set Advanced URL",
			"git:token@github.mycorp.net:org/repo.git",
			&Action{
				githubURL: "git:token@github.mycorp.net:org/repo.git",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &Action{}
			opt := WithGithubURL(tc.intput)
			opt(a)
			require.Equal(t, tc.expect, a)
		})
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		name   string
		opts   []Option
		expect *Action
		errStr string
	}{
		{
			"Basic",
			[]Option{
				WithBindPlaneRemoteURL("http://localhost:3001"),
				WithBindPlaneUsername("admin"),
				WithBindPlanePassword("password"),
				WithAutoRollout(false),
				WithOTELConfigWriteBack(false),
			},
			&Action{
				config: config.Config{
					Network: config.Network{
						RemoteURL: "http://localhost:3001",
					},
					Auth: config.Auth{
						Username: "admin",
						Password: "password",
					},
				},
				autoRollout:     false,
				enableWriteBack: false,
			},
			"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := New(zap.NewNop(), tc.opts...)
			if tc.errStr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errStr)
				return
			}

			// Set the client and logger to nil, we cannot compare the client directly.
			// We can make sure New did not return an error. Integration tests
			// will verify the client is created correctly.
			a.client = nil
			a.Logger = nil
			a.state = nil // TODO(jsirianni): Add state tests

			require.NoError(t, err)
			require.Equal(t, tc.expect, a)

		})
	}
}
