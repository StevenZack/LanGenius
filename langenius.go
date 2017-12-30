package LanGenius

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"os"
)

var (
	mEventHandler      EventHandler
	storagePath, mPort string
)

type EventHandler interface {
	OnClipboardReceived(string, bool)
	OnFileReceived(string)
	OnFolderReceived(string)
	OnDeviceOnline(string)
	OnDeviceOffline(string)
	OnRemoteControlCmdReceived(string)
	OnMessageReceived(string)
}
type Msg struct {
	State, Type, IP, Port, Info string
	Message                     string
	RemoteControlCmd            string
	RemoteControlStatus         bool
}

func Start(eh EventHandler, port, tmpDir, sPath string) {
	go Run(eh, port, tmpDir, sPath)
}
func Run(eh EventHandler, port, tmpDir, sPath string) {
	mEventHandler = eh
	if sPath[len(sPath)-1:] != "/" {
		storagePath = sPath + "/"
	} else {
		storagePath = sPath
	}
	mPort = port
	if tmpDir != "" {
		if tmpDir[len(tmpDir)-1:] != "/" {
			os.Setenv("TMPDIR", tmpDir+"/")
		} else {
			os.Setenv("TMPDIR", tmpDir)
		}
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
	http.HandleFunc("/camera", camera)

	e := http.ListenAndServe(port, nil)
	if e != nil {
		fmt.Println(e)
	}
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
		if len(v) > 8 && v[:8] == "192.168." {
			return v
		}
	}
	for _, v := range strs {
		if len(v) > 3 && v[:3] == "10." {
			return v
		}
	}
	for _, v := range strs {
		if len(v) > 4 && v[:4] == "172." {
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

func IsMyIP(str string) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.String() == str {
					return true
				}
			}
		}
	}
	return false
}
