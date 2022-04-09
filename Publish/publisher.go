package main

import (
	"fmt"
	"io/ioutil"
	"log"

	stan "github.com/nats-io/stan.go"
)

func main() {

	var sc = tryConnect()

	var fileName string
	for {
		fmt.Println("Enter a publish file name")
		fmt.Scanf("%s\n", &fileName)

		data, _ := ioutil.ReadFile(fileName)

		if err := sc.Publish("json.orders", data); err != nil {
			log.Println(err.Error())
		} else {
			log.Println("Order success sent!")
		}

	}

}

func tryConnect() stan.Conn {

	sc, err := stan.Connect("L0", "publisher")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("STAN connection established")
	}

	return sc
}
