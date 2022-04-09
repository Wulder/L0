package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

var conf Config

func main() {

	initConfig()
	fmt.Println("Start handle message at ", conf.LastMessageSequence)
	dbConnStr := "host=" + conf.DBHost + " user=" + conf.User + " dbname=" + conf.DBName + " password=" + conf.Password + " sslmode=disable"

	db = dbConnect(dbConnStr)

	NatsConnect(conf.ClusterID, conf.ClientID)

	http.HandleFunc("/order", OrderViewHandle)
	http.ListenAndServe(":8080", nil)
}

func OrderViewHandle(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()
	if keys["uid"] != nil {
		order := dbGetOrder(keys.Get("uid"))
		if order.OrderUid != "" {
			tmpl, _ := template.ParseFiles("interface.html")
			tmpl.Execute(w, order)
		} else {
			fmt.Fprintf(w, "Заказа с таким UID не найдено")
		}

	}

}

func initConfig() {
	data, _ := ioutil.ReadFile("config.json")

	if err := json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
	if conf.LastMessageSequence > 0 {
		conf.LastMessageSequence += 1
	}

}
