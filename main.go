package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultConfigFilename = "stingoftheviper"
const envPrefix = "STING"

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	color := ""
	number := 0

	rootCmd := &cobra.Command{
		Use:   "stingoftheviper",
		Short: "Cober and Viper together at last",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string)error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Your favorite color is:", color)
			fmt.Println("The magic number is:", number)
		},
	}

	rootCmd.Flags().IntVarP(&number, "number", "n", 7, "What is the magic number?")
	rootCmd.Flags().StringVarP(&color, "favorite-color", "c", "blue","Should come from flag first, then env var DEBUG_STUFF then the config file, then the default last")

	return rootCmd
}

func initializeConfig(cmd *cobra.Command)error {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

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