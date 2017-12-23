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
		deamonConn, e := net.ListenUDP("udp", a)
		if e != nil {
			fmt.Println(e)
			return
		}

		broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
		broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Online", Content: mPort, Info: osInfo, RemoteControlStatus: RemoteControlEnabled})
		deamonConn.WriteToUDP(broData, broadcastAddr)
		b := make([]byte, 4096)
		for {
			n, ra, e := deamonConn.ReadFromUDP(b)
			if e != nil {
				fmt.Println(e)
				break
			}
			if IsMyIP(ra.IP.String()) {
				continue
			}
			msg := Msg{}
			e = json.Unmarshal(b[:n], &msg)
			if e != nil {
				continue
			}
			msg.Content = ra.IP.String() + msg.Content
			go routeUdpMsg(msg)
		}
	}()
}
func routeUdpMsg(msg Msg) {
	switch msg.Type {
	case "LanGenius-Deamon":
		handleDeamon(msg)
	case "LanGenius-RemoteControlCmd":
		handleRemoteControlCmd(msg)
	}
}
func handleDeamon(msg Msg) {
	if msg.State == "Online" {
		mEventHandler.OnDeviceOnlineListener(msg)
	} else {
		mEventHandler.OnDeviceOfflineListener(msg)
	}
}
func StopDeamon() {
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Offline", Content: mPort})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
