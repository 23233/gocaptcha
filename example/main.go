package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/23233/gocaptcha"
)

const (
	dx = 180
	dy = 60
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/get/", Get)
	fmt.Println("服务已启动 -> http://127.0.0.1:8800")
	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = t.Execute(w, nil)
}
func Get(w http.ResponseWriter, r *http.Request) {

	_, bt, err := gocaptcha.GenerateCaptcha(dx, dy, 4, gocaptcha.CaptchaVeryEasy)

	if err != nil {
		fmt.Println(err)
	}

	_, _ = w.Write(bt)

}
