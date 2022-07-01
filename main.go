package main

import (
	"os"

	"jarvis-trading-bot/analyzer"
	"jarvis-trading-bot/api"
	"jarvis-trading-bot/broker"
	"jarvis-trading-bot/server"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"
)

func main() {
	arg := os.Args[1:]

	utils.LoadEnvVars()
	log.Init()

	if len(arg) > 0 && arg[0] == "server" {
		server.ServerStartAndListen()
	} else if len(arg) > 0 && arg[0] == "analyzer-server" {
		utils.InitDatabase()
		server.AnalyzerServerStart()
	} else if len(arg) > 0 && arg[0] == "analyzer-simulation-server" {
		utils.InitDatabase()
		server.AnalyzerServerSimulationStart()
	} else if len(arg) > 0 && arg[0] == "api-server" {
		utils.InitDatabase()
		api.StartAndListen()
	} else if len(arg) > 0 && arg[0] == "mock-sender" {
		server.StartMockSenderClient()
	} else if len(arg) > 0 && arg[0] == "balance" {
		symbol := ""
		if len(arg) > 1 {
			symbol = arg[1]
		}
		broker.GetBrokerBalance(symbol)
	} else if len(arg) == 3 && arg[0] == "query" {
		utils.InitDatabase()
		analyzer.QueryOperation(arg[1], arg[2])
	} else if len(arg) > 0 && arg[0] == "import" {
		utils.InitDatabase()
		utils.ImportDataFile()
	} else {
		log.WarningLogger.Println("No service selected. Nothing to do here!")
	}
}
