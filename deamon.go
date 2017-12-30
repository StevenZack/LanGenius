package LanGenius

import (
	"encoding/json"
	"fmt"
	"net"
)

var (
	DeamonPort string = ":12812"
	deamonConn *net.UDPConn
	osInfo     string
)

func StartDeamon(osI string) {
	osInfo = osI
	go func() {
		a, _ := net.ResolveUDPAddr("udp", DeamonPort)
		deamonConn, _ = net.ListenUDP("udp", a)
		DeamonBroadcastMe()
		b := make([]byte, 4096)
		for {
			n, ra, e := deamonConn.ReadFromUDP(b)
			if e != nil {
				fmt.Println(e)
				break
			}
			if IsMyIP(ra.IP.String()) {
				// fmt.Println("isMyIP:", string(b[:n]))
				continue
			}
			msg := Msg{}
			e = json.Unmarshal(b[:n], &msg)
			if e != nil {
				continue
			}
			msg.IP = ra.IP.String()
			go routeUdpMsg(msg)
		}
	}()
}
func DeamonBroadcastMe() {
	// fmt.Println("online broadcast")
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Online", Port: mPort, Info: osInfo, RemoteControlStatus: RemoteControlEnabled})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
func routeUdpMsg(msg Msg) {
	switch msg.Type {
	case "LanGenius-Deamon":
		handleDeamon(msg)
	case "LanGenius-RemoteControlCmd":
		handleRemoteControlCmd(msg)
	case "LanGenius-Message":
		handleMessage(msg)
	case "LanGenius-FolderReceived":
		handleFolderReceived(msg)
	}
}
func handleDeamon(msg Msg) {
	b, _ := json.Marshal(msg)
	if msg.State == "Online" {
		mEventHandler.OnDeviceOnline(string(b))
		DeamonBroadcastMe()
	} else if msg.State == "Offline" {
		mEventHandler.OnDeviceOffline(string(b))
	}
}
func StopDeamon() {
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Offline", Port: mPort})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
func handleMessage(msg Msg) {
	b, _ := json.Marshal(msg)
	mEventHandler.OnMessageReceived(string(b))
}
func SendMessage(data string) {
	// fmt.Println("sent:", data)
	msg := Msg{}
	json.Unmarshal([]byte(data), &msg)
	sendAddr, _ := net.ResolveUDPAddr("udp", msg.IP+DeamonPort)
	msg.Type = "LanGenius-Message"
	b, _ := json.Marshal(msg)
	deamonConn.WriteToUDP(b, sendAddr)
}
