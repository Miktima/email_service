package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

type appContext struct {
	ctx *context.Context
}

func main() {
	mux := http.NewServeMux()

	var ctx *appContext

	mux.HandleFunc("/login", LoginHandler(ctx))
	http.ListenAndServe(":8090", mux)
}

func LoginHandler(h *appContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Проверяем метод запроса - для логина принимаем только POST
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		login := req.PostFormValue("login")
		password := req.PostFormValue("password")

		// create context
		ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
		defer cancel()
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
		// run task list
		var res string
		ok := true
		err := chromedp.Run(ctx,
			chromedp.Navigate(`https://mail.rambler.ru/`),
			chromedp.WaitVisible(`#login`),
			chromedp.SendKeys(`#login`, login),
			chromedp.SendKeys(`#password`, password),
			chromedp.Click(`//form/button[@type="submit"]`, chromedp.BySearch),
			chromedp.WaitVisible(`//footer/div/button`, chromedp.BySearch),
			chromedp.Click(`//footer/div/button`, chromedp.BySearch),
			chromedp.WaitVisible(`.FolderItem-root-1t`, chromedp.ByQuery),
			chromedp.AttributeValue(".FolderItem-root-1t > a", "title", &res, &ok, chromedp.ByQuery),
		)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)

		if len(res) > 0 {
			resp["status"] = "Ok"
			resp["message"] = res
		} else {
			resp["status"] = "Error"
			resp["message"] = res
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}
}
