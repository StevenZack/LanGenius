package LanGenius

import (
	// "context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	// "errors"
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"math/big"
	"net"
	"net/http"
	// "net/url"
	"os"
	"time"
)

var (
	mux                      map[string]func(http.ResponseWriter, *http.Request)
	server                   http.Server
	javahandler              JavaHandler
	tdata                    TData
	str_storagePath          string
	flag_staticSiteIsRunning bool = false
	web_sync_cons            map[string]*websocket.Conn
)

func init() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	web_sync_cons = make(map[string]*websocket.Conn)
	tdata = TData{Clipboard: "Clipboard:", Copy: "copy", Send: "send", Files: "Files", CbContent: "", UploadButton: "upload", KC_enabled: false, CbEnabled: false}
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
	KC_enabled, CbEnabled                   bool
}

func Start(language string, jh JavaHandler, port string, showtimeFolder string) {
	ShowtimeFolder = showtimeFolder
	if language == "zh" {
		tdata.Clipboard = "剪切板内容"
		tdata.Copy = "复制"
		tdata.Send = "发送"
		tdata.Files = "共享的文件"
		tdata.UploadButton = "上传文件"
	}
	javahandler = jh
	http.HandleFunc("/", home)
	http.HandleFunc("/send", send)
	http.HandleFunc("/downloadFile", downloadFile)
	http.HandleFunc("/downloadKC", downloadKC)

	//showtime part
	http.Handle("/showtime/chat", websocket.Handler(pwint))
	http.Handle("/showtime/public/", http.StripPrefix("/showtime/public/", http.FileServer(http.Dir(ShowtimeFolder+"/public"))))
	http.HandleFunc("/showtime/live", live)
	http.HandleFunc("/showtime/camera", camera)
	http.Handle("/showtime", http.RedirectHandler("/showtime/live", 301))

	//ws sync part

	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

type WsMessage struct {
	State, Info string
}

func web_sync(ws *websocket.Conn) {
	defer ws.Close()
	web_sync_cons[ws.RemoteAddr().String()] = ws
	defer func() {
		delete(web_sync_cons, ws.RemoteAddr().String())
	}()
	for {
		msg := ""
		if e := websocket.Message.Receive(ws, &msg); e != nil {
			break
		}
		wm := WsMessage{}
		e := json.Unmarshal([]byte(msg), &wm)
		if e != nil {

		}
	}
}

func StartStaticSite(port, dir string) {
	go func() {
		err := http.ListenAndServe(port, http.StripPrefix("/", http.FileServer(http.Dir(dir))))
		if err != nil {
			fmt.Println(err)
		}
	}()
	flag_staticSiteIsRunning = true
}
func IsStaticSiteRunning() bool {
	return flag_staticSiteIsRunning
}
func home(w http.ResponseWriter, r *http.Request) {
	t := template.New("homeTPL")
	t.Parse(`<!DOCTYPE html>
<html>
<head>
	<title>LanGenius</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width maximum-scale=1 initial-scale=1">
	<script type="text/javascript">
		function copy(me) {
			var cb=document.getElementById("cb")
			cb.select()
			document.execCommand("Copy")
			var myvalue=me.value
			me.value="√"
			setTimeout("document.getElementById('btCopy').value='"+myvalue+"'",1000)
		}
		function send(me) {
			var myvalue=me.value
			var str=document.getElementById("cb").value
			var xhr=new XMLHttpRequest()
			var fd=new FormData()
			var obj=new Object()
			obj.Cb=str
			fd.append('cb',JSON.stringify(obj))
			xhr.onload=function(e) {
				if (this.status==200||this.status==302||this.status==304) {
					me.value="√"
					setTimeout("document.getElementById('btSend').value='"+myvalue+"'",1000)
				}
			}
			xhr.open("POST", "/send", true)
			xhr.send(fd)
		}
		function getUploadServerHost() {
			var strs=location.href.split(":")
			var port=new Number(strs[2].split("/")[0])
			return strs[0]+":"+strs[1]+":"+(port+1).toString()
		}
		function downloadKC(){
			var myos=detectOS()
			if (navigator.language=="zh-CN") {
				if (myos=="Windows") {
					window.location.href="downloadKC?os="+myos
				}else if (myos=="Linux") {
					window.location.href="downloadKC?os="+myos
				}else {
					window.location.href="downloadKC?os="+myos
				}
			}else {
				if (myos=="Windows") {
					window.location.href="downloadKC?os="+myos
				}else if (myos=="Linux") {
					window.location.href="downloadKC?os="+myos
				}else {
					window.location.href="downloadKC?os="+myos
				}
			}
		}
        function detectOS() {
            var sUserAgent = navigator.userAgent;
            var isWin = (navigator.platform == "Win32") || (navigator.platform == "Windows");
            if (isWin) return "Windows"
            var isMac = (navigator.platform == "Mac68K") || (navigator.platform == "MacPPC") || (navigator.platform == "Macintosh") || (navigator.platform == "MacIntel");
            if (isMac) return "Mac";
            var isLinux = (String(navigator.platform).indexOf("Linux") > -1)&&(String(navigator.platform).indexOf("Android") <1);
            if (isLinux) return "Linux";
            return "other";
        }
        function detectLang(){
            if (navigator.language=="zh-CN") {
            	document.getElementById("kcenabled").innerHTML="请先下载受控端程序，才能实现遥控"
            	document.getElementById("uploadPageBt").value="上传文件到手机"
            	document.getElementById("liveStreamPage").value="局域网直播"
            }
        }
        setTimeout("detectLang()", 200)
	</script>
	<style type="text/css">
		.wrapper{
			background-color: #fff;
			display: inline-block;
			padding: 5px;
			box-shadow: 2px 2px 10px #000;
			border-radius: 10px;
		}
	</style>
</head>
<body style="background-color: #58c6d5;font-family: sans-serif;">
<center>
{{if .CbEnabled}}
<div class="wrapper"><table>
	<tr>
		<th style="color: #D81B60">{{.Clipboard}}</th>
	</tr>
	<tr>
		<td>
			<textarea name="cb" id="cb" cols="30" rows="5">{{.CbContent}}</textarea>
		</td>
		<td>
			<input type="button" value="{{.Copy}}" id="btCopy" onclick="copy(this)"><br>
			<br>	
			<input type="button" value="{{.Send}}" id="btSend" onclick="send(this)">
		</td>
		<td><span id="spanInfo"></span><br><span></span></td>
	</tr>
	<tr>
	<td colspan="2"><hr></td>
	</tr>
</table></div>
{{end}}
<br><br>
<div class="wrapper">
<table>
	<tr>
		<th style="color: #1E88E5">{{.Files}}</th>
	</tr>
	<tr>
		<td colspan="2">
		{{range .FileSlice}}
			<a href="/downloadFile?filename={{.FileName}}">
			{{.FileName}}</a><br>
		{{end}}
		</td>
	</tr>
	<tr>
	<td colspan="2"><hr></td>
	</tr>
	<tr><td>
			<input type="button" int8 onclick="location.href=getUploadServerHost()" id="uploadPageBt" value="Upload"></input>
	</td></tr>	
	<tr><td>
			<input type="button" int8 onclick="location.href='/showtime'" id="liveStreamPage" value="Live"></input>
	</td></tr>
</table>
</div><br>

{{if .KC_enabled}}
<br><br>
	<b><u id="kcenabled" style="color: blue;cursor: pointer;" onclick="downloadKC()">Download Controlled End Program , to enable remote control</u></b>
<br>
{{end}}

</center>
</body>
</html>
`)
	t.Execute(w, tdata)

}
func SetCBEnabled(bo bool) {
	tdata.CbEnabled = bo
}
func AddFile(str string) {
	if !contains(str, "/") {
		return
	}
	tdata.FileSlice = append(tdata.FileSlice, MyFileEntry{FilePath: str, FileName: getFileName(str)})
	return
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

func DeleteFile(index int) {
	if index > -1 && index < len(tdata.FileSlice) {
		tdata.FileSlice = append(tdata.FileSlice[:index], tdata.FileSlice[index+1:]...)
	}
	return
}
func SetStoragePath(str string) {
	str_storagePath = str
}
func SetClipboard(str string) {
	tdata.CbContent = str
	fmt.Println("CbContent has been set to:" + str)
}
func send(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("cb"))
	var gobj struct {
		Cb string
	}
	e := json.Unmarshal([]byte(r.FormValue("cb")), &gobj)
	if e != nil {
		t := template.New("name")
		t.Parse(`<!DOCTYPE html>
<html>
<head>
	<title>LanGenius</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width maximum-scale=1 initial-scale=1">
</head>
<body>
<center><b>` + e.Error() + `</b></center>
<script type="text/javascript">
	setTimeout("location.href='/'", 1000)
</script>
</body>
</html>`)
		t.Execute(w, nil)
		return
	}
	javahandler.OnClipboardReceived(gobj.Cb)
	SetClipboard(gobj.Cb)
	http.Redirect(w, r, "/", http.StatusFound)
}
func downloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	for _, v := range tdata.FileSlice {
		if v.FileName == filename {
			w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			http.ServeFile(w, r, v.FilePath)
			return
		}
	}
	http.NotFound(w, r)
}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fhs := r.MultipartForm.File["myUploadFile"]
	for _, v := range fhs {
		file, err := v.Open()
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}
		javahandler.OnFileReceived(v.Filename)
		mf, err := os.OpenFile(str_storagePath+v.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}
		defer mf.Close()
		io.Copy(mf, file)
	}
	t := template.New("name")
	t.Parse(`<!DOCTYPE html>
<html>
<head>
	<title>LanGenius</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width maximum-scale=1 initial-scale=1">
</head>
<body>
<center><b>OK</b></center>
<script type="text/javascript">
	setTimeout("location.href='/'", 1000)
</script>
</body>
</html>`)
	t.Execute(w, nil)
}
func downloadKC(w http.ResponseWriter, r *http.Request) {
	var myos = r.FormValue("os")
	if myos == "Linux" {
		filename := "kc_linux_x64"
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		http.ServeFile(w, r, "/data/data/com.xchat.stevenzack.langenius/"+filename)
	} else {
		filename := "kc_windows_x64.exe"
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
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
		fmt.Println(err)
		return
	}
	go sendKCPulse()
	go readKC(c)
}
func StopKC() {
	tdata.KC_enabled = false
	c.Close()
}
func SendKC(cmd string, index int) {
	if len(kc_adds) > 0 && index > -1 && index < len(kc_adds) {
		c.WriteToUDP([]byte(cmd), kc_adds[index])
		return
	} else {
		return
	}
}
func readKC(c *net.UDPConn) {
	b := make([]byte, 512)
	for {
		n, ra, err := c.ReadFromUDP(b)
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(b[:n]) == "LanGenius-from-desktop" {
			if !isExisted(kc_adds, ra) {
				kc_adds = append(kc_adds, ra)
				kchandler.OnDeviceDetected(ra.String())
			}
		}
	}
}
func sendKCPulse() {
	for {
		c.WriteToUDP([]byte("LanGenius-from-android"), kc_adbr)
		time.Sleep(time.Second * 3)
	}
}
func isExisted(as []*net.UDPAddr, a *net.UDPAddr) bool {
	for _, v := range as {
		if v.String() == a.String() {
			return true
		}
	}
	return false
}
func GetUDPConnections() string {
	var a []string
	for _, v := range kc_adds {
		a = append(a, v.String())
	}
	return join(a, "#")
}
func join(a []string, s string) string {
	r := ""
	for k, v := range a {
		r += v
		if k != len(a)-1 {
			r += s
		}
	}
	return r
}
func getIP() string {
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

//showtime part ===========================================================================================================================

type userinfo struct {
	name   string
	img    string
	isself bool
}
type Client struct {
	id   string
	conn *websocket.Conn
	userinfo
}

type message struct {
	Data  string
	Mtype string
	Img   string
}

var member = make(map[string]*Client)

func getclient(ws *websocket.Conn) string {
	for k, v := range member {
		if v.conn == ws {
			return k
		}
	}
	return ""
}

func getnun() string {
	rnd, _ := rand.Int(rand.Reader, big.NewInt(12))
	num := fmt.Sprintf("%v", rnd)
	return num
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func guid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func (m *Client) addclient(ws *websocket.Conn) *Client {
	m.conn = ws
	return m
}

var username string
var uzb string = ""
var ShowtimeFolder string

func pwint(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
	uid := guid()
	//fmt.Println(i)
	if username == "" && uzb == "女主播" {
		username = uzb
	}
	user := userinfo{fmt.Sprintf("%s：", username), fmt.Sprintf("/public/images/%s.jpg", getnun()), false}
	username = ""
	client := Client{uid, ws, user}
	client.addclient(ws)
	member[uid] = &client
	for {
		var err error
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			break
		}
		for k, v := range member {
			if v.conn != ws {
				var mymes message
				json.Unmarshal([]byte(reply), &mymes)
				if mymes.Mtype == "mess" {
					mymes.Data = fmt.Sprintf("%s%s", member[getclient(ws)].userinfo.name, mymes.Data)
					mymes.Img = member[getclient(ws)].userinfo.img
				}
				msg, _ := json.Marshal(mymes)
				if err = websocket.Message.Send(v.conn, string(msg)); err != nil {
					delete(member, k)
					fmt.Println("LiveGoServer:", err)
					break
				}
			} else {
				if !v.userinfo.isself {
					var mymesss message
					json.Unmarshal([]byte(reply), &mymesss)
					//if mymesss.Mtype == "mess" {
					mymesss.Mtype = "self"
					mymesss.Img = v.userinfo.img
					msg, _ := json.Marshal(mymesss)
					member[k].userinfo.isself = true
					if err = websocket.Message.Send(ws, string(msg)); err != nil {
						delete(member, k)
						fmt.Println("LiveGoServer:", err)
						break
					}
					//}
				}
			}

		}
	}
}

func camera(w http.ResponseWriter, r *http.Request) {
	uzb = "zhibo"
	if r.Method == "GET" {
		t, _ := template.ParseFiles(ShowtimeFolder + "/views/camera.html")
		t.Execute(w, nil)
	} else {

	}
}
func live(w http.ResponseWriter, r *http.Request) {
	username = "name"
	if r.Method == "GET" {
		t, _ := template.ParseFiles(ShowtimeFolder + "/views/live.html")
		t.Execute(w, nil)
	}
}
