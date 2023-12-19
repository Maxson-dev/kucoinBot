package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const START = 1661248797

var (
	buyBody []byte

	wg sync.WaitGroup
)

func init() {
	var err error
	b := &BuyRequest{
		CategoryId:   "63035f03d723a30001c3e08e",
		DistributeId: "63035f03d723a30001c3e08c",
		Size:         "4000",
	}
	buyBody, err = json.Marshal(*b)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start()
}

func start() {
	b, err := os.ReadFile("./accounts.txt")
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("\r?\n")

	accs := re.Split(string(b), 30)

	for {
		if time.Now().Unix() == START {
			break
		}
	}

	for i := 0; i < len(accs); i++ {
		a := accs[i]
		data := strings.Split(a, "*")
		n := i

		if len(data) == 0 {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			acc := &Account{}
			acc.InitAcc(data[0], data[1], data[2])
			for {
				acc.Validate()
				ok := acc.BuyNFT(n)
				if ok {
					break
				}
			}
		}()

	}
	wg.Wait()
}
