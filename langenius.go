package LanGenius

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

func Start(eh EventHandler, port, tmpDir, sPath string) {
	mEventHandler = eh
	storagePath = sPath
	if tmpDir != "" {
		os.Setenv("TMPDIR", tmpDir)
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/send", send)
	http.HandleFunc("/download/", download)
	http.HandleFunc("/viewfile/", viewfile)
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
func contains(s, p string) bool {
	for _, v := range s {
		if string(v) == p {
			return true
		}
	}
	return false
}
func getFileName(s string) string {
	for i := len(s) - 2; i > -1; i-- {
		if s[i:i+1] == "/" {
			return s[i+1:]
		}
	}
	return ""
}
