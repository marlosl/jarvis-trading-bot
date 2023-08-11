package helpers

import (
    "fmt"
    
	"jarvis-trading-bot/cmd/cli/awsdeploy"
	"jarvis-trading-bot/cmd/cli/cloudfare"
)

func UpdateDNS() {
    domainName := awsdeploy.GetStackOutput("regionalDomainName")
    
    fmt.Printf("Updating DNS to %s\n", domainName)
    if domainName != "" {
        cloudfare.UpdateDNS(domainName)
    } else {
        fmt.Println("No DNS update will be made.")
    }
}