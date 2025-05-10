package server

import (
	"go1f/pkg/api"
	"net/http"
)

const webDir = "./web"

func StartServer(port string) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}
	return nil
}
