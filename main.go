package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"task/bot"
	"task/config"
	"task/database"
)

func main() {
	// flags
	flag.String("conf", `example_configuration.json`, "Configuration file path")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	conf := config.GetConfig()

	database.ConnectDatabase(conf)

	go bot.StartBot(conf.BotAPIKey)

	select {}

}
