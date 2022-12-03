package totocfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/config"
	"github.com/3d0c/toto-config/pkg/log"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "toto-config",
	Short: "Toto Config API Server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Config file hasn't been provided")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		panic(fmt.Sprintf("Error reading config file from '%s' - %s", cfgFile, err))
	}

	if err := viper.Unmarshal(config.TheConfig()); err != nil {
		panic(fmt.Sprintf("Failed to init config: %s", err))
	}

	// Override config database DSN if it's provided by ENV
	if dsn := os.Getenv("TOTO_DSN"); dsn != "" {
		config.TheConfig().Database.DSN = dsn
	}

	checkRequired(
		// http server required configs
		"server.apiVersion",
	)
}

func checkRequired(keys ...string) {
	for _, k := range keys {
		if !viper.IsSet(k) {
			log.TheLogger().Error("missing required configuration", zap.String("key", k))
			os.Exit(1)
		}
	}
}
