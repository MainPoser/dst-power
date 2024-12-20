package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	options_db "github.com/MainPoser/dst-power/pkg/db/options"
	"github.com/MainPoser/dst-power/pkg/server/options"
	"github.com/MainPoser/dst-power/pkg/util/signal"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/MainPoser/dst-power/internal/router"
)

// Options 包含运行服务的所有信息.
type Options struct {
	// configFile 配置文件路径.
	configFile string
	// gin启动模式
	ginMode string

	// 数据库配置
	db *options_db.DbOptions
	// 日志配置
	logLevel int
	// http配置
	serverOpt *options.ServerOptions
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	dbFlagSet := pflag.NewFlagSet("db", pflag.ExitOnError)
	dbFlagSet.StringVarP(&o.db.Url, "db-url", "", "mysql://admin:a123456@127.0.0.1:3306", "db url")
	fs.AddFlagSet(dbFlagSet)

	logFlagSet := pflag.NewFlagSet("log", pflag.ExitOnError)
	logFlagSet.IntVarP(&o.logLevel, "log-level", "v", 0, "log level")
	fs.AddFlagSet(logFlagSet)

	serverFlagSet := pflag.NewFlagSet("server", pflag.ExitOnError)
	serverFlagSet.IPVarP(&o.serverOpt.BindAddress, "bind-address", "", net.IP{}, "bindAddress")
	serverFlagSet.IntVarP(&o.serverOpt.BindPort, "bind-port", "", 8080, "bindPort")
	fs.AddFlagSet(serverFlagSet)

	fs.StringVarP(&o.configFile, "conf", "c", "~/.config", "config file path")
	fs.StringVarP(&o.ginMode, "gin-mode", "", gin.DebugMode, "ginMode")

}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) Complete() error {
	return nil
}

func (o *Options) Run(stopCh <-chan struct{}) error {
	logrus.SetLevel(logrus.Level(o.logLevel))
	logrus.Infof("%+v\n", *o)

	s := &http.Server{
		Addr:           o.serverOpt.BindAddress.String() + ":" + strconv.Itoa(o.serverOpt.BindPort),
		Handler:        router.NewRouter(o.ginMode),
		ReadTimeout:    time.Duration(o.serverOpt.ReadTimeout),
		WriteTimeout:   time.Duration(o.serverOpt.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()
	<-stopCh
	logrus.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Println("Server exiting")
	return nil
}

func NewServerRunOptions() *Options {
	return &Options{
		configFile: "",
		db:         &options_db.DbOptions{},
		logLevel:   0,
		serverOpt:  &options.ServerOptions{},
	}
}

func NewPowerCommand() *cobra.Command {
	o := NewServerRunOptions()
	cmd := &cobra.Command{
		Use:  "dst-power",
		Long: `The dst-power server.`,
		// 程序出错不打印 usage
		SilenceUsage: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// set default o
			if err := o.Complete(); err != nil {
				return err
			}

			// validate o
			if errs := o.Validate(); len(errs) != 0 {
				return errors.Join(errs...)
			}
			stopCh := signal.SetupSignalHandler()

			return o.Run(stopCh)
		},
	}

	fs := cmd.Flags()
	o.AddFlags(fs)

	return cmd
}
