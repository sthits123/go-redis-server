package main

import (
	"log"
	"net"
	"os"
	
)



func main() {
   log.Println("Reading config file...")
   
   readConf("./redis.conf")

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot listen on port 6379: ", err)
	}
	defer listener.Close()
	log.Println("listening on :6379")

	conn,err := listener.Accept()
	if err != nil {
		log.Println( err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println("connection accepted")

	for {
		// logic here
		v:=Value{typ: ARRAY}
		v.readArray(conn)
		handle(conn, &v)
	}
}
