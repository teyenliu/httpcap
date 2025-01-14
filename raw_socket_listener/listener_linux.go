package raw_socket

import (
	"fmt"
	_ "fmt"
	"log"
	"net"
	_ "strconv"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (t *Listener) readRAWSocket() {
	// AF_INET can't capture outgoing packets, must change to use AF_PACKET
	// https://github.com/golang/go/issues/7653
	// http://www.binarytides.com/packet-sniffer-code-in-c-using-linux-sockets-bsd-part-2/
	proto := (syscall.ETH_P_ALL<<8)&0xff00 | syscall.ETH_P_ALL>>8 // change to Big-Endian order
	fmt.Println("[DEBUG] proto: ", proto)
	//fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, proto)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatal("socket: ", err)
	}
	defer syscall.Close(fd)
	if t.addr != "" && t.addr != "0.0.0.0" {
		ifi, err := net.InterfaceByName(t.addr)
		if err != nil {
			log.Fatal("interfacebyname: ", err)
		}
		lla := syscall.SockaddrLinklayer{Protocol: uint16(proto), Ifindex: ifi.Index}
		if err := syscall.Bind(fd, &lla); err != nil {
			log.Fatal("bind: ", err)
		}
	}

	/*
		var pipefile = "/tmp/pipe.ipc"
		os.Remove(pipefile)
		err = syscall.Mkfifo(pipefile, 0666)
		if err != nil {
			log.Fatal("create named pipe error:", err)
		}

		fmt.Println("start schedule writing.")
		//f, err := os.OpenFile(pipefile, os.O_RDWR, 0777)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
	*/

	//named pipes for IPC communication
	//var src_ip string
	//var dest_ip string
	buf := make([]byte, 65536)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		if n <= 0 {
			continue
		}

		//packet := gopacket.NewPacket(buf[:n], layers.LayerTypeEthernet, gopacket.Default)
		packet := gopacket.NewPacket(buf[:n], layers.LayerTypeUDP, gopacket.Default)

		//fmt.Println("[DEBUG] packet:", packet.Dump())
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip := ipLayer.(*layers.IPv4)
			dst := ip.DstIP.String()
			src := ip.SrcIP.String()
			fmt.Printf("IP From %s to %s\n\n", src, dst)
			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				dst = fmt.Sprintf("%s:%d", dst, udp.DstPort)
				src = fmt.Sprintf("%s:%d", src, udp.SrcPort)
				fmt.Printf("UDP From %s to %s\n\n", src, dst)
			}
		}

		// send packets's byte data to named pipes
		//f.Write(packet.Data())

		//if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		//	tcp, _ := tcpLayer.(*layers.TCP)

		//	src_ip = packet.NetworkLayer().NetworkFlow().Src().String()
		//	dest_ip = packet.NetworkLayer().NetworkFlow().Dst().String()

		//	t.parsePacket(src_ip, dest_ip, tcp)
		//}
	}
}
