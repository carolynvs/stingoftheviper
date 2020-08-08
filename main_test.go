package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory
	tmpDir, err := ioutil.TempDir("", "stingoftheviper")
	require.NoError(t, err, "error creating a temporary test directory")
	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(testDir)
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")

	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		configB, err := ioutil.ReadFile(filepath.Join(testDir, "stingoftheviper.toml"))
		require.NoError(t, err, "error reading test config file")
		err = ioutil.WriteFile(filepath.Join(tmpDir, "stingoftheviper.toml"), configB, 0644)
		require.NoError(t, err, "error writing test config file")
		defer os.Remove(filepath.Join(tmpDir, "stingoftheviper.toml"))

		// Run ./stingoftheviper
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `Your favorite color is: blue
The magic number is: 7
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color from the config file and the number from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		// Run STING_FAVORITE_COLOR=purple ./stingoftheviper
		os.Setenv("STING_FAVORITE_COLOR", "purple")
		defer os.Unsetenv("STING_FAVORITE_COLOR")

		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `Your favorite color is: purple
The magic number is: 7
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color to use the environment variable value and the number to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./stingoftheviper --number 2
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--number", "2"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `Your favorite color is: red
The magic number is: 2
`
		assert.Equal(t, wantOutput, gotOutput, "expected the number to use the flag value and the color to use the flag default")
	})
}
