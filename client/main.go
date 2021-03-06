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
	"crypto/ecdsa"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/tronprotocol/grpc-gateway/api"
	pbc "github.com/tronprotocol/grpc-gateway/core"
	"github.com/zhangheng1536/tron-go/base58"
	"github.com/zhangheng1536/tron-go/common/hexutil"
	"github.com/zhangheng1536/tron-go/crypto"
	"os"
	"sync"
)

const (
	//address        = "34.233.96.87:50051"
	address        = "18.196.99.16:50051"
	defaultAddress = "TTG4qqSRwUBFRvDrhUk1VR8MsMF4vURD8k"
	toAddress      = "TPzvWP6qU44KxZKnGNR1idhbPNEmzguaTe"
	privateKey     = "586C04B441E9B33FF1C621003ECC9F9953D68BFB8A837971EA33A6DFB727D1B9"
)

func main() {
	// Set up a connection to the server.
	//getAccount()
	//createTransaction()
	findAccount()
	//sendCoin()
	//getBlance()
	//testNum(address)
	//lll := make([]string, 0)
	//for i := 0; i < 10; i++ {
	//	go setlll(&lll, i)
	//}
	//time.Sleep(100000)
	//for i := 0; i < 10; i++ {
	//	fmt.Println(lll[i])
	//}
}
func setlll(strings *[]string, i int) {
	fmt.Println("----", i)
	fmt.Println("---------", len(*strings))
	*strings = append(*strings, fmt.Sprintf("%d", i))
}

func testNum(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewDatabaseClient(conn)
	defer func() {
		fmt.Println("-------end--------")
		conn.Close()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetNowBlock(ctx, &pb.EmptyMessage{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("now num: %d", uint64(r.GetBlockHeader().GetRawData().Number))
	dp, err := c.GetDynamicProperties(ctx, &pb.EmptyMessage{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("solidity num: %d", uint64(dp.LastSolidityBlockNum))
}

func getBlance() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewWalletClient(conn)
	defer func() {
		fmt.Println("-------end--------")
		conn.Close()
	}()
	for true {
		addr := ""
		fmt.Scanln(&addr)
		fmt.Printf("address is :%s\n balance is:%d\n", addr, uint64(getAccount(c, addr).Balance))
	}
}

func sendCoin() {
	resultList := make([]string, 0)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewWalletClient(conn)
	defer conn.Close()
	for i := 0; i < 50; i++ {
		if i%10 == 0 {
			fmt.Print(i)
		}
		pk := createPK()
		createTransaction(c, privateKey, base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes()), 1000)
		saveAcc(pk)
		resultList = append(resultList, base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
	}
	saveResult(resultList)
}

func findAccount() {
	log.Println("start find address...")

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewWalletClient(conn)
	defer conn.Close()
	i := 0
	wg := sync.WaitGroup{}
	for true {
		wg.Add(1)
		go funcName(c, &wg)
		i++
		if 0 == i%200 {
			wg.Wait()

		}
		if 0 == i%10000 {
			i = 0
			log.Println("ten thousand continue ...")
		}
	}
}

func funcName(c pb.WalletClient, wg *sync.WaitGroup) {
	defer wg.Done()
	pk := createPK()
	to58Check := base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes())
	account := getAccount(c, to58Check)

	if to58Check[:3] == "Trx" ||
		to58Check[:3] == "TRX" ||
		to58Check[:4] == "TRON" ||
		to58Check[:4] == "Tron" ||
		(nil != account && account.Balance > 0) {
		saveAcc(pk)
	}
}

func saveAcc(pk *ecdsa.PrivateKey) {
	file, err := os.Create(fmt.Sprintf("res/%s.txt", base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes())))
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintf(file, hexutil.Encode(pk.D.Bytes()))
}
func saveResult(result []string) {
	file, err := os.Create("res/result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprint(file, result)
}

func createTransaction(c pb.WalletClient, priveteK string, toAddr string, amount int64) {
	pk, err := crypto.HexToECDSA(priveteK)
	if nil != err {
		log.Fatalf("did not priveteKey: %v", err)
	}
	toAddressBytes, _ := base58.From58Check(toAddr)
	contract := &pbc.TransferContract{
		OwnerAddress: crypto.PubkeyToAddress(pk.PublicKey).Bytes(),
		ToAddress:    []byte(toAddressBytes),
		Amount:       amount,
	}
	//anyOne,err:=ptypes.MarshalAny(param)
	//if nil!=err{
	//	log.Fatalf("did not anyOne: %v", err)
	//}
	//contract:=&pbc.Transaction_Contract{
	//	Parameter:anyOne,
	//	Type:pbc.Transaction_Contract_TransferContract,
	//}
	//raw:=&pbc.TransactionRaw{
	//	Contract:contract,
	//}
	//transaction := &pbc.Transaction{
	//	RawData: raw,
	//}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateTransaction(ctx, contract)
	if err != nil {
		log.Fatalf("could not transaction: %v", err)
	}
	log.Printf("transaction is: %v", r)

	data, err := proto.Marshal(r.RawData)
	if nil != err {
		log.Fatalf("did not RawData data: %v", err)
	}
	signData := base58.Sha256Do(data)
	sign, err := crypto.Sign(signData, pk)
	if nil != err {
		log.Fatalf("did not Sign: %v", err)
	}
	r.Signature = append(r.Signature, sign)
	ret, err := c.BroadcastTransaction(ctx, r)
	if nil != err {
		log.Fatalf("did not CreateTransaction: %v", err)
	}
	log.Printf("return is: %v", ret)
}

func createPK() *ecdsa.PrivateKey {
	pk, err := crypto.GenerateKey()
	if nil != err {
		log.Fatalf("did not GenerateKey: %v", err)
	}
	//fmt.Println("pk is :", hexutil.Encode(pk.D.Bytes()))
	////fmt.Println("pk is :", hexutil.Encode(pk.PublicKey.X.Bytes()))
	////fmt.Println("pk is :", hexutil.Encode(pk.PublicKey.Y.Bytes()))
	//fmt.Println("addressHex is :", hexutil.Encode(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
	//fmt.Println("addressCheck58 is :", base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
	return pk
}

func getAccount(c pb.WalletClient, addr string) *pbc.Account {

	// Contact the server and print out its response.
	address, _ := base58.From58Check(addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetAccount(ctx, &pbc.Account{Address: address})
	if err != nil {
		log.Println("could not greet: %v", err)
		return nil
	}
	//log.Printf("Greeting: %v", r)
	return r
}
