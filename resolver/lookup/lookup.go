package lookup

import (
	"context"
	"fmt"
	"net"
)

// host == www.google.com
func lookupHost(host, dns string) []string {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}
	ips, err := resolver.LookupHost(context.Background(), host)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	fmt.Println(ips)

	return ips
}
