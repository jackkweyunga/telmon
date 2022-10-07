package main

import (
	"fmt"
	"github.com/reiver/go-telnet"
)

func Monitor(addr string, port int, fromEmail string, password string, receivers []string) {

	conn, err := dial(addr, port)
	if err != nil {

		out := fmt.Sprintf("%v", err)
		fmt.Println("[TelnetClient] Encountered errors while connecting: ")
		fmt.Println(out)

		msg := "Subject: Error Report\r\n" +
			"\r\n" +
			out

		fmt.Println("Reporting this error")
		sendMail(msg, fromEmail, password, receivers)
		fmt.Println("[TelnetClient] Report sent successfully")

	} else {
		fmt.Println("[TelnetClient] Connection successfully established - OK")
		disconnect(conn)
	}

}

func dial(svrAddr string, port int) (*telnet.Conn, error) {

	svrAddrString := fmt.Sprintf("%v:%v", svrAddr, port)

	fmt.Println("[TelnetClient] Trying to connect to telnet server.")
	return telnet.DialTo(svrAddrString)
}

func disconnect(conn *telnet.Conn) {
	fmt.Println("[TelnetClient] Disconnecting ...")
	conn.Close()
	fmt.Println("[TelnetClient] Disconnected - OK")

	return
}
