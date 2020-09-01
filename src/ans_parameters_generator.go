package main

import "fmt"
import "os"
import "encoding/json" 
import "flag"
import "bufio"
import "io/ioutil"
import "strconv"

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
	dbgPtr := flag.Bool("debug", false, "debugging")

	flag.Parse()
	
	var S []Symbols

	if *NPtr == 0 {
		S = []Symbols{{"a", 0.3}, {"b", 0.2}, {"c", 0.4}, {"d", 0.1}}
	} else {
		for i := 'A'; int(i) < int('A') + *NPtr; i++ {
			fmt.Printf("Enter  probability for symbol %c: ", rune(i))
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			in := scanner.Text()
			p, _ := strconv.ParseFloat(in, 64)
			if *dbgPtr { fmt.Println(p) }
			S = append(S, Symbols{string(i), p})
		}
	}

	var IP InitializationParameters

	IP.R = int(*RPtr)
	IP.S = S

	IPj, _ := json.MarshalIndent(IP, "", "        ")
	if *dbgPtr { fmt.Println("Json IP ->", string(IPj)) }
	check(ioutil.WriteFile(*prefixPtr + "_parameters.json", IPj, 0644))
}
