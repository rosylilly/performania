//go:build release

package www

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed dist/*
	distFiles embed.FS

	distFileServer http.Handler
)

func init() {
	rootFS, err := fs.Sub(distFiles, "dist")
	if err != nil {
		panic(err)
	}

	distFileServer = http.FileServer(http.FS(rootFS))
}

func ServeStaticHandler(w http.ResponseWriter, r *http.Request) {
	distFileServer.ServeHTTP(w, r)
}
