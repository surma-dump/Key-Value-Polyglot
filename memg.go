#!/usr/bin/env gorun
package main;

import (
    "fmt"
    "net"
    "bufio"
    "strings"
    "os"
    "strconv"
)

const (
    PORT = "11211"
)

func main() {

    listener, err := net.Listen("tcp", "127.0.0.1:" + PORT)
    if err != nil {
        panic("Error listening on " + PORT + ": " + err.String())
    }

    cache := NewCache()

    netconn, err := listener.Accept()
    if err != nil {
        panic("Listener accept error: " + err.String())
        return
    }

    handleConn(netconn, cache)
}

/*
 * Networking
 */
func handleConn(conn net.Conn, cache *Cache) {

    reader := bufio.NewReader(conn)
    for {

        // Fetch

        content, err := reader.ReadString('\n')
        if err == os.EOF {
            break
        } else if err != nil {
            fmt.Println(err)
            return
        }

        content = content[:len(content) - 2]    // Chop \r\n

        // Handle

        parts := strings.Split(content, " ")
        cmd := parts[0]
        switch cmd {

        case "get":
            key := parts[1]
            val := cache.get(key)
            length := strconv.Itoa(len(val))

            conn.Write([]uint8("VALUE " + key + " 0 " + length + "\r\n"))
            conn.Write([]uint8(val + "\r\n"))
            conn.Write([]uint8("END\r\n"))

        case "set":
            key := parts[1]
            //exp := parts[2]
            //flags := parts[3]
            length, _ := strconv.Atoi(parts[4])
            // Really we should read exactly 'length' bytes + \r\n
            val, _ := reader.ReadString('\n')
            val = val[:length]    // Chop to length given by client
            cache.set(key, val)
            conn.Write([]uint8("STORED\r\n"))
        }
    }
}

/*
 * Cache
 */

type Cache struct {
    store map[string] string
};

func NewCache() *Cache {
    store := make(map[string] string)
    return &Cache{store: store}
}

func (self *Cache) set(key, val string) {
    self.store[key] = val;
}

func (self *Cache) get(key string) string {
    return self.store[key];
}
