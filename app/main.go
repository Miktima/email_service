package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	mux := http.NewServeMux()

	// create context
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf), chromedp.WithDebugf(log.Printf))
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	mux.HandleFunc("/login", LoginHandler(&ctx))
	mux.HandleFunc("/logout", LogoutHandler(&ctx))
	mux.HandleFunc("/mail", MailHandler(&ctx))
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

		login := req.PostFormValue("login")
		password := req.PostFormValue("password")

		// run task list
		var res string
		ok := true
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")

		// Проходим авторизацию в почту и возвращаем текст с количеством писем
		err := chromedp.Run(*ctx,
			chromedp.Navigate(`https://mail.rambler.ru/`),
			chromedp.WaitVisible(`#login`, chromedp.ByQuery),
			chromedp.SendKeys(`#login`, login, chromedp.ByQuery),
			chromedp.SendKeys(`#password`, password, chromedp.ByQuery),
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
			chromedp.Sleep(3*time.Second),
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
			ctx = nil
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

func MailHandler(ctx *context.Context) func(http.ResponseWriter, *http.Request) {
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

		var title, body string
		err := chromedp.Run(*ctx,
			chromedp.WaitVisible(`//div[@class="MailList-list-2L"]/div[@draggable="true"]`, chromedp.BySearch),
			chromedp.Click(`//div[@class="MailList-list-2L"]/div[@draggable="true"]`, chromedp.BySearch),
			chromedp.WaitVisible(`.ThemeBorder-root-26`, chromedp.ByQuery),
			chromedp.Text(`div.ThemeBorder-root-26 > div.LetterHeader-root-S0 > div.LetterHeader-title-1D`, &title, chromedp.ByQuery),
			chromedp.Text(`div.ThemeBorder-root-26 > div.LetterBody-root-3k`, &body, chromedp.ByQuery),
			chromedp.Click(`//div[@class="HeaderToolbar-toolbar-15"]/div[@class="ToolbarButton-root-1B"]/div`, chromedp.BySearch),
			chromedp.WaitVisible(`.MailList-list-2L`, chromedp.ByQuery),
		)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		if len(title) > 0 {
			resp["status"] = "Ok"
			resp["title"] = title
			resp["body"] = body
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
