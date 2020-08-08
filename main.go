package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "stingoftheviper"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "STING"
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	// Store the result of binding cobra flags and viper config. In a
	// real application these would be data structures, most likely
	// custom structs per command. This is simplified for the demo app and is
	// not recommended that you use one-off variables. The point is that we
	// aren't retrieving the values directly from viper or flags, we read the values
	// from standard Go data structures.
	color := ""
	number := 0

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "stingoftheviper",
		Short: "Cober and Viper together at last",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			// Print the final resolved value from binding cobra flags and viper config
			fmt.Fprintln(out, "Your favorite color is:", color)
			fmt.Fprintln(out, "The magic number is:", number)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant precedence)
	rootCmd.Flags().IntVarP(&number, "number", "n", 7, "What is the magic number?")
	rootCmd.Flags().StringVarP(&color, "favorite-color", "c", "red", "Should come from flag first, then env var DEBUG_STUFF then the config file, then the default last")

	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Bind flags to environment variables
	// Works great for simple flags, but needs help for flags
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --debug-plugins binds to PORTER_DEBUG_PLUGINS
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the configuration file value to the flag when the flag is not set
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
