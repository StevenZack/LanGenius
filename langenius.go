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
	http.Handle("/wsClipboard", websocket.Handler(wsClipboard))
	http.HandleFunc("/download/", download)
	http.HandleFunc("/viewfile/", viewfile)
	http.HandleFunc("/upload", upload)

	//live part
	http.Handle("/wsLive", websocket.Handler(wsLive))
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
func GetIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var strs []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				strs = append(strs, ip.String())
			case *net.IPAddr:
				// ip := v.IP
				// strs = append(strs, ip.String())
			}
		}
	}
	for _, v := range strs {
		if v[:3] == "10." {
			return v
		}
	}
	for _, v := range strs {
		if v[:4] == "172." {
			return v
		}
	}
	for _, v := range strs {
		if v[:8] == "192.168." {
			return v
		}
	}
	for _, v := range strs {
		if v != "127.0.0.1" && v != "::1" {
			return v
		}
	}
	return strs[0]
}
