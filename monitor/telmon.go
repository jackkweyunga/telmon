package monitor

import (
	"fmt"
	log "github.com/jackkweyunga/telmon/logging"
	"github.com/reiver/go-telnet"
	"os/exec"
)

func Monitor(addr string, port int, n Notify, r Restore) {

	conn, err := dial(addr, port)
	if err != nil {

		out := fmt.Sprintf("%v", err)
		log.Log.Println("[Telmon] Encountered errors while connecting: ")
		log.FailLogger.Error("[Telmon] %v \n", out)
		TotalFailures.WithLabelValues("failed").Inc()

		msg := "Subject: Error Report\r\n" +
			"\r\n" +
			out

		log.Log.Println("[Telmon] Reporting this error")
		n.report(msg)

		// try to restore connection
		if r.Cmd != "" {
			restoreConnection(r.Exec, r.Cmd, addr, port)
		}

	} else {
		log.Log.Println("[Telmon] Connection successfully established - OK")
		log.StableLogger.Info("[Telmon] connection is stable")
		TotalPasses.WithLabelValues("pass").Inc()
		disconnect(conn)
	}

}

func dial(svrAddr string, port int) (*telnet.Conn, error) {

	svrAddrString := fmt.Sprintf("%v:%v", svrAddr, port)

	log.Log.Println("[Telmon] Trying to connect to telnet server.")
	return telnet.DialTo(svrAddrString)
}

func disconnect(conn *telnet.Conn) {
	log.Log.Println("[Telmon] Disconnecting ...")
	err := conn.Close()
	if err != nil {
		log.Log.Fatalln(err)
	}
	log.Log.Println("[TelnetClient] Disconnected - OK")

	return
}

func restoreConnection(x string, cmd string, addr string, port int) {
	// run restoration command
	log.Log.Printf("[Telmon] Running restoration command: %v -c %v \n", x, cmd)

	c := exec.Command(x, "-c", cmd)
	_, err := c.CombinedOutput()
	if err != nil {
		log.Log.Printf("[Telmon] Restoration failed, command returned error: %v", err)
		return
	}

	// check if it's restored
	conn, err := dial(addr, port)
	if err != nil {
		log.Log.Printf("[Telmon] Restoration failed, Dial returned error: %v", err)
		return
	} else {
		log.RestoredLogger.Info("[Telmon] Restoration succeeded")
		disconnect(conn)
	}
}
