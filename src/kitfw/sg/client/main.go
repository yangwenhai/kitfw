package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"kitfw/sg/pb"

	protocol "kitfw/sg/protocol"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	rand.Seed(time.Now().UnixNano())
	userid := int64(rand.Intn(10000))%10000 + 10000

	//encode capnp
	req := &protocol.ConcatRequest{
		UserId: userid,
		Str1:   os.Args[1],
		Str2:   os.Args[2],
	}
	payload, _ := protocol.Encode(req)

	//send request
	starttime := time.Now()
	logid := fmt.Sprintf("%d", time.Now().UnixNano()%10000000)
	md := metadata.New(map[string]string{"userid": fmt.Sprintf("%d", userid), "logid": logid})
	// create a new context with this metadata
	ctx := metadata.NewContext(context.Background(), md)
	r, err := c.Process(ctx, &pb.KitfwRequest{
		Protoid: protocol.PROTOCOL_CONCAT_REQUEST,
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

	log.Printf("response: %s took: %fms", res.Val, time.Since(starttime).Seconds()*1000)
}
