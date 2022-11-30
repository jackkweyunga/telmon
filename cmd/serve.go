package cmd

import (
	"fmt"
	"github.com/jackkweyunga/telmon/web"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the stats webserver.",
	Long: `
           _                   
 _        | |                  
| |_  ____| |____   ___  ____  
|  _)/ _  ) |    \ / _ \|  _ \ 
| |_( (/ /| | | | | |_| | | | |
 \___)____)_|_|_|_|\___/|_| |_|
                              
Start a webserver running at a passed port. Default port is 8080.

usage:
telmon serve -p 9090

help/description:
telmon serve -h
`,
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
		fmt.Printf(`
           _                   
 _        | |                  
| |_  ____| |____   ___  ____  
|  _)/ _  ) |    \ / _ \|  _ \ 
| |_( (/ /| | | | | |_| | | | |
 \___)____)_|_|_|_|\___/|_| |_|
                              
Started a webserver listening at port %v `, port)

		web.Run(port)

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serveCmd.PersistentFlags().String("port", "", "The port for the webserver to run on.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//serveCmd.Flags().BoolP("port", "p", false, "The port for the webserver to run on.")
}
