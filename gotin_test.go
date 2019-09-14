package gotin_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/poloxue/gotin"
)

func TestOne(t *testing.T) {
	fmt.Println(sort.SearchInts([]int{1, 2, 6, 8, 9, 11}, 6))
	fmt.Println(sort.SearchInts([]int{1, 2, 6, 8, 9, 11}, 7))
}

func TestIn(t *testing.T) {
	cases := []struct {
		haystack interface{}
		needle   interface{}
		expect   bool
		err      error
	}{
		// unsupported haystack
		{haystack: 2, needle: 2, expect: false, err: gotin.ErrUnSupportHaystack},

		// slice
		// 1. int
		{haystack: []int{1, 2, 3}, needle: 3, expect: true, err: nil},
		{haystack: []int{1, 2, 3}, needle: 4, expect: false, err: nil},
		{haystack: []int{1, 2, 3}, needle: "three", expect: false, err: nil},

		// 2. string
		{haystack: []string{"one", "two", "three"}, needle: "one", expect: true, err: nil},

		// composite type
		{haystack: []interface{}{"one", "two", "three"}, needle: "one", expect: true, err: nil},
		{haystack: []interface{}{1, "two", 3}, needle: 3, expect: true, err: nil},

		// array
		{haystack: [3]int{1, 2, 3}, needle: 3, expect: true, err: nil},
		{haystack: [3]int{1, 2, 3}, needle: 4, expect: false, err: nil},

		{haystack: [10]string{"one", "two", "three"}, needle: "one", expect: true, err: nil},
	}

	for _, c := range cases {
		actual, err := gotin.In(c.haystack, c.needle)
		if actual != c.expect || err != c.err {
			t.Errorf("gotin.In(%v, %v) = (%t, %v), expect (%v, %v)", c.haystack, c.needle, actual, err, c.expect, c.err)
		}
	}
}

