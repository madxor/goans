package goans

import "math"
import "fmt"
import "strconv"
import "sort"
import "math/rand"

type StackState struct {
	v float64
	s byte
}

type SortStack []StackState

type State struct {
	s byte
	x int
}

type Configuration struct {
	R int
	F int // Frame size
	P map[byte]float64
	L map[byte]int
	D map[int]State
	C map[State]int
}

type EncodedFrame struct {
	F int
	B []byte
	//B string
}

type Decoded []byte

func (a SortStack) Len() int { return len(a) }
func (a SortStack) Less(i, j int) bool { return a[i].v < a[j].v }
func (a SortStack) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func Configure(R int, L map[byte]int, D map[int]State, C map[State]int) Configuration {
	var cfg Configuration
	// this should be calculated after calculating P
	if (R <= int(math.Ceil(math.Log2(float64(len(cfg.P)))))) {
		cfg.R = int(math.Ceil(math.Log2(float64(len(cfg.P))))) + 1
	}
	return cfg
}

func CalculateProbabilitiesFromGeometricDistribution(G float64, N int) map[byte]float64 {
	var PSum float64
	P := make(map[byte]float64)

	if N > 256 || N <= 0 {
		return nil
	}

	if G >= 1.0 || G < 0.5 {
		return nil
	}

	for i := 0; i < N; i++ {
		p := G * math.Pow(float64(1) - G, float64(i + 1)) // if i==0 then this is broken, therefore the +1 here
		P[byte(i)] =  p
		PSum += p
		if i + 2 == N {
			P[byte(i + 1)] = float64(1) - PSum
			break
		}
	}
	return P
}

func CalculateProbabilitiesFromSample(sample []byte) map[byte]float64 {
	P := make(map[byte]float64)

	for i := 0; i < len(sample); i++ {
		P[sample[i]] += float64(1) / float64(len(sample))
	}

	return P
}

func CalculateL(cfg Configuration) map[byte]int {
	L := make(map[byte]int)

	eL := math.Pow(2, float64(cfg.R))

	for k, v := range cfg.P {
		L[k] = int(math.Ceil(float64(eL) * float64(v)))
	}

	// Correct L for non-trivial cases

	var sMax byte // symbol with highest representation
	rMax := 0  // representation of the above symbol

	LMax := int(eL)
	sign := 1

	for ;; {
		LSum := 0
		for i := range L {
			LSum += L[i]
			if L[i] > rMax {
				sMax = i
				rMax = L[i]
			}
		}

		if LSum == LMax {
			break
		}

		if LSum > LMax {
			sign = -1
			// This trick should be true for geometric distribution
			if rMax > ((LSum - LMax) * 2) {
				L[sMax] -= LSum - LMax
				break
			}
		}

		for i := range L {
			if i != sMax && L[i] > 1 {
				L[i] += sign
				LSum += sign
			}
			if LSum == LMax {
				break
			}
		}
	}

	// Double check if there are no symbols with zero representation.
	// If found reduce the most represented symbol
	for i := range L {
		for j := range L {
			if L[j] > rMax {
				sMax = j
				rMax = L[j]
			}
		}
		if i != sMax && L[sMax] > 1 && L[i] == 0 {
			L[i]++
			L[sMax]--
		}
	}

	return L
}

func CalculateTables(cfg Configuration) (map[State]int, map[int]State) {
	C := make(map[State]int)
	D := make(map[int]State)

	eL := int(math.Pow(2, float64(cfg.R)))

	// Recalculate probabilities based on L values
	P := make(map[byte]float64)
	for i := range cfg.L {
		P[i] = float64(cfg.L[i]) / float64(eL)
	}

	X := make(map[byte]int)

	var stack []StackState

	for s, _ := range P {
		t := StackState{0.5/P[s], s}
		stack = append(stack, t)
		X[s] = cfg.L[s]
	}

	for x := eL; x < eL * 2; x++ {
		sort.Sort(SortStack(stack))
		t := stack[0]
		stack = stack[1:]
		stack = append(stack, StackState{t.v + 1 / P[t.s], t.s})

		D[x] = State{t.s, X[t.s]}
		C[State{t.s, X[t.s]}] = x
		X[t.s]++
	}

	return C, D
}

func GetConfigurationFromGeometricDistribution(G float64, N int, F int) Configuration {
	var cfg Configuration

	cfg.F = F
	cfg.P = CalculateProbabilitiesFromGeometricDistribution(G, N)
	cfg.R = int(math.Ceil(math.Log2(float64(len(cfg.P))))) + 1
	cfg.L = CalculateL(cfg)
	cfg.C, cfg.D = CalculateTables(cfg)

	return cfg
}

