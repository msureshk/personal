package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Config structure
type Config struct {
	Server struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"server"`

	Bing struct {
		Endpoint     string `json:"endpoint"`
		Token        string `json:"token"`
		Resultlimit  int    `json:"resultlimit"`
		Searchsuffix string `json:"searchsuffix"`
	} `json:"bing"`
	Colly struct {
		Maxdepth            int    `json:"maxdepth"`
		Parallelism         int    `json:"parallelism"`
		Exclusionsregex     string `json:"exclusionsregex"`
		Disallowedurlsregex string `json:"Disallowedurlsregex"`
	} `json:"colly"`
}

//Cfg global variable for the config structure
var Cfg Config

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		qt := r.URL.Query().Get("qt")
		if qt == "" {
			return
		}
		searchterm := qt + " " + Cfg.Bing.Searchsuffix
		log.Println("visiting", searchterm)
		var urllist []string

		//search for the query term using bing web search

		urllist, err := search(searchterm, &Cfg)
		if err != nil {
			log.Println("failed to get response:", err)
			return
		}

		if err != nil {
			log.Println("failed to serialize response:", err)
			return
		}
		jsonyamllinks, err := crawlurls(urllist, &Cfg)
		if err != nil {
			log.Println("failed to get crawl response:", err)
			return
		}
		c, err := json.Marshal(jsonyamllinks)

		if err != nil {
			log.Println("failed to serialize response:", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(c)
	}
}

func readFile(cfg *Config) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Println("failed to open config:", err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	err = json.Unmarshal(byteValue, cfg)
	if err != nil {
		log.Println("failed to decode config:", err)
	}
}

//load the config and listens to the port
func main() {

	readFile(&Cfg)

	http.HandleFunc("/", handler)
	addr := Cfg.Server.Host + ":" + Cfg.Server.Port
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
