package main

import "fmt"
import "math"
import "encoding/json"
import "flag"
import "sort"
//import "strconv"
import "io/ioutil"
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

type SymbolEncVariance struct {
	MinK int
	MaxK int
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

func depthSearch(sLen int, bLen int, s int, S []string, V map[string]SymbolEncVariance, sS []string, sK []int, dbg bool) (bool) {
	if (dbg) {fmt.Println("depthSearch -> ", sLen, bLen, s, S, V)}

	var found bool

	found = false

	if sLen == 0 {
		if bLen == 0 {
			if (dbg) {fmt.Print("Found for:", sLen, bLen, s, S, V)}
			fmt.Println("sS, sK ->", sS, sK)
			return true
		} else {
			/* do nothing */
			return false
		}
	}

	for i := s; i < len(S); i++ {
		if (dbg) {fmt.Println(S[i])}
		sS := append(sS, S[i])
		if bLen >= 0 {
			for k := V[S[i]].MinK; k <= V[S[i]].MaxK; k++ {
				if (dbg) {fmt.Println("k ->", k, V[S[i]].MinK, V[S[i]].MaxK)}
				sK := append(sK, k)
				if depthSearch(sLen - 1, bLen - k, i, S, V, sS, sK, dbg) {
					if (dbg) {fmt.Println(s, "Matched for ->", S[i], k)}
					found = true
				} else {
					found = false
				}
			}
		}
	}
	if (dbg) {fmt.Println()}
	return found
}

func main() {
	prefixPtr := flag.String("prefix", "test", "prefix for a configuration file")
	dPtr := flag.Bool("debug", false, "debugging")
	sPtr := flag.Int("s", 0, "length of symbols")
	bPtr := flag.Int("b", 0, "length of binary output")

	flag.Parse()

	var A ANSConfiguration

	Af, err := ioutil.ReadFile(*prefixPtr + "_config.json")
	check(err)

	check(json.Unmarshal(Af, &A))

	var stateStart int
	var stateCounter int

	stateStart = 9999999
	stateCounter = 0

	C := make(map[CodingState]int)
	V := make(map[string]SymbolEncVariance)

	for k, v := range A.CJ {
		if *dPtr { fmt.Println(k) }
		sTmp := str.Split(k,"+")
		xTmp, _ := b64.StdEncoding.DecodeString(sTmp[1])
		C[CodingState{sTmp[0], int(xTmp[0])}] = v
		if v < stateStart {
			stateStart = v
		}
		if *dPtr { fmt.Println(sTmp, xTmp, v) }
		stateCounter++
	}

	if *dPtr {
		fmt.Println("Encoding table recreated from file ->", C)
	}

	for X := range A.D {
		k := int(float64(A.R) - math.Floor(math.Log2(float64(A.D[X].X))))

		if *dPtr {
			fmt.Println("\tk := R - Floor(Log2(D[X].X)) \t-->\t", k, ":=", A.R, "-", "Floor(Log2(", A.D[X].X, "))")
		}

		if k <= V[A.D[X].S].MinK || V[A.D[X].S].MaxK == 0 {
			V[A.D[X].S] = SymbolEncVariance{k, V[A.D[X].S].MaxK}
		}

		if k >= V[A.D[X].S].MaxK {
			V[A.D[X].S] = SymbolEncVariance{V[A.D[X].S].MinK, k}
		}

		if *dPtr {
			fmt.Println(A.D[X].S, "->", V[A.D[X].S].MinK, V[A.D[X].S].MaxK)
		}
	}

	symbols := make([]string, 0, len(V))
	for s := range V {
		symbols = append(symbols, s)
		if *dPtr {
			fmt.Println(s)
			fmt.Println(symbols)
		}
	}
	sort.Strings(symbols)

	fmt.Println(V)

	if *dPtr {
		fmt.Println(symbols[0], symbols[len(symbols) - 1])
	}

	var solutionS []string
	var solutionK []int

	depthSearch(*sPtr, *bPtr, 0, symbols, V, solutionS, solutionK, *dPtr)
}
