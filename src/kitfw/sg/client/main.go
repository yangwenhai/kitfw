package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"kitfw/sg/pb"

	protocol "kitfw/sg/protocol"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8081"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKitfwClient(conn)

	// Contact the server and print out its response.
	if len(os.Args) < 3 {
		log.Fatalf("args not enough")
		return
	}

	//encode capnp
	userid, _ := strconv.ParseInt(os.Args[1], 10, 64)
	req := &protocol.ConcatRequest{
		UserId: userid,
		Str1:   os.Args[2],
		Str2:   os.Args[3],
	}
	payload, _ := protocol.Encode(req)

	//send request
	r, err := c.Process(context.Background(), &pb.KitfwRequest{
		Protoid: protocol.PROTOCOL_CONCAT_REQUEST,
		Logid:   "123456",
		Payload: payload,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//decode capnp
	res := &protocol.ConcatReply{}
	err = protocol.Decode(res, r.Payload)
	if err != nil {
		fmt.Println("err!", err)
		return
	}

	log.Printf("response: %s", res.Val)
}
