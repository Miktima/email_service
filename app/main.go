package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func main() {
	mux := http.NewServeMux()

	// create context
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	mux.HandleFunc("/login", LoginHandler(&ctx))
	mux.HandleFunc("/logout", LogoutHandler(&ctx))
	http.ListenAndServe(":8090", mux)
}

func LoginHandler(ctx *context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := make(map[string]string)
		// Проверяем метод запроса - для логина принимаем только POST
		if req.Method != http.MethodPost {
			resp["status"] = "Error"
			resp["message"] = "Method not allowed"
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}

		log.Println("ctx login1:", *ctx)
		login := req.PostFormValue("login")
		password := req.PostFormValue("password")

		// Определяем user agent
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
		// run task list
		var res string
		ok := true
		// Проходим авторизацию в почту и возвращаем текст с количеством писем
		err := chromedp.Run(*ctx,
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
		// формируем статус ответа
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		if len(res) > 0 {
			resp["status"] = "Ok"
			resp["message"] = res
		} else {
			resp["status"] = "Error"
			resp["message"] = "Response is empty"
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}
}

func LogoutHandler(ctx *context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := make(map[string]string)

		if ctx == nil {
			resp["status"] = "Error"
			resp["message"] = "Browser context is empty"
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
		var res string
		err := chromedp.Run(*ctx,
			chromedp.WaitVisible(`//div[@class="rc__mAFe4"]/div[@class="rc__xgBcB"]/button`, chromedp.BySearch),
			chromedp.Click(`//div[@class="rc__mAFe4"]/div[@class="rc__xgBcB"]/button`, chromedp.BySearch),
			chromedp.WaitVisible(`//div[@class="rc__IP0ui"]/div[4]/button`, chromedp.BySearch),
			chromedp.Click(`//div[@class="rc__IP0ui"]/div[4]/button`, chromedp.BySearch),
			chromedp.WaitVisible(`.ad_branding_site`, chromedp.ByQuery),
			chromedp.Location(&res),
		)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		if len(res) > 0 {
			resp["status"] = "Ok"
			resp["message"] = res
		} else {
			resp["status"] = "Error"
			resp["message"] = "Response is empty"
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)

	}
}
