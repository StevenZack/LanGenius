package LanGenius

import (
	"encoding/json"
	"fmt"
	"net"
)

var (
	DeamonPort string = ":12812"
	deamonConn *net.UDPConn
)

func StartDeamon() {
	go func() {
		a, _ := net.ResolveUDPAddr("udp", DeamonPort)
		deamonConn, e := net.ListenUDP("udp", a)
		if e != nil {
			fmt.Println(e)
			return
		}

		broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
		broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Online", Content: mPort})
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
			go routeUdpMsg(msg, ra.IP.String())
		}
	}()
}
func routeUdpMsg(msg Msg, rip string) {
	switch msg.Type {
	case "LanGenius-Deamon":
		handleDeamon(msg, rip)
	}
}
func handleDeamon(msg Msg, rip string) {
	if msg.State == "Online" {
		mEventHandler.OnDeviceOnlineListener(rip + msg.Content)
	} else {
		mEventHandler.OnDeviceOfflineListener(rip + msg.Content)
	}
}
func StopDeamon() {
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Offline", Content: mPort})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
