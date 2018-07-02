package main

import (
	"fmt"
	"github.com/zhangheng1536/tron-go/base58"
	"github.com/zhangheng1536/tron-go/common/hexutil"
	"github.com/zhangheng1536/tron-go/crypto"
	"log"
	"net/http"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Path[1:])
	if nil != err {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprintf(w, "Hi there, I love %d!\n\n", count)
	for i := 1; i <= count; i++ {
		fmt.Fprintf(w, "---------num %d------------\n", i)
		createPK(&w)
		fmt.Fprintf(w, "------------end------------\n\n")
	}
}

func main() {
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8088", nil))
	for true {
		go time.Sleep(100000)
	}
}

func createPK(w *http.ResponseWriter) {
	pk, err := crypto.GenerateKey()
	if nil != err {
		log.Fatalf("did not GenerateKey: %v", err)
	}
	fmt.Fprintf(*w, "pk is:\t%s\n", hexutil.Encode(pk.D.Bytes()))
	fmt.Fprintf(*w, "addressHx is:\t%s\n", hexutil.Encode(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
	fmt.Fprintf(*w, "addressCheck58 is:\t%s\n", base58.To58Check(crypto.PubkeyToAddress(pk.PublicKey).Bytes()))
}
