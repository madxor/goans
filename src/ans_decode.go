package main

import "fmt"
import "math"
import "encoding/json" 
import "flag"
import "io/ioutil"
import "strconv"

/*
	Input:
	D <-- encoding table
	Ls <-- set of symbol L values
	B <-- bitstream to decode
	X <-- final state

	Output:
	S <-- decoded message
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
	pPtr := flag.String("prefix", "test", "prefix of a configuration file")
	ePtr := flag.String("encoded", "abc", "original message to be decoded")
	//mPtr := flag.String("m", "1000", "message to be encoded")
	//XPtr := flag.Int("X", 9, "final state")
	dPtr := flag.Bool("debug", false, "debugging")

	flag.Parse()

	var A ANSConfiguration
	var E EncodedMessage
	
	Af, err := ioutil.ReadFile(*pPtr + "_config.json")
	check(err)
	
	check(json.Unmarshal(Af, &A))
	
	Ef, err := ioutil.ReadFile(*pPtr + "_" + *ePtr + "_encoded.json")
	check(err)
	
	check(json.Unmarshal(Ef, &E))

	var b []byte

	B := []byte(E.M)

	X := E.F // set the final state
	
	if *dPtr {
		fmt.Println(B)
	}
	
	var S []string

	for ;; {
		S = append(S, A.D[X].S)		// si++ = D[x]

		if *dPtr {
			fmt.Println(len(S)) 
			fmt.Println("\tS = D[X] \t\t\t-->\t", A.D[X].S, " = D[", X, "]")
		}

		k := int(float64(A.R) - math.Floor(math.Log2(float64(A.D[X].X))))

		if *dPtr {
			fmt.Println("\tk := R - Floor(Log2(D[X].X)) \t-->\t", k, ":=", A.R, "-", "Floor(Log2(", A.D[X].X, "))")
			fmt.Println("\tpre B ->", B)
		}
		if len(B) < k {
			b = B
			B = nil
			bb, _ := strconv.ParseInt(string(b), 2, 64)
			X = int(float64(math.Pow(2, float64(k))) * float64(A.D[X].X) + float64(int(bb)))
			S = append(S, A.D[X].S)		// si++ = D[x]

			if *dPtr {
				fmt.Println(len(S))
				fmt.Println("\tS = D[X] \t\t\t-->\t", A.D[X].S, " = D[", X, "]")
				fmt.Println("\tX = Pow(2, k)) * D[X].X + bb \t-->\t", X, "= Pow(2, ",k ,")) * ", A.D[X].X, "+", bb)
				fmt.Println("\tpost B ->", B)
				fmt.Println("\tb ->", b)
				fmt.Println("\tbb ->", bb)
			}
		} else {
			b = B[:k]
			B = B[k:]
		
			bb, _ := strconv.ParseInt(string(b), 2, 64)
			X = int(float64(math.Pow(2, float64(k))) * float64(A.D[X].X) + float64(int(bb)))

			if *dPtr {
				fmt.Println("\tX = Pow(2, k)) * D[X].X + bb \t-->\t", X, "= Pow(2, ",k ,")) * ", A.D[X].X, "+", bb)
				fmt.Println("\tpost B ->", B)
				fmt.Println("\tb ->", b)
				fmt.Println("\tbb ->", bb)
			}
		}

		if len(B) == 0 {
			break
		}
	}

	fmt.Println("Decoded string ->", S)
}
