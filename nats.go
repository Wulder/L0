package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"

	"github.com/nats-io/stan.go"
)

var sc stan.Conn

func NatsConnect(cluster string, sub string) {

	sc, _ := stan.Connect(cluster, sub)

	sc.Subscribe("json.orders", HandleMsg, stan.StartAtSequence(uint64(conf.LastMessageSequence)))

}

func HandleMsg(msg *stan.Msg) {

	conf.LastMessageSequence = uint(msg.Sequence)
	configForWriting, err := json.Marshal(conf)
	if err != nil {
		log.Fatalf(err.Error())
	}
	ioutil.WriteFile("config.json", configForWriting, fs.ModeAppend)

	var m Model
	err = json.Unmarshal([]byte(msg.Data), &m)
	if err != nil {
		fmt.Printf("JSON unmarshal error: %s", err.Error())
	} else {
		dbWriteOrder(m)
	}

}
