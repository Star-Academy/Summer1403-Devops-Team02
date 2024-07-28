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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Description struct {
	Address  string `json:"address"`
	Duration string `json:"duration"`
}

type Record struct {
	ID          int         `json:"id"`
	Description Description `json:"description"`
}

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func saveToRedis(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func homePage(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	if addr == "" {
		http.Error(w, "Missing 'addr' parameter", http.StatusBadRequest)
		return
	}

	records := traceRoute(addr)
	jsonResponse, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamp := time.Now()
	err = saveToRedis(addr+":"+timestamp.String(), records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}

const (
	ProtocolICMP = 1
)

func traceRoute(addr string) []Record {
	var records []Record
	c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for i := 1; i <= 30; i++ {
		finished, dst, dur, err := getNthHop(c, addr, i)
		if err != nil {
			records = append(records, Record{ID: i, Description: Description{Address: "unknown", Duration: "*"}})
			fmt.Printf("%d * * *\n", i)
		} else {
			records = append(records, Record{ID: i, Description: Description{Address: dst.String(), Duration: dur.String()}})
			fmt.Printf("%d %s (%s)\n", i, dst, dur)

			if finished {
				fmt.Println("DONE!")
				break
			}
		}
	}
	return records
}

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

func getNthHop(c *icmp.PacketConn, addr string, ttl int) (bool, *net.IPAddr, time.Duration, error) {
	// Start listening for icmp replies
	c.IPv4PacketConn().SetTTL(ttl)

	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		panic(err)
	}

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, // << uint(seq), // TODO
			Data: []byte(""),
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
	err = c.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
	if err != nil {
		return false, dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return false, dst, 0, err
	}
	duration := time.Since(start)

	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return false, dst, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return true, dst, duration, nil
	case ipv4.ICMPTypeTimeExceeded:
		return false, &net.IPAddr{IP: peer.(*net.IPAddr).IP}, duration, nil
	default:
		return false, dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}
