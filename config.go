package main

type Config struct {
	ClusterID           string `json:"clusterId"`
	ClientID            string `json:"clientId"`
	DBHost              string `json:"dbHost"`
	DBPort              int    `json:"dbPort"`
	DBName              string `json:"dbName"`
	User                string `json:"user"`
	Password            string `json:"password"`
	LastMessageSequence uint   `json:"LastMessageSequence`
}
