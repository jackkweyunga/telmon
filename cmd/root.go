package cmd

import (
	"github.com/jackkweyunga/telmon/monitor"
	"github.com/jackkweyunga/telmon/web"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telmon",
	Short: "Monitor any telnet connection frm a client to a server.",
	Long: `
           _                   
 _        | |                  
| |_  ____| |____   ___  ____  
|  _)/ _  ) |    \ / _ \|  _ \ 
| |_( (/ /| | | | | |_| | | | |
 \___)____)_|_|_|_|\___/|_| |_|
                               
Use telmon to simply monitor telnet connections.

Usage 
telmon [command]
telmon -h for help

example:
________________________________________________________________
start monitoring:		telnet monitor					
start a stats webserver:	telnet server [-p [PORT_NUMBER]]	
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var port int
		var err error

		port = 8080 // default port
		if len(args) > 1 {
			port, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
		}

		go web.Run(port)
		monitor.Play()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().Int("port", 8080, "webserver port")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//serveCmd.Flags().BoolP("port", "p", false, "The port for the webserver to run on.")

}
