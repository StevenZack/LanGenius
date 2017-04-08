package LanGenius

import (
	"context"
	"errors"
	f "fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	mux             map[string]func(http.ResponseWriter, *http.Request)
	server          http.Server
	javahandler     JavaHandler
	tdata           TData
	str_storagePath string
)

func init() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	tdata = TData{Clipboard: "Clipboard:", Copy: "copy", Send: "send", Files: "Files", CbContent: "", UploadButton: "upload"}
	str_storagePath = "/sdcard/"
}

// func main() {
// 	var jh = JavaClass{}
// 	str_storagePath = "/home/steven/Documents/"
// 	Start("zh", jh)
// 	f.Scanf("%s")
// 	Stop()
// }

type JavaHandler interface {
	OnClipboardReceived(string)
	OnFileReceived(string)
}
type MyFileEntry struct {
	FileName string
	FilePath string
}
type TData struct {
	Clipboard, Copy, Send, CbContent, Files string
	FileSlice                               []MyFileEntry
	UploadButton                            string
}
type MyHandler struct{}

func (*MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.RequestURI)
	if h, ok := mux[u.Path]; ok {
		h(w, r)
		return
	}
	http.NotFound(w, r)
}

type JavaClass struct{}

func (JavaClass) OnClipboardReceived(str string) {
	f.Println("ClipboardReceived():" + str)
}
func (JavaClass) OnFileReceived(str string) {
	f.Println("OnFileReceived():" + str)
}
func Start(language string, jh JavaHandler) {
	if language == "zh" {
		tdata.Clipboard = "复制内容"
		tdata.Copy = "复制"
		tdata.Send = "发送"
		tdata.Files = "共享的文件"
		tdata.UploadButton = "上传文件"
	}
	javahandler = jh
	mux["/"] = home
	mux["/send"] = send
	mux["/downloadFile"] = downloadFile
	mux["/uploadFile"] = uploadFile

	go func() {
		server = http.Server{Addr: ":4444", Handler: &MyHandler{}}
		err := server.ListenAndServe()
		if err != nil {
			f.Println(err)
		}
	}()
}
func Stop() {
	err := server.Shutdown(context.Background())
	if err != nil {
		f.Println(err)
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	t := template.New("homeTPL")
	t.Parse(`<!DOCTYPE html>
<html>
<head>
	<title>shareMe</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width maximum-scale=1 initial-scale=1">
	<script type="text/javascript">
		function copy() {
			var cb=document.getElementById("cb")
			cb.select()
			document.execCommand("Copy")
			document.getElementById("spanInfo").innerHTML="OK"
		}
		function send() {
			var str=document.getElementById("cb").value
			location.href="send?url="+encodeURIComponent(str)
		}
		function onFileSelected() {
			document.getElementById('submit_button').disabled=false
		}
	</script>
</head>
<body>
<center>
<table>
	<tr>
		<th style="color: #D81B60">{{.Clipboard}}</th>
	</tr>
	<tr>
		<td>
			<textarea name="cb" id="cb" cols="30" rows="5">{{.CbContent}}</textarea>
		</td>
		<td>
			<input type="button" value="{{.Copy}}" onclick="copy()"><br>
			<br>	
			<input type="button" value="{{.Send}}" onclick="send()">
		</td>
		<td><span id="spanInfo"></span><br><span></span></td>
	</tr>
	<tr>
	<td colspan="2"><hr></td>
	</tr>
	<tr>
		<th style="color: #1E88E5">{{.Files}}</th>
	</tr>
	<tr>
		<td>
		{{range .FileSlice}}
			<a href="/downloadFile?filename={{.FileName}}">
			{{.FileName}}</a><br>
		{{end}}
		</td>
	</tr>
	<tr>
	<td colspan="2"><hr></td>
	</tr>
	<tr>
		<form action="/uploadFile" method="post" enctype="multipart/form-data">
			<td><input type="file" name="myUploadFile" onchange="onFileSelected()"></td>
			<td><input type="submit" id="submit_button" disabled value="{{.UploadButton}}"></td>
		</form>
	</tr>
</table>
</center>
</body>
</html>`)
	t.Execute(w, tdata)
}
func AddFile(str string) error {
	fns := strings.Split(str, "/")
	if len(fns) < 1 {
		return errors.New("Bad file path")
	}
	tdata.FileSlice = append(tdata.FileSlice, MyFileEntry{FilePath: str, FileName: fns[len(fns)-1]})
	return nil
}
func SetStoragePath(str string) {
	str_storagePath = str
}
func SetClipboard(str string) {
	tdata.CbContent = str
	f.Println("CbContent has been set to:" + str)
}
func send(w http.ResponseWriter, r *http.Request) {
	f.Println(r.FormValue("url"))
	str, _ := url.QueryUnescape(r.FormValue("url"))
	javahandler.OnClipboardReceived(str)
	SetClipboard(str)
	http.Redirect(w, r, "/", http.StatusFound)
}
func downloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	for _, v := range tdata.FileSlice {
		if v.FileName == filename {
			w.Header().Add("Content-Disposition", f.Sprintf("attachment; filename=%s", filename))
			http.ServeFile(w, r, v.FilePath)
			return
		}
	}
	http.NotFound(w, r)
}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("myUploadFile")
	if err != nil {
		f.Println(err)
		f.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><meta charset="utf-8">
		<meta name="viewport" content="width=device-width initial-scale=1 maximum-scale=1"
		</head>
		<body>
		<center>
		<a href="/"><h4>`+err.Error()+`</h4></a>
		</center>
		</body>
		</html>
		`)
		return
	}
	javahandler.OnFileReceived(handler.Filename)
	defer file.Close()
	myfile, err := os.OpenFile(str_storagePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		f.Println(err)
		f.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><meta charset="utf-8">
		<meta name="viewport" content="width=device-width initial-scale=1 maximum-scale=1"
		</head>
		<body>
		<center>
		<a href="/"><h4>`+err.Error()+`</h4></a>
		</center>
		</body>
		</html>
		`)
		return
	}
	defer myfile.Close()
	io.Copy(myfile, file)
	f.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><meta charset="utf-8">
		<meta name="viewport" content="width=device-width initial-scale=1 maximum-scale=1"
		</head>
		<body>
		<center>
		<a href="/"><h1>OK</h1></a>
		</center>
		<script type="text/javascript">
		 setTimeout("window.history.back(-1)",2000)
		</script>
		</body>
		</html>
		`)
}
