package LanGenius

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

func Start(eh EventHandler, port, pkg string) {
	mEventHandler = eh
	os.Setenv("TMPDIR", pkg)
	http.HandleFunc("/", home)
	http.HandleFunc("/send", send)
	http.HandleFunc("/downloadFile/", downloadFile)
	http.HandleFunc("/upload", upload)

	//live part
	http.Handle("/live/chat", websocket.Handler(chat))
	http.HandleFunc("/live", live)
	http.HandleFunc("/live/camera", camera)

	go func() {
		e := http.ListenAndServe(port, nil)
		if e != nil {
			fmt.Println(e)
		}
	}()
}
