package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Human struct {
	Name       string    `json:"name"`
	SecondName string    `json:"second_name"`
	Date       time.Time `json:"date"`
}
type SecondHuman struct {
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
	Date       string `json:"date"`
}

func main() {

	ip := "localhost"
	port := ":9999"

	// Rout --- что это такое
	http.HandleFunc("/hello", GetQuery)

	address := ip + port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Println("listen and serve", err)
	}
}

func GetQuery(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("URI:", request.RequestURI)
	//fmt.Println("body:", request.Body)
	//fmt.Println("length:", request.ContentLength)
	//fmt.Println("header:", request.Header)
	//fmt.Println("method:", request.Method)
	//fmt.Println("Host:", request.Host)
	//fmt.Println("URL:", request.URL)
	x := r.URL.Query()

	timeNow := time.Now()
	WriteToJson(x.Get("name"), x.Get("secondName"), timeNow)
	// Второе условие
	secondTime := time.Now().Format("15:04:05 02-01-06")
	var human SecondHuman = SecondHuman{
		Name:       x.Get("name"),
		SecondName: x.Get("secondName"),
		Date:       secondTime,
	}
	content, err := json.MarshalIndent(human, "", "    ")
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
	}

}

func WriteToJson(name string, secondName string, time time.Time) {
	var human Human = Human{
		Name:       name,
		SecondName: secondName,
		Date:       time,
	}
	buf, err := json.MarshalIndent(human, "", "    ")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile("./student.json", buf, 0777)
	if err != nil {
		log.Println("did not written")
	}
}

// Порт 80 и 480
// Программи postman, insomnia
// Статус коды
