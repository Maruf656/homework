package router

import (
	"encoding/json"
	"errors"
	"github.com/Abdulhalim92/server/internal/models"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func StartRouter() error {

	// вывести в отдельные функции
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Println(err)
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}

	var config models.Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Println(err)
		return err
	}

	address := net.JoinHostPort(config.Host, config.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/history", GetHistory)
	mux.HandleFunc("/calculate", Calculate)

	err = http.ListenAndServe(address, mux)
	if err != nil {
		log.Println("listen and serve", err)
		return err
	}

	return nil
}

func Calculate(w http.ResponseWriter, r *http.Request) {
	// Получение query - параметров
	queries := r.URL.Query()
	// Получение первого числа
	numOne := queries.Get("num_one")
	firstNum, err := strconv.ParseFloat(numOne, 64)
	if err != nil {
		return
	}
	// Получение второго числа
	numTwo := queries.Get("num_two")
	secondNum, err := strconv.ParseFloat(numTwo, 64)
	if err != nil {
		return
	}
	// Получение арифметической операции
	operation := queries.Get("operation")

	// Нахождение результата операции между двумя числами
	var result float64

	switch operation {
	case "+":
		result = firstNum + secondNum
	case "-":
		result = firstNum - secondNum
	case "*":
		result = firstNum * secondNum
	case "/":
		result = firstNum / secondNum
	default:
		log.Println("not an operation")
	}
	// Запись данных в структуру
	HistoryElement := models.HistoryElement{
		NumberOne: firstNum,
		NumberTwo: secondNum,
		Operation: operation,
		Result:    result,
	}
	// Провека
	jsonFile, err := os.OpenFile("./history.json", os.O_RDWR, 0777)
	if err != nil {
		log.Println("didn't do")
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {
			log.Println("the file didn't closed")
		}
	}(jsonFile)
	buf := make([]byte, 4096)
	n, err := jsonFile.Read(buf)
	if err != nil {
		log.Println("didn't read")
	}
	buf = buf[:n]
	var HistoryEl models.HistoryElement
	err = json.Unmarshal(buf, &HistoryEl)
	if err != nil {
		log.Println("didn't do")
	}
	log.Println(string(buf))
	log.Println(HistoryEl)
	// Сохранение истории операций над элементами
	var HistoryElements []models.HistoryElement
	HistoryElements = append(HistoryElements, HistoryElement)
	// Сериализация структуры HistoryElement в JSON
	bytes, err := json.MarshalIndent(&HistoryElements, "", "    ")
	if err != nil {
		return
	}
	// Запись сериализованного JSON в файл JSON
	err = os.WriteFile("./history.json", bytes, 0777)
	if err != nil {
		return
	}
}

func GetHistory(w http.ResponseWriter, r *http.Request) {

	err := errors.New("fake error")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	// получение истории калькулятора
	//var History []*models.HistoryElement // response ваш ответ
}
