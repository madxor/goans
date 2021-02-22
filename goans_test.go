package goans

import "testing"
import "bytes"
import "math/rand"

func TestEncodeDecode(t *testing.T) {
	buf := []byte("testtesttest")

	cfg := GetConfigurationFromSample(buf)

	e := EncodeFrame(buf, 0, cfg)
	m := DecodeFrame(e, cfg)

	if bytes.Compare(m, buf) != 0 {
		t.Errorf("%s != %s", buf, m)
	}
}

func TestEncodeDecodeGeometric(t *testing.T) {
	rand.Seed(int64(1))

	var buf = []byte{ byte(rand.Intn(256)) }

	for i := 1; i < 100 ; i++ {
		buf = append(buf, byte(rand.Intn(256)))
	}

	for p := 0.5; p <= 0.9; p += 0.1 {
		cfg := GetConfigurationFromGeometricDistribution(p, 256, 100)
		e := EncodeFrame(buf, 1, cfg)
		m := DecodeFrame(e, cfg)

		if bytes.Compare(m, buf) != 0 {
			t.Errorf("%s != %s", buf, m)
			return
		}
	}
}

func TestEncodeDecodeWithRedundancy(t *testing.T) {
	buf1 := []byte("teeeeeeeeeesttest")
	buf2 := []byte("eeeeetesteeeeeeee")

	cfg := GetConfigurationFromSample(buf1)

	e1 := EncodeFrame(buf1, 0, cfg)
	m1 := DecodeFrame(e1, cfg)

	if bytes.Compare(m1, buf1) != 0 {
		t.Errorf("%s != %s", buf1, m1)
		return
	}

	e2 := EncodeFrame(buf2, 0, cfg)
	m2 := DecodeFrame(e2, cfg)

	if bytes.Compare(m2, buf2) != 0 {
		t.Errorf("%s != %s", buf2, m2)
	}
}

func TestInitialStateIterateRandom(t *testing.T) {
	cfg := GetConfigurationFromGeometricDistribution(0.5, 256, 100)
	buf := GenerateRandomFrame(cfg)

	for i := 0; i < (2 << cfg.R + 1); i++ {
		e := EncodeFrame(buf, i, cfg)
		m := DecodeFrame(e, cfg)

		if bytes.Compare(m, buf) != 0 {
			t.Errorf("%d: %s != %s", i, buf, m)
			return
		}
	}
}

