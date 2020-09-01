package main

import "fmt"
import "sort"
import "math"
import "encoding/json" 
import "flag"
import "io/ioutil"
import b64 "encoding/base64"

/*
	Input:
	S <-- set of symols
	P <-- symbols probability distribution
	R <-- encoding quality parameter

	Output:
	C <-- encoding table
	D <-- decoding table
 */

type Symbols struct {
	S string
	P float64
}

type State struct {
	v float64
	s string
}

type CodingState struct {
	S string
	X int
}

type InitializationParameters struct {
	R int
	S []Symbols
}

type ANSConfiguration struct {
	R	int
	Ls	map[string]int
	D	map[int]CodingState
	//C	map[CodingState]int
	CJ	map[string]int
}

type SortState []State

func (a SortState) Len() int 		{ return len(a) }
func (a SortState) Less(i, j int) bool	{ return a[i].v < a[j].v }
func (a SortState) Swap(i, j int)	{ a[i], a[j] = a[j], a[i] }

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	prefixPtr := flag.String("prefix", "test", "prefix for parameters and configuration files")
	dbgPtr := flag.Bool("debug", false, "debugging")

	flag.Parse()

	var IP InitializationParameters

	var S []Symbols
	var R float64

	P := make(map[string]float64)
	Ls := make(map[string]int)
	Xs := make(map[string]int)
	D := make(map[int]CodingState)
	C := make(map[CodingState]int)
	C2j := make(map[string]int)

	var stateStack []State
	var t State

	IPFile, err := ioutil.ReadFile(*prefixPtr + "_parameters.json")
	check(err)
	if *dbgPtr {
		fmt.Println("Json IP file ->", string(IPFile))
	}
	
	check(json.Unmarshal(IPFile, &IP))

	if *dbgPtr {
		fmt.Println(IP)
	}

	R = float64(IP.R)
	S = IP.S

	for i := 0; i < len(S); i++ {
		P[S[i].S] = S[i].P
		if *dbgPtr {
			fmt.Println(i, S[i].S, P[S[i].S])
		}
	}

	var L = int(math.Pow(2, R))

	if *dbgPtr {
		fmt.Println(L)
	}

	for i := 0; i < len(S); i++ {
		Ls[S[i].S] = int(float64(L) * P[S[i].S]) //+ 1
		if Ls[S[i].S] == 0 {
			Ls[S[i].S]++
		}
		if *dbgPtr {
			fmt.Println(i, S[i].S, P[S[i].S], Ls[S[i].S])
		}
	}

	for i := 0; i < len(S); i++ {
		t = State{ 0.5/P[S[i].S], S[i].S }
		stateStack = append(stateStack, t)
		Xs[S[i].S] = Ls[S[i].S]

		if *dbgPtr {
			fmt.Println(i, Xs[S[i].S], t)
		}
	}

	for x := int(L); x < int(L)*2; x++ {
		if *dbgPtr {
			fmt.Println(stateStack)
		}

		sort.Sort(SortState(stateStack))

		if *dbgPtr {
			fmt.Println(stateStack)
		}

		t, stateStack = stateStack[0], stateStack[1:]		// (v,s) = getmin()
		stateStack = append(stateStack, State{t.v + 1/P[t.s], t.s})	// put((v+1/ps, s))


		D[x] = CodingState{t.s, Xs[t.s]}				// D[x] = (s, xs)
		C[CodingState{t.s, Xs[t.s]}] = x				// C[xs, s] = x
		tmp := b64.StdEncoding.EncodeToString([]byte(string(Xs[t.s])))
		//C2j[t.s + "-" + string(Xs[t.s])] = x				// Coding table json writer helper
		C2j[t.s + "+" + tmp] = x				// Coding table json writer helper
		Xs[t.s]++							// xs ++
	}

	A := ANSConfiguration{IP.R, Ls, D, C2j}
	
	AJ, err := json.MarshalIndent(A, "", "        ")
	check(err)

	if *dbgPtr {
		fmt.Println("ANS Configuration ->", string(AJ))
	}

	check(ioutil.WriteFile(*prefixPtr + "_config.json", AJ, 0644))

	fmt.Println("ANS initialization done")
}

