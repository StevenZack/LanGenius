package LanGenius

import (
	"fmt"
	"net/http"
)

var (
	flag_staticSiteIsRunning bool
)

func StartStaticSite(port, dir string) {
	go func() {
		e := http.ListenAndServe(port, http.StripPrefix("/", http.FileServer(http.Dir(dir))))
		if e != nil {
			fmt.Println(e)
		}
	}()
	flag_staticSiteIsRunning = true
}
func IsStaticSiteRunning() bool {
	return flag_staticSiteIsRunning
}
