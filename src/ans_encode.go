package main

import "fmt"
import "math"
import "encoding/json" 
import "flag"
import "io/ioutil"
import "strconv"
import str "strings"
import b64 "encoding/base64"

/*
	Input:
	C <-- encoding table
	Ls <-- set of symbol L values
	M <-- message to encode
	x <-- initial state

	Output:
	E <-- encoded message
	X <-- final state
 */

type State struct {
	v float64
	s string
}

type CodingState struct {
	S string
	X int
}

type EncodedMessage struct{
	M string
	F int
}

type ANSConfiguration struct {
	R	int
	Ls	map[string]int
	D	map[int]CodingState
	//C	map[CodingState]int
	CJ	map[string]int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	prefixPtr := flag.String("prefix", "test", "prefix for a configuration file")
	messagePtr := flag.String("m", "abc", "message to be encoded")
	XPtr := flag.Int("X", 12, "initial state")
	dbgPtr := flag.Bool("debug", false, "debugging")

	flag.Parse()
	
	var A ANSConfiguration
	
	Af, err := ioutil.ReadFile(*prefixPtr + "_config.json")
	check(err)
	
	check(json.Unmarshal(Af, &A))
	
	M := []byte(*messagePtr)
	
	if *dbgPtr {
		fmt.Println(M)
	}
	
	var s string
	var E string
	var b []int
	var X int
	var x float64
	var stateStart int
	
	stateStart = 9999999
	
	C := make(map[CodingState]int)

	for k, v := range A.CJ {
		if *dbgPtr { fmt.Println(k) }
		sTmp := str.Split(k,"+")
		xTmp, _ := b64.StdEncoding.DecodeString(sTmp[1])
		C[CodingState{sTmp[0], int(xTmp[0])}] = v
		if v < stateStart {
			stateStart = v
		}
		if *dbgPtr { fmt.Println(sTmp, xTmp, v) }
	}

	if *dbgPtr {
		fmt.Println("Encoding table recreated from file ->", C)
		fmt.Println("State start ->", stateStart)
	}

	//x = float64(stateStart) + float64(*XPtr) // set the initial state
	
	x = float64(*XPtr) // set the initial state

	for i := len(M); i > 0; i-- {
		s = string(M[i-1])
		
		if *dbgPtr {
			fmt.Println("Encoding ->", s)
		}

		k := math.Floor(math.Log2(x / float64(A.Ls[s])))
		if *dbgPtr { fmt.Println("k = floor(log2(x/Ls) \t-->\t", k, "= floor(log2(",x, "/", A.Ls[s], "))") }
		if k > 0 {
			b = append([]int{int(math.Mod(x, math.Pow(2, k)))}, b...)
			if *dbgPtr { fmt.Println(i, "->", b) }
		}
		x = float64(C[CodingState{s, int(math.Floor(x / math.Pow(2, k)))}])
	}
	X = int(x)

	fmt.Println("Final state ->", X)
	
	if *dbgPtr { fmt.Println("Encoded proto-string ->", b) }
	
	for i := 0; i < len(b); i++ {
		E = E + strconv.FormatInt(int64(b[i]), 2)
	}
	fmt.Println("Encoded bitstring ->", E)
	
	O := EncodedMessage{E, X}
	
	OJ, err := json.MarshalIndent(O, "", "        ")
	check(err)

	if *dbgPtr {
		fmt.Println("Encoded message ->", O)
	}

	check(ioutil.WriteFile(*prefixPtr + "_" + *messagePtr + "_encoded.json", OJ, 0644))
}
