package main_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("PULUMI_CONFIG_PASSPHRASE", "")
	os.Setenv("PULUMI_SKIP_UPDATE_CHECK", "true")

	cmd := exec.Command("pulumi", "login", "--local")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()

	cmd = exec.Command("pulumi", "logout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	os.Unsetenv("PULUMI_CONFIG_PASSPHRASE")
	os.Unsetenv("PULUMI_SKIP_UPDATE_CHECK")

	os.Exit(exitCode)
}

func TestLocalProject(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_new_stack_local_source_secrets", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceSecrets"
		workDir := filepath.Join(".", "localproject_main")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			os.Unsetenv("FOO")
			os.Unsetenv("BAZ")
			os.Unsetenv("BAR_TOKEN")
			err := s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		s.Workspace().SetEnvVars(map[string]string{
			"FOO": "test-foo",
			"BAZ": "test-baz",
		})

		envvars := s.Workspace().GetEnvVars()
		require.Equal(t, "test-foo", envvars["FOO"])
		require.Equal(t, "test-baz", envvars["BAZ"])

		err = s.SetAllConfig(ctx, auto.ConfigMap{
			"bar:token": auto.ConfigValue{
				Value:  "test-bar-token",
				Secret: true,
			},
			"buzz:owner": auto.ConfigValue{
				Value:  "xyz",
				Secret: true,
			},
		})
		require.NoError(t, err)

		values, err := s.GetAllConfig(ctx)
		require.NoError(t, err)
		require.Equal(t, "test-bar-token", values["bar:token"].Value)
		require.True(t, values["bar:token"].Secret)
		require.Equal(t, "xyz", values["buzz:owner"].Value)
		require.True(t, values["buzz:owner"].Secret)
	})
	t.Run("Test_new_stack_local_source_config", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceConfig"
		workDir := filepath.Join(".", "localproject_main")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		require.Equal(t, stackName, s.Name())

		err = s.SetAllConfig(ctx, auto.ConfigMap{
			"bar:token": auto.ConfigValue{
				Value:  "test-bar-token",
				Secret: true,
			},
			"buzz:owner": auto.ConfigValue{
				Value:  "xyz",
				Secret: true,
			},
		})
		require.NoError(t, err)

		values, err := s.GetAllConfig(ctx)
		require.NoError(t, err)
		require.Equal(t, "test-bar-token", values["bar:token"].Value)
		require.True(t, values["bar:token"].Secret)
		require.Equal(t, "xyz", values["buzz:owner"].Value)
		require.True(t, values["buzz:owner"].Secret)
	})
	t.Run("Test_local_source_workspace_env_vars", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceWorkspaceEnvVars"
		workDir := filepath.Join(".", "localproject_main")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)

			s.Workspace().UnsetEnvVar("FOO")
			s.Workspace().UnsetEnvVar("BAZ")
		}()

		require.Equal(t, stackName, s.Name())

		err = s.Workspace().SetEnvVars(map[string]string{
			"FOO": "BAR",
			"BAZ": "QUX",
		})
		require.NoError(t, err)

		envvars := s.Workspace().GetEnvVars()
		require.Equal(t, "BAR", envvars["FOO"])
		require.Equal(t, "QUX", envvars["BAZ"])

		s.Workspace().UnsetEnvVar("FOO")
		s.Workspace().UnsetEnvVar("BAZ")

		envvars = s.Workspace().GetEnvVars()
		require.NotContains(t, envvars, "FOO")
		require.NotContains(t, envvars, "BAZ")
	})
	t.Run("Test_pulumi_whoami", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		cmd := exec.Command("pulumi", "whoami")
		cmd.Stdout = &out
		cmd.Stderr = &out

		err := cmd.Run()
		require.NoError(t, err)
		username := strings.TrimSpace(out.String())

		var expectedOut bytes.Buffer
		expectedCmd := exec.Command("whoami")
		expectedCmd.Stdout = &expectedOut
		expectedCmd.Stderr = &expectedOut

		err = expectedCmd.Run()
		require.NoError(t, err)
		expectedUsername := strings.TrimSpace(expectedOut.String())

		require.Equal(t, expectedUsername, username)
	})
}
