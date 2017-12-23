package LanGenius

import (
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
)

var members []*websocket.Conn

func wsLive(ws *websocket.Conn) {
	defer ws.Close()
	members = append(members, ws)
	for {
		str := ""
		if e := websocket.Message.Receive(ws, &str); e != nil {
			break
		}
		go sendToAll(str)
	}
}
func live(w http.ResponseWriter, r *http.Request) {
	t := template.New("live")
	t.Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Live stream</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width">
	<script>function changeLang(){"zh-CN"==navigator.language&&(document.getElementById("btCamera").innerHTML="我要直播")}setTimeout("changeLang()",50)</script>
	</head>
	<body>
	<a href="/live/camera" id="btCamera">Live stream</a><br>
	<img src="" id="receiver" alt="online(0)">
	<script>var image=document.getElementById("receiver"),domain=location.href.split("/")[2],socket=new WebSocket("ws://"+domain+"/wsLive");socket.onmessage=function(e){var o=JSON.parse(e.data);"video"==o.Mtype&&(image.src=o.Data)},socket.onopen=function(e){console.log("connected")},socket.onclose=function(e){console.log("closed")}</script>
	</body>
	</html>
		`)
	// t, _ := template.ParseFiles("/home/asd/go/src/LanGenius/views/live.html")
	t.Execute(w, nil)
}
func camera(w http.ResponseWriter, r *http.Request) {
	// t, _ := template.ParseFiles("/home/asd/go/src/LanGenius/views/camera.html")
	// t.Execute(w, nil)
	t := template.New("camera")
	t.Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<title></title>
	<meta name="viewport" content="width=device-width">
	</head>
	<body style="padding:0;margin:0">
	<div id="result"></div>
	<video id="sourcevid" controls></video>
	<br>
	<canvas id="output" style="display:none"></canvas>
	<script>var data,v=document.getElementById("sourcevid"),mcavas=document.getElementById("output"),mcavasContext=mcavas.getContext("2d"),draw=function(){try{mcavas.width=v.videoWidth,mcavas.height=v.videoHeight,mcavasContext.drawImage(v,0,0,v.videoWidth,v.videoHeight)}catch(e){if("NS_ERROR_NOT_AVAILABLE"==e.name)return console.log("NS_ERROR_NOT_AVAILABLE"),setTimeout(draw,33);console.log(e)}v.src&&(data={Data:mcavas.toDataURL("image/jpeg",.5),Mtype:"video"},socket.send(JSON.stringify(data))),setTimeout(draw,33)},socket=new WebSocket(location.href.replace("http","ws").replace("live/camera","wsLive"));socket.onopen=function(e){console.log("connected"),draw()},socket.onclose=function(){console.log("closed")},socket.onerror=function(e){console.log("err:"+e.data)},navigator.getUserMedia=navigator.getUserMedia||navigator.webkitGetUserMedia||navigator.mozGetUserMedia||navigator.msGetUserMedia,navigator.getUserMedia?navigator.getUserMedia({audio:!1,video:!0},function(e){v.src=window.URL.createObjectURL(e),v.onloadedmetadata=function(e){v.play()}},function(e){console.log("The following error occurred: "+e.name),"zh-CN"==navigator.language?document.getElementById("result").innerHTML="由于chrome安全限制，直播者只能使用火狐浏览器，但是观看者仍然可以使用chrome":document.getElementById("result").innerHTML="Live streamer please use firefox web browser,but the viewer can still use chrome"}):console.log("getUserMedia not supported")</script>
	</body>
	</html>
		`)
	t.Execute(w, nil)
}
func sendToAll(str string) {
	for k, v := range members {
		if e := websocket.Message.Send(v, str); e != nil {
			members = append(members[:k], members[k+1:]...) //delete offline Conn in member
			break
		}
	}
}
