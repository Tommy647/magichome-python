package main

import (
	"log"
	"math/rand"
	"net"
	"time"
)

const port = "5577"

const (
	// modes - 3rd bit
	seven_color_cross_fade   = 0x25
	red_gradual_change       = 0x26
	green_gradual_change     = 0x27
	blue_gradual_change      = 0x28
	yellow_gradual_change    = 0x29
	cyan_gradual_change      = 0x2a
	purple_gradual_change    = 0x2b
	white_gradual_change     = 0x2c
	red_green_cross_fade     = 0x2d
	red_blue_cross_fade      = 0x2e
	green_blue_cross_fade    = 0x2f
	seven_color_strobe_flash = 0x30
	red_strobe_flash         = 0x31
	green_strobe_flash       = 0x32
	blue_stobe_flash         = 0x33
	yellow_strobe_flash      = 0x34
	cyan_strobe_flash        = 0x35
	purple_strobe_flash      = 0x36
	white_strobe_flash       = 0x37
	seven_color_jumping      = 0x38
)
const max = 254

var status = []byte{0x81, 0x8A, 0x8B, 0x96}
var on = []byte{0x71, 0x23, 0x0F, 0xA3}
var off = []byte{0x71, 0x24, 0x0F, 0xA4}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var red uint8 = 0
	var green uint8 = 0
	var blue uint8 = 0

	ip := "192.168.0.179"
	log.Print("Establishing connection with the device.")
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		panic(err.Error())
	}

	getStatus(conn)
	turnOn(conn)
	for i := 0; i < 10; i++ {
		red = 0x00 // uint8(r.Int63n(max))
		green = uint8(r.Int63n(max))
		blue = 0x00 // uint8(r.Int63n(max))
		time.Sleep(2 * time.Second)
		update(conn, red, green, blue)
		getStatus(conn)
	}
	turnOff(conn)
	conn.Close()
}

func getStatus(conn net.Conn) {
	_ = conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	i, err := conn.Write(status)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("sent %d bytes", i)

	buf := make([]byte, 14)
	i, err = conn.Read(buf)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("recieved %d bytes", i)
	log.Printf("got a status message %#v", buf)
}

func turnOn(conn net.Conn) {
	_ = conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	log.Print("turn on")
	i, err := conn.Write(on)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("sent %d bytes", i)

}

func turnOff(conn net.Conn) {
	log.Print("turn off")
	_ = conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	i, err := conn.Write(off)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("sent %d bytes", i)
}

func update(conn net.Conn, red, green, blue uint8) {
	log.Printf("updating to red:%d green:%d blue:%d", red, green, blue)
	message := []byte{0x31, red, green, blue, 0xff, 0x00, 0x0f}
	checksum := checksum(message)
	message = append(message, uint8(checksum))
	log.Printf("updating %v", message)
	_ = conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	i, err := conn.Write(message)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("sent %d bytes", i)
}

func checksum(message []byte) uint8 {
	x := sum(message)
	log.Printf("& 0xff %d", x&0xff)
	return x & 0xff
}

func sum(message []byte) uint8 {
	var x uint8 = 0x00
	for i := range message {
		log.Printf("summing bytes %d", x)
		x += uint8(message[i])
	}
	log.Printf("summed up bytes %d", x)

	return x
}
