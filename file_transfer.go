package LanGenius

import (
	"html/template"
	"net/http"
)

type EventHandler interface {
	OnClipboardReceived(string)
	OnFileReceived(string)
}

var (
	mEventHandler EventHandler
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
	t := template.New("homeTPL")
	t.Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>LanGenius - Transfer file easily</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-widthinitial-scale=1">
	<script>function ChangeLanguage(){(navigator.language="zh-CN")?(document.getElementById("btCopy").value="复制",document.getElementById("btSend").value="发送",document.getElementById("txtClipboard").innerHTML="剪切板",document.getElementById("txtFiles").innerHTML="文件",document.getElementById("liveStreamPage").value="直播页面",window.Succeed="成功",document.getElementById("submit_button").value="上传文件"):window.Succeed=" Succeed"}function copy(e){document.getElementById("cb").select(),document.execCommand("Copy");var t=e.value;e.value+=window.Succeed,e.disabled="disabled",setTimeout("document.getElementById('btCopy').value='"+t+"';document.getElementById('btCopy').disabled=''",1e3)}function send(e){var t=e.value,n=document.getElementById("cb").value,d=new XMLHttpRequest,o=new FormData,u=new Object;u.Cb=n,o.append("cb",JSON.stringify(u)),d.onload=function(n){200!=this.status&&302!=this.status&&304!=this.status||(e.value+=window.Succeed,setTimeout("document.getElementById('btSend').value='"+t+"'",1e3))},d.open("POST","/send",!0),d.send(o)}function DoUpload(){for(var e=new FormData,t=document.getElementById("myUploadFile"),n=0;n<t.files.length;n++)e.append("myUploadFile",t.files[n]);uploading();var d=new XMLHttpRequest;d.upload.addEventListener("progress",function(e){if(e.lengthComputable){var t=Math.round(100*e.loaded/e.total);document.getElementById("uploadProgress").value=t,document.getElementById("result").innerHTML=t.toString()+"%"}else document.getElementById("result").innerHTML="unable to compute"},!1),d.onreadystatechange=function(){4==d.readyState&&200==d.status&&(document.getElementById("result").innerHTML=d.responseText,uploadDone())},d.addEventListener("load",function(e){document.getElementById("result").innerHTML=e.target.responseText,uploadDone()},!1),d.addEventListener("error",function(e){document.getElementById("result").innerHTML="Something went wrong .",uploadDone()},!1),d.addEventListener("abort",function(e){document.getElementById("result").innerHTML="Abort .",uploadDone()},!1),d.open("POST","/upload"),d.send(e)}setTimeout("ChangeLanguage()",50)</script>
	<style>.wrapper{background-color:#fff;display:inline-block;padding:5px;box-shadow:2px 2px 10px #000;border-radius:10px}</style>
	</head>
	<body style="background-color:#58c6d5;font-family:sans-serif;padding:10px;margin:0">
	<center>
	{{if .ClipboardEnabled}}
	<div class="wrapper">
	<table>
	<tr>
	<th style="color:#d81b60" id="txtClipboard">Clipboard</th>
	</tr>
	<tr>
	<td>
	<textarea name="cb" id="cb" cols="30" rows="5">{{.CbContent}}</textarea>
	</td>
	<td>
	<input type="button" value="Copy" id="btCopy" onclick="copy(this)">
	<br>
	<br>
	<input type="button" value="Send" id="btSend" onclick="send(this)">
	</td>
	<td><span id="spanInfo"></span>
	<br><span></span></td>
	</tr>
	<tr>
	<td colspan="2">
	<hr>
	</td>
	</tr>
	</table>
	</div>
	{{end}}
	<br>
	<br>
	<div class="wrapper">
	<table>
	<tr>
	<th style="color:#1e88e5" id="txtFiles">Files</th>
	</tr>
	<tr>
	<td colspan="2">
	{{range .Files}}
	<a href="/downloadFile?filename={{.Name}}">
	{{.Name}}</a>
	<br>{{end}}
	</td>
	</tr>
	<tr>
	<td colspan="2">
	<hr>
	</td>
	</tr>
	<tr>
	<td>
	<input type="file" name="myUploadFile" id="myUploadFile" onchange='document.getElementById("submit_button").disabled=""' multiple>
	<input type="button" id="submit_button" disabled value="Upload" onclick="DoUpload()">
	</td>
	</tr>
	<tr><td>
	<progress value="0" max="100" id="uploadProgress" style="height:4px;width:100%"></progress>
	<div id="result"></div>
	</td></tr>
	<tr>
	<td>
	<input type="button" int8 onclick='location.href="/live"' id="liveStreamPage" value="Live">
	</td>
	</tr>
	</table>
	</div>
	<br>
	</center>
	</body>
	</html>
		`)
	t.Execute(w, homeData)
}
func send(w http.ResponseWriter, r *http.Request) {

}
func downloadFile(w http.ResponseWriter, r *http.Request) {

}
func upload(w http.ResponseWriter, r *http.Request) {

}
