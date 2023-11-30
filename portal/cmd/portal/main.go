package main

import (
	"net/http"

	"github.com/rosylilly/performania/portal/www"
)

func main() {
	// _, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	http.ListenAndServe(":8080", http.HandlerFunc(www.ServeStaticHandler))
}
