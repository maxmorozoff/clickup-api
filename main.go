package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"bytes"
)
func main() {
	auth_this := os.Getenv("AUTH_THIS")
	clickup_auth := os.Getenv("CLICKUP_AUTH")
	clickup_list_id := os.Getenv("CLICKUP_LIST_ID")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	
	clickup_api := "https://api.clickup.com/api/v2/list/"+clickup_list_id+"/task"

	client := &http.Client{}
	
	http.HandleFunc("/create-task", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/create-task", "host:", r.Host, "ip:", r.RemoteAddr)
		if r.Method == "POST" && r.Header.Get("Authorization") == auth_this {
			log.Println("[AUTH_OK]")

			body, err := ioutil.ReadAll(r.Body) 
			if err != nil {
				log.Println("Errored when reading body of incoming request")
				return
			}

			req, _ := http.NewRequest("POST", clickup_api , bytes.NewBuffer(body))

			req.Header.Add("Authorization", clickup_auth)
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)

			if err != nil {
				log.Println("Errored when sending request to the server")
				w.WriteHeader(http.StatusBadGateway)
				return
			}

			defer resp.Body.Close()
			resp_body, _ := ioutil.ReadAll(resp.Body)

			log.Println(resp.Status)
			log.Println(string(resp_body))

			// w.Write(resp_body)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	})

	
	log.Fatal("HTTP server error: ", http.ListenAndServe(":" + port, nil))
}
