package LanGenius

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
	"os"
)

type EventHandler interface {
	OnClipboardReceived(string)
	OnFileReceived(string)
}
type FileEntry struct {
	Name, Path string
}

var (
	mEventHandler EventHandler
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

var (
	homeData struct {
		Files     []FileEntry
		Clipboard string
	}
)

func home(w http.ResponseWriter, r *http.Request) {
	t := template.New("homeTPL")
	t.Parse(``)
	t.Execute(w, homeData)
}
