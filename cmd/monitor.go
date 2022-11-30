package cmd

import (
	"fmt"
	"github.com/jackkweyunga/telmon/monitor"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "start the monitoring service",
	Long: `
           _                   
 _        | |                  
| |_  ____| |____   ___  ____  
|  _)/ _  ) |    \ / _ \|  _ \ 
| |_( (/ /| | | | | |_| | | | |
 \___)____)_|_|_|_|\___/|_| |_|
                              
Starts a monitoring service without starting the webserver.

usage:
telmon monitor

help/description:
telmon monitor -h

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
           _                   
 _        | |                  
| |_  ____| |____   ___  ____  
|  _)/ _  ) |    \ / _ \|  _ \ 
| |_( (/ /| | | | | |_| | | | |
 \___)____)_|_|_|_|\___/|_| |_|
                              
[Telmon] All logs are here telmon.log
[Telmon] Running telmon without logs api.
[Telmon] To access prometheus stats api(/prometheus) run command: telmon
`)

		monitor.Play()
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
