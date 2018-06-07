/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"log"
	//"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"github.com/zhangheng1536/tron-go/common/base58"
	pb "github.com/zhangheng1536/tron-go/gen/api"
	pbc "github.com/zhangheng1536/tron-go/gen/core"
)

const (
	address        = "34.214.241.188:50051"
	defaultAddress = "TWsm8HtU2A5eEzoT8ev8yaoFjHsXLLrckb"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWalletClient(conn)

	// Contact the server and print out its response.
	address, _ := base58.From58Check(defaultAddress)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetAccount(ctx, &pbc.Account{Address: address})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r)
}
