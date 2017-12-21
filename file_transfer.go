package LanGenius

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
)

type EventHandler interface {
	OnClipboardReceived(string)
	OnFileReceived(string)
}

var (
	mEventHandler EventHandler
	storagePath   string
)

type FileEntry struct {
	Name, Path string
}

var (
	homeData struct {
		Files            []FileEntry
		Clipboard        string
		ClipboardEnabled bool
	}
)

func home(w http.ResponseWriter, r *http.Request) {
	// t := template.New("homeTPL")
	// t.Parse(`
	// 	`)
	// t.Execute(w, homeData)
	t, e := template.ParseFiles("/home/asd/go/src/LanGenius/views/index.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	t.Execute(w, homeData)
}
func send(w http.ResponseWriter, r *http.Request) {
	var gobj struct {
		Cb string
	}
	e := json.Unmarshal([]byte(r.FormValue("cb")), &gobj)
	if e != nil {
		fmt.Fprint(w, e.Error())
		return
	}
	mEventHandler.OnClipboardReceived(gobj.Cb)
	fmt.Fprint(w, "OK")
}
func download(w http.ResponseWriter, r *http.Request) {
	filename, _ := url.QueryUnescape(r.URL.RequestURI()[len("/download/"):])
	Logd(filename)
	for _, v := range homeData.Files {
		if v.Name == filename {
			w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			http.ServeFile(w, r, v.Path)
			return
		}
	}
	http.NotFound(w, r)
}
func viewfile(w http.ResponseWriter, r *http.Request) {
	filename, _ := url.QueryUnescape(r.URL.RequestURI()[len("/viewfile/"):])
	Logd(filename)
	for _, v := range homeData.Files {
		if v.Name == filename {
			http.ServeFile(w, r, v.Path)
			return
		}
	}
	http.NotFound(w, r)
}
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fhs := r.MultipartForm.File["myUploadFile"]
	for _, v := range fhs {
		file, e := v.Open()
		if e != nil {
			fmt.Println(e)
			fmt.Fprint(w, e.Error())
			return
		}
		mEventHandler.OnFileReceived(v.Filename)
		mf, e := os.OpenFile(storagePath+v.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if e != nil {
			fmt.Println(e)
			fmt.Fprint(w, e.Error())
			return
		}
		defer mf.Close()
		io.Copy(mf, file)
	}
	fmt.Fprint(w, "OK")
}
func AddFile(str string) {
	if !contains(str, "/") {
		return
	}
	homeData.Files = append(homeData.Files, FileEntry{Name: getFileName(str), Path: str})
}
func RemoveFile(index int) {
	if index > -1 && index < len(homeData.Files) {
		homeData.Files = append(homeData.Files[:index], homeData.Files[index+1:]...)
	} else {
		fmt.Println("remove file index out of bound: index=", index)
	}
}
func SetStoragePath(str string) {
	storagePath = str
}
func SetClipboardEnabled(b bool) {
	homeData.ClipboardEnabled = b
}
func SetClipboard(str string) {
	homeData.Clipboard = str
}
func Logd(e interface{}) {
	fmt.Println(e)
}
