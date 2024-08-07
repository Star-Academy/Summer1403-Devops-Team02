// Copyright © 2016 Alex
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	ProtocolICMP = 1
	//ProtocolIPv6ICMP = 58
	ListenAddr = "0.0.0.0"
)

type TraceHop struct {
	IPAddr net.IP
	RTT    time.Duration
}

func TraceRoute(addr string, maxHops int) ([]*TraceHop, error) {
	c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		return nil, err
	}

	res := make([]*TraceHop, maxHops)

	for i := 0; i < maxHops; i++ {
		finished, dst, dur, err := getNthHop(c, dst, i+1)
		if err != nil {
			res[i] = nil
		} else {
			res[i] = &TraceHop{IPAddr: dst.IP, RTT: dur}
			if finished {
				return res[:i+1], nil
			}
		}

	}

	return res, nil
}

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

func getNthHop(c *icmp.PacketConn, dst *net.IPAddr, ttl int) (bool, *net.IPAddr, time.Duration, error) {
	// Start listening for icmp replies
	c.IPv4PacketConn().SetTTL(ttl)

	data := make([]byte, 64)
	rand.Read(data)

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, //<< uint(seq), // TODO
			Data: data,
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return false, dst, 0, err
	}

	// Send it
	start := time.Now()

	n, err := c.WriteTo(b, dst)
	if err != nil {
		return false, dst, 0, err
	} else if n != len(b) {
		return false, dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	if err != nil {
		return false, dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		// fmt.Println("Unable to read!")
		return false, dst, 0, err
	}
	duration := time.Since(start)

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return false, dst, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return true, dst, duration, nil
	case ipv4.ICMPTypeTimeExceeded:
		// Convert peer to IPAddr
		return false, &net.IPAddr{IP: peer.(*net.IPAddr).IP}, duration, nil
	default:
		return false, dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}
