package sdr

import (
	"math"
	"testing"
)

const approxErrorLimit = 0.011

var (
	atanBenchTable      = [][2]float32{}
	atanBenchTableFixed = [][2]int{}
)

func init() {
	for y := -1.0; y <= 1.0; y += 0.5 {
		for x := -1.0; x <= 1.0; x += 0.5 {
			atanBenchTable = append(atanBenchTable, [2]float32{float32(x), float32(y)})
			atanBenchTableFixed = append(atanBenchTableFixed, [2]int{int(x * (1 << 14)), int(y * (1 << 14))})
		}
	}
}

func TestAtan2(t *testing.T) {
	for y := -1.0; y <= 1.0; y += 0.01 {
		for x := -1.0; x <= 1.0; x += 0.01 {
			expected := float32(math.Atan2(y, x))
			if err := math.Abs(float64(expected - FastAtan2(float32(y), float32(x)))); err > approxErrorLimit {
				t.Errorf("FastAtan2 gave an error of %f for x=%f y=%f", err, x, y)
			}
			if err := math.Abs(float64(expected - FastAtan2_2(float32(y), float32(x)))); err > approxErrorLimit {
				t.Errorf("FastAtan2_2 gave an error of %f for x=%f y=%f", err, x, y)
			}
		}
	}
	x, y := 0.0, 0.0
	expected := float32(math.Atan2(y, x))
	if err := math.Abs(float64(expected - FastAtan2(float32(y), float32(x)))); err > approxErrorLimit {
		t.Errorf("FastAtan2 gave an error of %f for x=%f y=%f", err, x, y)
	}
	if err := math.Abs(float64(expected - FastAtan2_2(float32(y), float32(x)))); err > approxErrorLimit {
		t.Errorf("FastAtan2_2 gave an error of %f for x=%f y=%f", err, x, y)
	}
}

func TestFastAtan2Error(t *testing.T) {
	maxE := 0.0
	sumE := 0.0
	count := 0
	for y := -1.0; y <= 1.0; y += 0.01 {
		for x := -1.0; x <= 1.0; x += 0.01 {
			ai := float64(FastAtan2(float32(y), float32(x)))
			af := math.Atan2(y, x)
			e := math.Abs(ai - af)
			sumE += e
			if e > maxE {
				maxE = e
			}
			count++
		}
	}
	if maxE > 0.0102 {
		t.Errorf("Expected max error of 0.0102 got %f", maxE)
	}
	t.Logf("Max error %f\n", maxE)
	t.Logf("Mean absolute error %f", sumE/float64(count))
}

func TestFastAtan2_2Error(t *testing.T) {
	maxE := 0.0
	sumE := 0.0
	count := 0
	for y := -1.0; y <= 1.0; y += 0.01 {
		for x := -1.0; x <= 1.0; x += 0.01 {
			ai := float64(FastAtan2_2(float32(y), float32(x)))
			af := math.Atan2(y, x)
			e := math.Abs(ai - af)
			sumE += e
			if e > maxE {
				maxE = e
			}
			count++
		}
	}
	if maxE > 0.005 {
		t.Errorf("Expected max error of 0.005 got %f", maxE)
	}
	t.Logf("Max error %f\n", maxE)
	t.Logf("Mean absolute error %f", sumE/float64(count))
}

func TestScaleF32(t *testing.T) {
	input := make([]float32, 257)
	for i := 0; i < len(input); i++ {
		input[i] = float32(i)
	}
	expected := make([]float32, len(input))
	output := make([]float32, len(input))
	scalef32(input, expected, 1.0/256.0)
	Scalef32(input, output, 1.0/256.0)
	for i, v := range expected {
		if output[i] != v {
			t.Fatalf("Output doesn't match expected:\n%+v\n%+v", output, expected)
		}
	}

	// Unaligned
	input = input[1:]
	expected = make([]float32, len(input)+1)[1:]
	output = make([]float32, len(input)+1)[1:]
	scalef32(input, expected, 1.0/256.0)
	Scalef32(input, output, 1.0/256.0)
	for i, v := range expected {
		if output[i] != v {
			t.Fatalf("Output doesn't match expected:\n%+v\n%+v", output, expected)
		}
	}
}

func BenchmarkConj32(b *testing.B) {
	in := complex64(complex(1.0, -0.2))
	for i := 0; i < b.N; i++ {
		_ = Conj32(in)
	}
}

func BenchmarkFastAtan2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, xy := range atanBenchTable {
			FastAtan2(xy[1], xy[0])
		}
	}
}

func BenchmarkFastAtan2_Go(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, xy := range atanBenchTable {
			fastAtan2(xy[1], xy[0])
		}
	}
}

func BenchmarkFastAtan2_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, xy := range atanBenchTable {
			FastAtan2_2(xy[1], xy[0])
		}
	}
}

func BenchmarkFastAtan2_2_Go(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, xy := range atanBenchTable {
			fastAtan2_2(xy[1], xy[0])
		}
	}
}

func BenchmarkAtan2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, xy := range atanBenchTable {
			math.Atan2(float64(xy[1]), float64(xy[0]))
		}
	}
}

func BenchmarkScalef32(b *testing.B) {
	input := make([]float32, benchSize)
	output := make([]float32, len(input))
	b.SetBytes(benchSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scalef32(input, output, 1.0/benchSize)
	}
}

func BenchmarkScalef32_Go(b *testing.B) {
	input := make([]float32, benchSize)
	output := make([]float32, len(input))
	b.SetBytes(benchSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scalef32(input, output, 1.0/benchSize)
	}
}
