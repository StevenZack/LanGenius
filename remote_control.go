package LanGenius

import (
	"encoding/json"
	// "fmt"
	"net"
)

var (
	RemoteControlEnabled bool
)

func handleRemoteControlCmd(msg Msg) {
	b, _ := json.Marshal(msg)
	if RemoteControlEnabled {
		mEventHandler.OnRemoteControlCmdReceived(string(b))
	}
}
func SetRemoteControlStatus(b bool) {
	RemoteControlEnabled = b
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Online", Content: mPort, Info: osInfo, RemoteControlStatus: RemoteControlEnabled})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
func SendRemoteControl(msg Msg) {
	sendAddr, _ := net.ResolveUDPAddr("udp", msg.Content)
	msg.Type = "LanGenius-RemoteControlCmd"
	b, _ := json.Marshal(msg)
	deamonConn.WriteToUDP(b, sendAddr)
}
