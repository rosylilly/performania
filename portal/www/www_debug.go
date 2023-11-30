//go:build !release

package www

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var viteDevServer *httputil.ReverseProxy

func init() {
	viteDevServerURL := "http://localhost:5173"
	if u, ok := os.LookupEnv("VITE_DEV_SERVER_URL"); ok {
		viteDevServerURL = u
	}
	u, err := url.Parse(viteDevServerURL)
	if err != nil {
		panic(err)
	}

	viteDevServer = httputil.NewSingleHostReverseProxy(u)
}

func ServeStaticHandler(w http.ResponseWriter, r *http.Request) {
	viteDevServer.ServeHTTP(w, r)
}
