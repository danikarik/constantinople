package main

import (
	"flag"
	"fmt"

	"github.com/danikarik/constantinople/pkg/app"
	"github.com/danikarik/constantinople/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Version stands for build version.
const Version = "v0.0.1"

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Start http server",
		Run: func(cmd *cobra.Command, args []string) {
			server()
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()
	rootCmd.AddCommand(versionCmd)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.SetEnvPrefix("CBKZ")
	viper.AutomaticEnv()
	viper.SetDefault("ADDR", ":3000")
	viper.SetDefault("AUTH_SERV", "127.0.0.1:8000")
	viper.SetDefault("REDIS_HOST", "127.0.0.1:6379")
	viper.SetDefault("REDIS_PASS", "")
	viper.SetDefault("ORIGINS", []string{"*"})
	viper.SetDefault("USERNAME", "carbase")
	viper.SetDefault("PASSWORD", "FtgVCxZVZeRk5K9S")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		util.Exit("%s", err.Error())
	}
}

func server() {
	app, err := app.New(viper.GetString("ADDR"), app.Options{
		Origins:     viper.GetStringSlice("ORIGINS"),
		AuthService: viper.GetString("AUTH_SERV"),
		RedisHost:   viper.GetString("REDIS_HOST"),
		RedisPass:   viper.GetString("REDIS_PASS"),
		Username:    viper.GetString("USERNAME"),
		Password:    viper.GetString("PASSWORD"),
	})
	if err != nil {
		util.Exit("%s", err.Error())
	}
	if err = app.Serve(); err != nil {
		util.Exit("%s", err.Error())
	}
}
