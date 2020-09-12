package main

import "fmt"
import "os"
import "encoding/json"
import "flag"
import "bufio"
import "io/ioutil"
import "strconv"
import "math"

/*
	Generator reads from standard input and generates a json file
	with parameters required by the initialization initialization
	procedure.
 */

type Symbols struct {
	S string
	P float64
}

type InitializationParameters struct {
	R int
	S []Symbols
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	prefixPtr := flag.String("prefix", "test", "prefix name for parameters output files")
	RPtr	:= flag.Int("R", 3, "encoding quality parameter")
	NPtr	:= flag.Int("N", 0, "number of symbols to read")
	GPtr	:= flag.Float64("G", 0.0, "parameter for geometric distribution of symbols")
	dbgPtr := flag.Bool("debug", false, "debugging")

	flag.Parse()

	if *GPtr >= 1.0 {
		fmt.Printf("G parameter must be smaller than 1")
	}

	if *GPtr < 0.0 {
		fmt.Printf("G parameter must be larger than 0")
	}

	if *NPtr > *RPtr {
		RPtr = NPtr
	}

	var S []Symbols

	if *NPtr == 0 {
		S = []Symbols{{"A", 0.3}, {"B", 0.2}, {"C", 0.4}, {"D", 0.1}}
	} else {
		for i := 'A'; int(i) < int('A') + *NPtr; i++ {
			if *GPtr == 0.0 {
				fmt.Printf("Enter  probability for symbol %c: ", rune(i))
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				in := scanner.Text()
				p, _ := strconv.ParseFloat(in, 64)
				if *dbgPtr { fmt.Println(p) }
				S = append(S, Symbols{string(i), p})
			} else {
				// Use gemoetric distribution of symbols
				S = append(S, Symbols{string(i), math.Pow(*GPtr, float64(1 + int(i) - int('A')))})
				if int(i) - int('A') + 2 == *NPtr  {
					S = append(S, Symbols{string(i + 1), math.Pow(*GPtr, float64(1 + int(i) - int('A')))})
					break
				}
			}
		}
	}

	var IP InitializationParameters

	IP.R = int(*RPtr)
	IP.S = S

	IPj, _ := json.MarshalIndent(IP, "", "        ")
	if *dbgPtr { fmt.Println("Json IP ->", string(IPj)) }
	check(ioutil.WriteFile(*prefixPtr + "_parameters.json", IPj, 0644))
}
