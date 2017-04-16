package main

import (
	f "fmt"
	"github.com/go-vgo/robotgo"
	// "crypto/json"
	"net"
	"strings"
)

var (
	addr, rabr, ras *net.UDPAddr
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":9942")
	rabr, _ := net.ResolveUDPAddr("udp", "255.255.255.255:9943")
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		f.Println(err)
		return
	}
	defer c.Close()
	c.WriteToUDP([]byte("LanGenius-from-desktop"), rabr)
	b := make([]byte, 512)
	f.Println("Welcome to Keyboard Controller for windows_x64")
	f.Println("Started successfully !!")
	f.Println("=== My IP : " + getIP() + " =====")
	f.Println("\n\nListening ...")
	for {
		n, ra, err := c.ReadFromUDP(b)
		if err != nil {
			f.Println(err)
			break
		}
		strtemp := string(b[:n])
		if strtemp == "LanGenius-from-android" {
			c.WriteToUDP([]byte("LanGenius-from-desktop"), rabr)
		} else {
			var sts = strings.Split(strtemp, "#")
			ifs := make([]interface{}, len(sts))
			for i, v := range sts {
				ifs[i] = v
			}
			f.Println(reverse(sts), "from", ra.String())
			robotgo.KeyTap(ifs...)
		}
	}
}
func reverse(sts []string) string {
	var sts_r []string
	for i := len(sts) - 1; i > -1; i-- {
		sts_r = append(sts_r, sts[i])
	}
	return strings.Join(sts_r, " ")
}

// func handlePkg(str string, ra *net.UDPAddr, c *net.UDPConn) {
// 	f.Println(str, "( from ", ra.String(), " )")
// 	cmd := strings.Split(str, "#")
// 	f.Println(cmd[0] == "requestIP")
// 	if cmd[0] == "requestIP" {
// 		c.WriteToUDP([]byte(getIP()), ra)
// 	} else {
// 		is := make([]interface{}, len(cmd))
// 		for i, v := range cmd {
// 			is[i] = v
// 		}
// 		robotgo.KeyTap(is...)
// 	}
// }
func getIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		f.Println(err)
		return ""
	}
	var strs []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			f.Println(err)
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
		if strings.HasPrefix(v, "192.168.") {
			return v
		}
	}
	for _, v := range strs {
		if strings.HasPrefix(v, "10.") {
			return v
		}
	}
	for _, v := range strs {
		if strings.HasPrefix(v, "172.") {
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
