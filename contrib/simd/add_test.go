package simd

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchClever(t *testing.T) {
	Search := Clever
	keys := make([]uint64, 512)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}

	for i := 0; i < len(keys); i++ {
		idx := int(Search(keys, uint64(i)))
		require.Equal(t, (i+1)/2, idx, "%v\n%v", i, keys)
	}
	require.Equal(t, 256, int(Search(keys, math.MaxUint64>>1)))
	require.Equal(t, 256, int(Search(keys, math.MaxUint64)))
}

func TestSearchNaive(t *testing.T) {
	Search := Naive
	keys := make([]uint64, 512)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}

	for i := 0; i < len(keys); i++ {
		idx := int(Search(keys, uint64(i)))
		require.Equal(t, (i+1)/2, idx, "%v\n%v", i, keys)
	}
	require.Equal(t, 256, int(Search(keys, math.MaxUint64>>1)))
	require.Equal(t, 256, int(Search(keys, math.MaxUint64)))
}

func TestSearchSIMD(t *testing.T) {
	Search := skernel
	keys := make([]uint64, 512)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}

	for i := 0; i < len(keys); i++ {
		idx := int(Search(keys, uint64(i)))
		require.Equal(t, (i+1)/2, idx, "%v\n%v", i, keys)
	}
	require.Equal(t, 256, int(Search(keys, math.MaxInt64>>1)))
	require.Equal(t, 256, int(Search(keys, math.MaxInt64)))
}

func TestSearchParallel(t *testing.T) {
	Search := Parallel
	keys := make([]uint64, 512)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}

	for i := 0; i < len(keys); i++ {
		idx := int(Search(keys, uint64(i)))
		require.Equal(t, (i+1)/2, idx, "%v\n%v", i, keys)
	}
	require.Equal(t, 256, int(Search(keys, math.MaxInt64>>1)))
	require.Equal(t, 256, int(Search(keys, math.MaxInt64)))
}

func TestSIMDKernel(t *testing.T) {
	data := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	out0 := skernel(data, 0)
	out1 := skernel(data, 1)
	out2 := skernel(data, 2)
	out7 := skernel(data, 7)
	out10 := skernel(data, 10)
	out50 := skernel(data, 50)
	t.Logf("out %v %v %v %v %v %v", out0, out1, out2, out7, out10, out50)
}
func TestNaive(t *testing.T) {
	data := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	out0 := Naive(data, 0)
	out1 := Naive(data, 1)
	out2 := Naive(data, 2)
	out7 := Naive(data, 7)
	out10 := Naive(data, 10)
	out50 := Naive(data, 50)
	t.Logf("out %v %v %v %v %v %v", out0, out1, out2, out7, out10, out50)
}

func Test_cmp2(t *testing.T) {
	data := [2]uint64{0, 1}
	pk := [2]uint64{0, 0}
	for i := range data {
		// fill pk with i
		for j := range pk {
			pk[j] = uint64(i)
		}
		s := cmp2(data, pk)
		s_n := cmp2_native(data, pk)
		require.Equal(t, i, int(s))
		require.Equal(t, s_n, s)
	}
}

func Test_cmp4(t *testing.T) {
	data := [4]uint64{0, 1, 2, 3}
	pk := [4]uint64{0, 0, 0, 0}
	for i := range data {
		// fill pk with i
		for j := range pk {
			pk[j] = uint64(i)
		}
		s := cmp4(data, pk)
		s_n := cmp4_native(data, pk)
		require.Equal(t, i, int(s))
		require.Equal(t, s_n, s)
	}
}

func Test_cmp8(t *testing.T) {
	data := [8]uint64{0, 1, 2, 3, 4, 5, 6, 7}
	pk := [4]uint64{0, 0, 0, 0}
	for i := range data {
		// fill pk with i
		for j := range pk {
			pk[j] = uint64(i)
		}
		s := cmp8(data, pk)
		s_n := cmp8_native(data, pk)
		require.Equal(t, i, int(s))
		require.Equal(t, s_n, s)
	}
	/* keys that are greater than maxint64
	var n1, n2 int64 = -1, -2
	data[1] = *(*uint64)(unsafe.Pointer(&n1))
	data[2] = *(*uint64)(unsafe.Pointer(&n2))
	for i := range data {
		// fill pk with i
		for j := range pk {
			pk[j] = uint64(i)
		}
		s := cmp8(data, pk)
		s_n := cmp8_native(data, pk)
		require.Equal(t, i, int(s))
		require.Equal(t, s_n, s)
	}
	*/
}

func Benchmark_cmp2_native(b *testing.B) {
	b.StopTimer()
	twos := [2]uint64{2, 1}
	pks := [2]uint64{2, 2}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp2_native(twos, pks)
	}
	_ = idx
}

func Benchmark_cmp2_sse(b *testing.B) {
	b.StopTimer()
	twos := [2]uint64{2, 1}
	pks := [2]uint64{2, 2}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp2(twos, pks)
	}
	_ = idx
}

func Benchmark_cmp4_native(b *testing.B) {
	b.StopTimer()
	fours := [4]uint64{1, 2, 3, 4}
	pk := [4]uint64{2, 2, 2, 2}
	b.ResetTimer()
	b.StartTimer()

	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp4_native(fours, pk)
	}
	_ = idx
}

func Benchmark_cmp4_avx2(b *testing.B) {
	b.StopTimer()
	fours := [4]uint64{1, 2, 3, 4}
	pk := [4]uint64{2, 2, 2, 2}
	b.ResetTimer()
	b.StartTimer()

	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp4(fours, pk)
	}
	_ = idx
}

func Benchmark_cmp8_native(b *testing.B) {
	b.StopTimer()
	data := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	pk := [4]uint64{2, 2, 2, 2}
	b.ResetTimer()
	b.StartTimer()

	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp8_native(data, pk)
	}
	_ = idx
}

func Benchmark_cmp8_avx2(b *testing.B) {
	b.StopTimer()
	data := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	pk := [4]uint64{2, 2, 2, 2}
	b.ResetTimer()
	b.StartTimer()

	var idx int16
	for i := 0; i < b.N; i++ {
		idx = cmp8(data, pk)
	}
	_ = idx
}

const BENCHKEYS = 16384

type kv struct {
	k, v uint64
}

type kvs []kv

func (l kvs) Len() int           { return len(l) }
func (l kvs) Less(i, j int) bool { return l[i].k < l[j].k }

func BenchmarkSearchNaive(b *testing.B) {
	b.StopTimer()
	keys := make([]kv, BENCHKEYS/2)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(keys); j++ {
			idx = Naive(keys, uint64(j))
		}
	}
	_ = idx
}

func BenchmarkClever(b *testing.B) {
	b.StopTimer()
	keys := make([]uint64, BENCHKEYS)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(keys); j++ {
			idx = Clever(keys, uint64(j))
		}
	}
	_ = idx
}

func BenchmarkSearchAVX2(b *testing.B) {
	b.StopTimer()
	keys := make([]uint64, BENCHKEYS)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(keys); j++ {
			idx = skernel(keys, uint64(j))
		}
	}
	_ = idx

}

func BenchmarkSearchParallel(b *testing.B) {
	b.StopTimer()
	keys := make([]uint64, BENCHKEYS)
	for i := 0; i < len(keys); i += 2 {
		keys[i] = uint64(i)
		keys[i+1] = 1
	}
	b.ResetTimer()
	b.StartTimer()
	var idx int16
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(keys); j++ {
			idx = skernel(keys, uint64(j))
		}
	}
	_ = idx

}
