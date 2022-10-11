package web

import (
	"bufio"
	"encoding/json"
	"log"
)

type Stats struct {
	Count int            `json:"count"`
	Data  []LogStructure `json:"data"`
}

type LogStructure struct {
	Kind  string `json:"kind"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

func processData(s *bufio.Scanner) Stats {
	s.Split(bufio.ScanLines)

	var lines []string
	var t Stats

	for s.Scan() {
		var l LogStructure

		lines = append(lines, s.Text())

		err := json.Unmarshal(s.Bytes(), &l)
		if err != nil {
			log.Fatal(err)
		}
		if l.Kind != "" {
			t.Data = append(t.Data, l)
		}
	}

	t.Count = len(lines)

	//fmt.Println("Number of lines", len(lines))
	return t
}