func GetConfigurationFromSample(sample []byte) Configuration {
	var cfg Configuration

	cfg.F = len(sample)
	cfg.P = CalculateProbabilitiesFromSample(sample)
	cfg.R = int(math.Ceil(math.Log2(float64(len(cfg.P))))) + 1
	cfg.L = CalculateL(cfg)
	cfg.C, cfg.D = CalculateTables(cfg)

	return cfg
}

func C(s byte, x, k int, cfg Configuration) int {
	return int(cfg.C[State{s, int(math.Floor(float64(x) / math.Pow(2, float64(k))))}])
}

func EncodeFrame(m []byte, X int, cfg Configuration) EncodedFrame {
	var e EncodedFrame
	var B string
	var M = []byte{ 0 }

	step := 8

	eL := math.Pow(2, float64(cfg.R))

	x := int(eL + math.Mod(float64(X), eL))

	for i := len(m) - 1; i >= 0; i-- {
		s := m[i]
		k := int(math.Floor(math.Log2(float64(x) / float64(cfg.L[s]))))
		if k > 0 && C(s, x, k, cfg) > 0 {
			// TODO: Use bit operations instead of string
			t := "%0" + strconv.FormatUint(uint64(k), 10) + "b"
			b := fmt.Sprintf(t, int(math.Mod(float64(x), math.Pow(2, float64(k)))))
			B = b + B
		}
		x = C(s, x, k, cfg)
	}

	// Add zeros to the end in order to fill full bytes
	for i := 0; i < len(B) % 8; i++ {
		B += "0"
	}

	T := []byte(B)

	//fmt.Println("--------")
	//fmt.Println(m)
	for i := 0; i < int(math.Ceil(float64(len(B)) / 8)); i++ {

		t := 0
		for j := 0; j < step; j++ {
			t += int(T[i * 8 + j] - '0') << (7 - j)
		}
/*
		t, err := strconv.ParseUint(B[i * 8 : i * 8 + step], 2, 32)
	//	fmt.Println("EncodeFrame: t = ", t, "B[", i * 8, ":", i * 8 + step, "]", "len(B)", len(B))

		if err != nil {
			fmt.Println("EncodeFrame Error:", err)
			break
		}

	//	fmt.Println("EncodeFrame: t = ", t, "; err ->", err)
*/
		if i == 0 {
			M[0] = byte(t)
		} else {
			M = append(M, byte(t))
		}
	}
	//fmt.Println(M)
	//fmt.Println(B)
	//fmt.Println("--------")

	e.F = x
	e.B = M
	return e
}

func DecodeFrame(e EncodedFrame, cfg Configuration) []byte {
	var m, b, B []byte
	var T string

	for i := 0; i < len(e.B); i++ {
		T += fmt.Sprintf("%08b", e.B[i])
	}
	x := e.F
	B = []byte(T)

	//fmt.Println("--------")
	//fmt.Println(e.B)
	//fmt.Println(B)
	//fmt.Println(T)

	for i := 0; i < cfg.F; i++ {
		bb := int64(0)
		m = append(m, cfg.D[x].s)
		if i == cfg.F - 1 { // break after the last symbol
			break
		}
		//fmt.Println("DecodeFrame:", i, " m = ", m)
		k := int(float64(cfg.R) - math.Floor(math.Log2(float64(cfg.D[x].x))))
		//fmt.Println("DecodeFrame:", "k = ", k, " B = ", B, "len(B) =", len(B))
		if len(B) > 0 {
			b = B[:k]
			B = B[k:]
			bb, _ = strconv.ParseInt(string(b), 2, 64)
		}
		//fmt.Println("DecodeFrame:", i, " bb = ", bb, "len(B) =", len(B))
		x = int(math.Pow(2, float64(k)) * float64(cfg.D[x].x) + float64(int(bb)))
	}
	return m
}

func GenerateRandomFrame(cfg Configuration) []byte {
	rand.Seed(int64(1))
	eL := int(math.Pow(2, float64(cfg.R)))

	var frame = []byte { byte((cfg.D[eL + rand.Intn(eL)]).s) }

	for i := 1; i < cfg.F; i++ {
		frame = append(frame, byte((cfg.D[eL + rand.Intn(eL)]).s))
	}

	return frame
}

