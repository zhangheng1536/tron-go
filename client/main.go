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
	"fmt"
	pb "github.com/tronprotocol/grpc-gateway/api"
	pbc "github.com/tronprotocol/grpc-gateway/core"
	"github.com/zhangheng1536/tron-go/base58"
	"github.com/zhangheng1536/tron-go/common/hexutil"
	"github.com/zhangheng1536/tron-go/crypto"
)

const (
	address        = "47.104.214.27:50051"
	defaultAddress = "TRXSsMSfYgqhFhBzCYNHKD9adJrXpWZigb"
)

func main() {
	// Set up a connection to the server.
	//getAccount()
	pk, err := crypto.GenerateKey()
	if nil != err {
		log.Fatalf("did not GenerateKey: %v", err)
	}
	fmt.Println("pk is :", hexutil.Encode(pk.D.Bytes()))
	fmt.Println("pk is :", hexutil.Encode(pk.PublicKey.X.Bytes()))
	fmt.Println("pk is :", hexutil.Encode(pk.PublicKey.Y.Bytes()))
	fmt.Println("addressHex is :", hexutil.Encode(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
	fmt.Println("addressCheck58 is :", base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))

}

func getAccount() {
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
