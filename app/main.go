// Command multi is a chromedp example demonstrating how to use headless-shell
// and a container (Docker, Podman, other). See README.md.
package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {

	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)

}

// func run(ctx context.Context, verbose bool, wait time.Duration, out string, urls ...string) error {
func hello(w http.ResponseWriter, req *http.Request) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://ifconfig.me/`),
		chromedp.Text(`#ip_address`, &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(strings.TrimSpace(res))
}
