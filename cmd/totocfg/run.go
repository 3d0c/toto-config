package totocfg

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/apiserver"
	"github.com/3d0c/toto-config/pkg/config"
	"github.com/3d0c/toto-config/pkg/log"
)

const (
	envPrefix = "TOTO_CFG"
)

var (
	globalCtx context.Context
	globalWG  *sync.WaitGroup
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Toto Config API Server",
	Long:  `runs Toto Config API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		log.InitLogger(config.TheConfig().Logger)
		log.TheLogger().Debug("toto-config component",
			zap.String("config", fmt.Sprintf("%#v", config.TheConfig())))

		runProcesses()
		globalWG.Wait()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	viper.SetEnvPrefix(envPrefix)

	var cancelFn func()
	globalCtx, cancelFn = context.WithCancel(context.Background())
	globalWG = &sync.WaitGroup{}

	globalWG.Add(1)
	go signalHandler(cancelFn)
}

func signalHandler(fn func()) {
	defer globalWG.Done()
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.TheLogger().Info("stop execution", zap.String("signal", sig.String()))
	fn()
	close(sigs)
}

func runProcesses() {
	var (
		apiSrv *apiserver.APIHTTPServer
		err    error
	)

	globalWG.Add(1)
	defer globalWG.Done()

	if apiSrv, err = apiserver.NewAPIHTTPServer(config.TheConfig().Server); err != nil {
		log.TheLogger().Fatal("error initializing API server", zap.Error(err))
	}

	apiSrv.Run(globalCtx)
}
