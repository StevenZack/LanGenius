package LanGenius

import (
	"context"
	"errors"
	f "fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	tdata = TData{Clipboard: "Clipboard:", Copy: "copy", Send: "send", Files: "Files", CbContent: "", UploadButton: "upload", KC_enabled: false}
	str_storagePath = "/sdcard/"
}

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
	KC_enabled                              bool
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

type javaClass struct{} //just for test

func (javaClass) OnClipboardReceived(str string) {
	f.Println("ClipboardReceived():" + str)
}
func (javaClass) OnFileReceived(str string) {
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
	mux["/downloadKC"] = downloadKC
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
		function DoUpload(){
			if (document.getElementById("myUploadFile").value != "") {
				document.getElementById("submit_button").disabled=true;
	            var fileObj = document.getElementById("myUploadFile").files[0];
	            //创建xhr
	            var xhr = new XMLHttpRequest();
	            var url = "uploadFile";
	            //FormData对象
	            var fd = new FormData();
	            fd.append("myUploadFile", fileObj); 
	            fd.append("acttime",new Date().toString());    //本人喜欢在参数中添加时间戳，防止缓存（--、）
	            xhr.onreadystatechange = function () {
	                if (xhr.readyState == 4 && xhr.status == 200) {
	                    var result = xhr.responseText;
	                    document.getElementById("result").innerHTML = result;
	                }
	            }
	            //进度条部分
	            xhr.upload.onprogress = function (evt) {
	                if (evt.lengthComputable) {
	                    var percentComplete = Math.round(evt.loaded * 100 / evt.total);
	                    document.getElementById('uploadProgress').value = percentComplete;
	                }
	            };
	            xhr.open("POST", url, true);
	            xhr.send(fd);
	        }else{
				document.getElementById("submit_button").disabled=false;
	        }
		}
		function downloadKC(){
			var myos=detectOS()
			if (navigator.language=="zh-CN") {
				if (confirm("是否下载受控端 for "+myos+" ?")) {
					window.location.href="downloadKC?os="+myos
				}
				if (myos=="Windows") {
					if (confirm("是否下载受控端 for Windows x64 ?")) {
					window.location.href="downloadKC?os="+myos
					}
				}else if (myos=="Linux") {
					if (confirm("是否下载受控端 for Linux x64 ?")) {
					window.location.href="downloadKC?os="+myos
					}
				}else {
					if (confirm("是否下载受控端 for Windows x64 ? 暂时不支持你的操作系统")) {
					window.location.href="downloadKC?os="+myos
					}
				}
			}else {
				if (myos=="Windows") {
					if (confirm("Download Controlled End executable for Windows x64 ?")) {
					window.location.href="downloadKC?os="+myos
					}
				}else if (myos=="Linux") {
					if (confirm("Download Controlled End executable for Linux x64 ?")) {
					window.location.href="downloadKC?os="+myos
					}
				}else {
					if (confirm("Download Controlled End executable for Windows x64 ? We don't support Mac yet")) {
					window.location.href="downloadKC?os="+myos
					}
				}
			}
		}
        function detectOS() {
            var sUserAgent = navigator.userAgent;
            var isWin = (navigator.platform == "Win32") || (navigator.platform == "Windows");
            if (isWin) return "Windows"
            var isMac = (navigator.platform == "Mac68K") || (navigator.platform == "MacPPC") || (navigator.platform == "Macintosh") || (navigator.platform == "MacIntel");
            if (isMac) return "Mac";
            var isLinux = (String(navigator.platform).indexOf("Linux") > -1)||(String(navigator.platform).indexOf("Android") > -1)||((navigator.platform == "X11") && !isWin && !isMac);
            if (isLinux) return "Linux";
            return "other";
        }
        function detectLang(){
            var language_zh_cn = "zh-CN";  
            var currentLang = navigator.language;  
        }
	</script>
	<style type="text/css">
		.Mybutton{
			border: 5px solid #0088ff;
			border-radius: 10px;
			height: 30px;
			line-height: 30px;
			
			cursor: pointer;
		}
	</style>
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
		<td><input type="file" name="myUploadFile" id="myUploadFile" onchange="onFileSelected()"></td>
		<td><input type="button" id="submit_button" disabled value="{{.UploadButton}}" onclick="DoUpload()"></td>
	</tr>
	<tr><td colspan="2">
	    <progress value="0" max="100" id="uploadProgress" style="height: 4px; width: 100%"></progress>
	</td></tr>
	<tr><td colspan="2">
	    <div id="result"></div>
	</td></tr>
	<tr><td><br><br><br><br></td></tr>
	{{if .KC_enabled}}
	<tr><td colspan="2">
		<div class="Mybutton" id="kcbt" align="center" onmouseover="this.setAttribute('style','box-shadow: 2px 2px 15px #0088ff')" onmouseout="this.setAttribute('style','')" onclick="downloadKC()">{{.KCenabled}}</div>
	</td></tr>
	{{end}}
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
	r.ParseMultipartForm(1 << 30)
	file, handler, err := r.FormFile("myUploadFile")
	if err != nil {
		f.Println(err)
		f.Fprintf(w, err.Error())
		return
	}
	javahandler.OnFileReceived(handler.Filename)
	defer file.Close()
	myfile, err := os.OpenFile(str_storagePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		f.Println(err)
		f.Fprintf(w, err.Error())
		return
	}
	defer myfile.Close()
	io.Copy(myfile, file)
	f.Fprintf(w, `OK`)
}
func downloadKC(w http.ResponseWriter, r *http.Request) {
	var myos = r.FormValue("os")
	if myos == "Linux" {
		filename := "KC_Linux"
		w.Header().Add("Content-Disposition", f.Sprintf("attachment; filename=%s", filename))
		http.ServeFile(w, r, "/data/data/com.xchat.stevenzack.langenius/"+filename)
	} else {
		filename := "KC_Windows.exe"
		w.Header().Add("Content-Disposition", f.Sprintf("attachment; filename=%s", filename))
		http.ServeFile(w, r, "/data/data/com.xchat.stevenzack.langenius/"+filename)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

//keyboard controller
var (
	kc_addr, kc_adbr *net.UDPAddr
	kc_adds          []*net.UDPAddr
	c                *net.UDPConn
	kchandler        JavaKCHandler
)

type JavaKCHandler interface {
	OnDeviceDetected(string)
}

func StartKC(jh JavaKCHandler) {
	kchandler = jh
	tdata.KC_enabled = true
	kc_addr, _ = net.ResolveUDPAddr("udp", ":9943")
	kc_adbr, _ = net.ResolveUDPAddr("udp", "255.255.255.255:9942")
	c, err := net.ListenUDP("udp", kc_addr)
	if err != nil {
		f.Println(err)
		return
	}
	go sendKCPulse()
	go readKC(c)
}
func StopKC() {
	tdata.KC_enabled = false
	c.Close()
}
func SendKC(cmd string, index int) error {
	if len(kc_adds) > 0 && index > -1 && index < len(kc_adds) {
		c.WriteToUDP([]byte(cmd), kc_adds[index])
		return nil
	} else {
		return errors.New("bad index")
	}
}
func readKC(c *net.UDPConn) {
	b := make([]byte, 512)
	for {
		n, ra, err := c.ReadFromUDP(b)
		if err != nil {
			f.Println(err)
			return
		}
		if string(b[:n]) == "LanGenius-from-desktop" {
			kc_adds = append(kc_adds, ra)
			kchandler.OnDeviceDetected(ra.String())
		}
	}
}
func sendKCPulse() {
	for {
		c.WriteToUDP([]byte("LanGenius-from-android"), kc_adbr)
		time.Sleep(time.Second * 3)
	}
}
