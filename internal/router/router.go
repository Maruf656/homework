package router

import (
	"encoding/json"
	"github.com/Abdulhalim92/server/config"
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

	getConfig, err := config.GetConfig()
	if err != nil {
		log.Println("Не получилось получить настройки")
		return err
	}

	address := net.JoinHostPort(getConfig.Host, getConfig.Port)

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
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Получение второго числа
	numTwo := queries.Get("num_two")
	secondNum, err := strconv.ParseFloat(numTwo, 64)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Получение арифметической операции
	operation := queries.Get("operation")

	// Нахождение результата операции между двумя числами
	result := Calc(firstNum, secondNum, operation)
	// Запись данных в структуру
	HistoryElement := models.HistoryElement{
		NumberOne: firstNum,
		NumberTwo: secondNum,
		Operation: operation,
		Result:    result,
	}
	// Сохранение истории операций над элементами
	// Открываем и читаем файл JSON
	contentJSON, err := FileOpen()
	if err != nil {
		log.Println("Не получилось получить контент")
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Добавляем в срез объектов новую историю результатов
	var HistoryEL []models.HistoryElement
	err = json.Unmarshal(contentJSON, &HistoryEL)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("didn't unmarshal")
		return
	}
	lengthHistory := len(HistoryEL) // длина среза объектов
	log.Println(lengthHistory)
	HistoryEL = append(HistoryEL, HistoryElement)
	// Сериализация структуры HistoryElement в JSON
	bytes, err := json.MarshalIndent(HistoryEL, "", "    ")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Запись сериализованного JSON в файл JSON
	err = os.WriteFile("./history.json", bytes, 0777)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
func Calc(firstNum float64, secondNum float64, operation string) float64 {
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
	return result
}
func FileOpen() ([]byte, error) {
	// Открываем файл
	file, err := os.OpenFile("./history.json", os.O_RDWR, 0777)
	if err != nil {
		log.Println("didn't open")
		return nil, err
	}
	// Читаем файл
	contentJSON, err := io.ReadAll(file)
	if err != nil {
		log.Println("didn't read")
		return nil, err
	}

	return contentJSON, nil

}
func GetHistory(w http.ResponseWriter, r *http.Request) {

	// получение истории калькулятора
	var History []models.HistoryElement // response ваш ответ
	// Открываем и читаем файл JSON
	contentJSON, err := FileOpen()
	if err != nil {
		log.Println("Не получилось получить контент")
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// log.Println(string(contentJSON))
	err = json.Unmarshal(contentJSON, &History)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	log.Println(History)
	_, err = w.Write(contentJSON)
	if err != nil {
		log.Println("Не получилось")
		return
	}
}
