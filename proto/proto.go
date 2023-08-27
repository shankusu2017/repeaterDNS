package proto

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func Buf2DNSReq(buf []byte) *layers.DNS {
	packet := gopacket.NewPacket(buf, layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	req, _ := dnsPacket.(*layers.DNS)
	return req
}