func TestInIntSlice(t *testing.T) {
	cases := []struct {
		haystack []int
		needle   int
		expect   bool
	}{
		{haystack: []int{1, 2, 3, 4, 5}, needle: 5, expect: true},
		{haystack: []int{1, 2, 3, 4, 5}, needle: 0, expect: false},
	}

	for _, c := range cases {
		actual := gotin.InIntSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InIntSlice(%v, %d) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestInStringSlice(t *testing.T) {
	cases := []struct {
		haystack []string
		needle   string
		expect   bool
	}{
		{haystack: []string{"one", "two", "three", "four"}, needle: "one", expect: true},
		{haystack: []string{"one", "two", "three"}, needle: "five", expect: false},
	}

	for _, c := range cases {
		actual := gotin.InStringSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InSliceString(%v, %s) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestInIntSliceSortedFunc(t *testing.T) {
	cases := []struct {
		haystack []int
		needle   int
		expect   bool
	}{
		{haystack: []int{5, 1, 3, 2, 8, 6, 7}, needle: 1, expect: true},
		{haystack: []int{5, 2, 1, 3, 7, 8, 6}, needle: 4, expect: false},
	}

	for _, c := range cases {
		actual := gotin.InIntSliceSortedFunc(c.haystack)(c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InSliceIntSortedFunc(%v)(%d) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestInStringSliceSortedFunc(t *testing.T) {
	cases := []struct {
		haystack []string
		needle   string
		expect   bool
	}{
		{haystack: []string{"one", "two"}, needle: "one", expect: true},
		{haystack: []string{"one", "two", "three"}, needle: "four", expect: false},
	}

	for _, c := range cases {
		actual := gotin.InStringSliceSortedFunc(c.haystack)(c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InStringSliceSortedFunc(%v, %s) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestSortInIntSlice(t *testing.T) {
	cases := []struct {
		haystack []int
		needle   int
		expect   bool
	}{
		{haystack: []int{5, 1, 3, 2, 8, 6, 7}, needle: 1, expect: true},
		{haystack: []int{5, 2, 1, 3, 7, 8, 6}, needle: 4, expect: false},
	}

	for _, c := range cases {
		actual := gotin.SortInIntSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.SortInIntSlice(%v, %d) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestSortInStringSlice(t *testing.T) {
	cases := []struct {
		haystack []string
		needle   string
		expect   bool
	}{
		{haystack: []string{"one", "two"}, needle: "one", expect: true},
		{haystack: []string{"one", "two", "three"}, needle: "four", expect: false},
	}

	for _, c := range cases {
		actual := gotin.SortInStringSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.SortInStringSlice(%v, %s) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestInIntSliceMapKeyFunc(t *testing.T) {
	cases := []struct {
		haystack []int
		needle   int
		expect   bool
	}{
		{haystack: []int{4, 3, 1, 2, 0}, needle: 1, expect: true},
		{haystack: []int{8, 1, 4, 0}, needle: 5, expect: false},
	}

	for _, c := range cases {
		actual := gotin.InIntSliceMapKeyFunc(c.haystack)(c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InIntSliceMapKeyFunc(%v)(%d) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestInStringSliceMapKeyFunc(t *testing.T) {
	cases := []struct {
		haystack []string
		needle   string
		expect   bool
	}{
		{haystack: []string{"one", "two"}, needle: "one", expect: true},
		{haystack: []string{"one", "two", "three"}, needle: "four", expect: false},
	}

	for _, c := range cases {
		actual := gotin.InStringSliceMapKeyFunc(c.haystack)(c.needle)
		if actual != c.expect {
			t.Errorf("gotin.InStringSliceMapKeyFunc(%v)(%s) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestMapKeyInIntSlice(t *testing.T) {
	cases := []struct {
		haystack []int
		needle   int
		expect   bool
	}{
		{haystack: []int{4, 2, 3, 8, 6, 1, 5}, needle: 4, expect: true},
		{haystack: []int{5, 2, 1, 8, 7, 3}, needle: 4, expect: false},
	}

	for _, c := range cases {
		actual := gotin.MapKeyInIntSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.MapKeyInStringSlice(%v, %d) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func TestMapKeyInStringSlice(t *testing.T) {
	cases := []struct {
		haystack []string
		needle   string
		expect   bool
	}{
		{haystack: []string{"one", "two"}, needle: "one", expect: true},
		{haystack: []string{"one", "two", "three"}, needle: "four", expect: false},
	}

	for _, c := range cases {
		actual := gotin.MapKeyInStringSlice(c.haystack, c.needle)
		if actual != c.expect {
			t.Errorf("gotin.MapKeyInStringSlice(%v, %s) = %t, expect %t", c.haystack, c.needle, actual, c.expect)
		}
	}
}

func randomHaystackAndNeedle(size int) ([]int, int){
	haystack := make([]int, size)

	for i := 0; i<size ; i++{
		haystack[i] = rand.Int()
	}

	return haystack, rand.Int()
}

func BenchmarkIn_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gotin.In(haystack, needle)
	}
}

func BenchmarkIn_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gotin.In(haystack, needle)
	}
}

func BenchmarkIn_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gotin.In(haystack, needle)
	}
}

func BenchmarkInIntSlice_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.InIntSlice(haystack, needle)
	}
}

func BenchmarkInIntSlice_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.InIntSlice(haystack, needle)
	}
}

func BenchmarkInIntSlice_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.InIntSlice(haystack, needle)
	}
}

func BenchmarkInIntSliceSortedFunc_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	in := gotin.InIntSliceSortedFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkInIntSliceSortedFunc_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)

	in := gotin.InIntSliceSortedFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkInIntSliceSortedFunc_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	in := gotin.InIntSliceSortedFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkSortInIntSlice_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.SortInIntSlice(haystack, needle)
	}
}

func BenchmarkSortInIntSlice_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.SortInIntSlice(haystack, needle)
	}
}

func BenchmarkSortInIntSlice_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.SortInIntSlice(haystack, needle)
	}
}

func BenchmarkInIntSliceMapKeyFunc_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	in := gotin.InIntSliceMapKeyFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkInIntSliceMapKeyFunc_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)

	in := gotin.InIntSliceMapKeyFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkInIntSliceMapKeyFunc_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	in := gotin.InIntSliceMapKeyFunc(haystack)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = in(needle)
	}
}

func BenchmarkMapKeyInIntSlice_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.MapKeyInIntSlice(haystack, needle)
	}
}

func BenchmarkMapKeyInIntSlice_1000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e3)


	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.MapKeyInIntSlice(haystack, needle)
	}
}

func BenchmarkMapKeyInIntSlice_1000000(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(1e6)

	b.ResetTimer()
	for i := 0; i<b.N; i++ {
		_ = gotin.MapKeyInIntSlice(haystack, needle)
	}
}
