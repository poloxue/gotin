# gotin

in 是一个很常用的功能，有些情况可能也称为 contains，虽然不同语言的表示不同，但基本都是有的。不过可惜的是，Go 却没有，它即没有提供类似 Python 操作符 in，也没有像其他语言那样提供这样的标准库函数，如 PHP 中 in_array。

Go 的哲学是追求少即是多。我想或许 Go 团队觉得这是一个实现起来不足为道的功能吧。

为何说微不足道？如果要自己实现，又该如何做呢？

我所想到的有三种实现方式，一是遍历，二是 sort 的二分查找，三是 map 的 key 索引。

本文相关源码已经上传在我的 github 上，[poloxue/gotin](https://github.com/poloxue/gotin)。

## 遍历

遍历应该是我们最容易想到的最简单的实现方式。

示例如下：

```go
func InIntSlice(haystack []int, needle int) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}
```

上面演示了如何在一个 []int 类型变量中查找指定 int 是否存在的例子，是不是非常简单，由此我们也可以感受到我为什么说它实现起来微不足道。

这个例子有个缺陷，它只支持单一类型。如果要支持像解释语言一样的通用 in 功能，则需借助反射实现。

代码如下：

```go
func In(haystack interface{}, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, ErrUnSupportHaystack
}
```

为了更加通用，In 函数的输入参数 haystack 和 needle 都是 interface{} 类型。

简单说说输入参数都是 interface{} 的好处吧，主要有两点，如下：

其一，haystack 是 interface{} 类型，使 in 支持的类型不止于 slice，还包括 array。我们看到，函数内部通过反射对 haystack 进行了类型检查，支持 slice（切片）与 array（数组）。如果是其他类型则会提示错误，增加新的类型支持，如 map，其实也很简单。但不推荐这种方式，因为通过 _, ok := m[k] 的语法即可达到 in 的效果。

其二，haystack 是 interface{}，则 []interface{} 也满足要求，并且 needle 是 interface{}。如此一来，我们就可以实现类似解释型语言一样的效果了。

怎么理解？直接示例说明，如下：

```go
gotin.In([]interface{}{1, "two", 3}, "two")
```

haystack 是 []interface{}{1, "two", 3}，而且 needle 是 interface{}，此时的值是 "two"。如此看起来，是不是实现了解释型语言中，元素可以是任意类型，不必完全相同效果。如此一来，我们就可以肆意妄为的使用了。

但有一点要说明，In 函数的实现中有这样一段代码：

```go
if sVal.Index(i).Interface() == needle {
	...
}
```

Go 中并非任何类型都可以使用 == 比较的，如果元素中含有 slice 或 map，则可能会报错。

## 二分查找

以遍历确认元素是否存在有个缺点，那就是，如果数组或切片中包含了大量数据，比如 1000000 条数据，即一百万，最坏的情况是，我们要遍历 1000000 次才能确认，时间复杂度 On。

有什么办法可以降低遍历次数？

自然而然地想到的方法是二分查找，它的时间复杂度 log2(n) 。但这个算法有前提，需要依赖有序序列。

于是，第一个要我们解决的问题是使序列有序，Go 的标准库已经提供了这个功能，在 sort 包下。

示例代码如下：

```go
fmt.Println(sort.SortInts([]int{4, 2, 5, 1, 6}))
```

对于 []int，我们使用的函数是 SortInts，如果是其他类型切片，sort 也提供了相关的函数，比如 []string 可通过 SortStrings 排序。

完成排序就可以进行二分查找，幸运的是，这个功能 Go 也提供了，[]int 类型对应函数是 SearchInts。

简单介绍下这个函数，先看定义：

```go
func SearchInts(a []int, x int) int
```

输入参数容易理解，从切片 a 中搜索 x。重点要说下返回值，这对于我们后面确认元素是否存在至关重要。返回值的含义，返回查找元素在切片中的位置，如果元素不存在，则返回，在保持切片有序情况下，插入该元素应该在什么位置。

比如，序列如下：

```
1 2 6 8 9 11
```

假设，x 为 6，查找之后将发现它的位置在索引 2 处；x 如果是 7，发现不存在该元素，如果插入序列，将会放在 6 和 8 之间，索引位置是 3，因而返回值为 3。

代码测试下：

```go
fmt.Println(sort.SearchInts([]int{1, 2, 6, 8, 9, 11}, 6)) // 2
fmt.Println(sort.SearchInts([]int{1, 2, 6, 8, 9, 11}, 7)) // 3
```

如果判断元素是否在序列中，只要判断返回位置上的值是否和查找的值相同即可。

但还有另外一种情况，如果插入元素位于序列最后，例如元素值为 12，插入位置即为序列的长度 6。如果直接查找 6 位置上的元素就可能发生越界的情况。那怎么办呢？其实判断返回是否大于切片长度即可，大于则说明元素不在切片序列中。

完整的实现代码如下：

```go
func SortInIntSlice(haystack []int, needle int) bool {
	sort.Ints(haystack)

	index := sort.SearchInts(haystack, needle)
	return index < len(haystack) && haystack[index] == needle
}
```

但这还有个问题，对于无序的场景，如果每次查询都要经过一次排序并不划算。最后能实现一次排序，稍微修改下代码。

```go
func InIntSliceSortedFunc(haystack []int) func(int) bool {
	sort.Ints(haystack)

	return func(needle int) bool {
		index := sort.SearchInts(haystack, needle)
		return index < len(haystack) && haystack[index] == needle
	}
}
```

上面的实现，我们通过调用 InIntSliceSortedFunc 对 haystack 切片排序，并返回一个可多次使用的函数。

使用案例如下：

```go
in := gotin.InIntSliceSortedFunc(haystack)

for i := 0; i<maxNeedle; i++ {
	if in(i) {
		fmt.Printf("%d is in %v", i, haystack)
	}
}
```

二分查找的方式有什么不足呢？

我想到的重要一点，要实现二分查找，元素必须是可排序的，如 int，string，float 类型。而对于结构体、切片、数组、映射等类型，使用起来就不是那么方便，当然，如果要用，也是可以的，不过需要我们进行一些适当扩展，按指定标准排序，比如结构的某个成员。

到此，二分查找的 in 实现就介绍完毕了。

## map key

本节介绍 map key 方式。它的算法复杂度是 O1，无论数据量多大，查询性能始终不变。它主要依赖的是 Go 中的 map 数据类型，通过 hash map 直接检查 key 是否存在，算法大家应该都比较熟悉，通过 key 可直接映射到索引位置。

我们常会用到这个方法。

```go
_, ok := m[k]
if ok {
	fmt.Println("Found")
}
```

那么它和 in 如何结合呢？一个案例就说明白了这个问题。

假设，我们有一个 []int 类型变量，如下：

```go
s := []int{1, 2, 3}
```

为了使用 map 的能力检查某个元素是否存在，可以将 s 转化 map[int]struct{}。

```go
m := map[interface{}]struct{}{
	1: struct{}{},
	2: struct{}{},
	3: struct{}{},
	4: struct{}{},
}
```

如果检查某个元素是否存在，只需要通过如下写法即可确定：

```go
k := 4
if _, ok := m[k]; ok {
	fmt.Printf("%d is found\n", k)
}
```

是不是非常简单？

补充一点，关于这里为什么使用 struct{}，可以阅读我之前写的一篇关于 [Go 中如何使用 set](https://mp.weixin.qq.com/s/a0BWRTikJNTPc6VXVn43OQ) 的文章。

按照这个思路，实现函数如下：

```go
func MapKeyInIntSlice(haystack []int, needle int) bool {
	set := make(map[int]struct{})

	for _ , e := range haystack {
		set[e] = struct{}{}
	}

	_, ok := set[needle]
	return ok
}
```

实现起来不难，但和二分查找有着同样的问题，开始要做数据处理，将 slice 转化为 map。如果是每次数据相同，稍微修改下它的实现。

```go
func InIntSliceMapKeyFunc(haystack []int) func(int) bool {
	set := make(map[int]struct{})

	for _ , e := range haystack {
		set[e] = struct{}{}
	}

	return func(needle int) bool {
		_, ok := set[needle]
		return ok
	}
}
```

对于相同的数据，它会返回一个可多次使用的 in 函数，一个使用案例如下：

```go
in := gotin.InIntSliceMapKeyFunc(haystack)

for i := 0; i<maxNeedle; i++ {
	if in(i) {
		fmt.Printf("%d is in %v", i, haystack)
	}
}
```

对比前两种算法，这种方式的处理效率最高，非常适合于大数据的处理。接下来的性能测试，我们将会看到效果。

## 性能

介绍完所有方式，我们来实际对比下每种算法的性能。测试源码位于 gotin_test.go 文件中。

基准测试主要是从数据量大小考察不同算法的性能，本文中选择了三个量级的测试样本数据，分别是 10、1000、1000000。

为便于测试，首先定义了一个用于生成 haystack 和 needle 样本数据的函数。

代码如下：

```go
func randomHaystackAndNeedle(size int) ([]int, int){
	haystack := make([]int, size)

	for i := 0; i<size ; i++{
		haystack[i] = rand.Int()
	}

	return haystack, rand.Int()
}
```

输入参数是 size，通过 rand.Int() 随机生成切片大小为 size 的 haystack 和 1 个 needle。在基准测试用例中，引入这个随机函数生成数据即可。

举个例子，如下：

```go
func BenchmarkIn_10(b *testing.B) {
	haystack, needle := randomHaystackAndNeedle(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gotin.In(haystack, needle)
	}
}
```

首先，通过 randomHaystackAndNeedle 随机生成了一个含有 10 个元素的切片。因为生成样本数据的时间不应该计入到基准测试中，我们使用 b.ResetTimer() 重置了时间。

其次，压测函数是按照 `Test+函数名+样本数据量` 规则编写，如案例中 BenchmarkIn_10，表示测试 In 函数，样本数据量为 10。如果我们要用 1000 数据量测试 InIntSlice，压测函数名为 BenchmarkInIntSlice_1000。

测试开始吧！简单说下我的笔记本配置，Mac Pro 15 版，16G 内存，512 SSD，4 核 8 线程的 CPU。

测试所有函数在数据量在 10 的情况下的表现。

```bash
$ go test -run=none -bench=10$ -benchmem
```

匹配所有以 10 结尾的压测函数。

测试结果：

```bash
goos: darwin
goarch: amd64
pkg: github.com/poloxue/gotin
BenchmarkIn_10-8                         3000000               501 ns/op             112 B/op         11 allocs/op
BenchmarkInIntSlice_10-8                200000000                7.47 ns/op            0 B/op          0 allocs/op
BenchmarkInIntSliceSortedFunc_10-8      100000000               22.3 ns/op             0 B/op          0 allocs/op
BenchmarkSortInIntSlice_10-8            10000000               162 ns/op              32 B/op          1 allocs/op
BenchmarkInIntSliceMapKeyFunc_10-8      100000000               17.7 ns/op             0 B/op          0 allocs/op
BenchmarkMapKeyInIntSlice_10-8           3000000               513 ns/op             163 B/op          1 allocs/op
PASS
ok      github.com/poloxue/gotin        13.162s
```

表现最好的并非 SortedFunc 和 MapKeyFunc，而是最简单的针对单类型的遍历查询，平均耗时 7.47ns/op，当然，另外两种方式表现也不错，分别是 22.3ns/op 和 17.7ns/op。

表现最差的是 In、SortIn（每次重复排序） 和 MapKeyIn（每次重复创建 map）两种方式，平均耗时分别为 501ns/op 和 513ns/op。

测试所有函数在数据量在 1000 的情况下的表现。

```bash
$ go test -run=none -bench=1000$ -benchmem
```

测试结果：

```bash
goos: darwin
goarch: amd64
pkg: github.com/poloxue/gotin
BenchmarkIn_1000-8                         30000             45074 ns/op            8032 B/op       1001 allocs/op
BenchmarkInIntSlice_1000-8               5000000               313 ns/op               0 B/op          0 allocs/op
BenchmarkInIntSliceSortedFunc_1000-8    30000000                44.0 ns/op             0 B/op          0 allocs/op
BenchmarkSortInIntSlice_1000-8             20000             65401 ns/op              32 B/op          1 allocs/op
BenchmarkInIntSliceMapKeyFunc_1000-8    100000000               17.6 ns/op             0 B/op          0 allocs/op
BenchmarkMapKeyInIntSlice_1000-8           20000             82761 ns/op           47798 B/op         65 allocs/op
PASS
ok      github.com/poloxue/gotin        11.312s
```

表现前三依然是 InIntSlice、InIntSliceSortedFunc 和 InIntSliceMapKeyFunc，但这次顺序发生了变化，MapKeyFunc 表现最好，17.6 ns/op，与数据量 10 的时候相比基本无变化。再次验证了前文的说法。

同样的，数据量 1000000 的时候。

```bash
$ go test -run=none -bench=1000000$ -benchmem
```

测试结果如下：

```bash
goos: darwin
goarch: amd64
pkg: github.com/poloxue/gotin
BenchmarkIn_1000000-8                                 30          46099678 ns/op         8000098 B/op    1000001 allocs/op
BenchmarkInIntSlice_1000000-8                       3000            424623 ns/op               0 B/op          0 allocs/op
BenchmarkInIntSliceSortedFunc_1000000-8         20000000                72.8 ns/op             0 B/op          0 allocs/op
BenchmarkSortInIntSlice_1000000-8                     10         138873420 ns/op              32 B/op          1 allocs/op
BenchmarkInIntSliceMapKeyFunc_1000000-8         100000000               16.5 ns/op             0 B/op          0 allocs/op
BenchmarkMapKeyInIntSlice_1000000-8                   10         156215889 ns/op        49824225 B/op      38313 allocs/op
PASS
ok      github.com/poloxue/gotin        15.178s
```

MapKeyFunc 依然表现最好，每次操作用时 17.2 ns，Sort 次之，而 InIntSlice 呈现线性增加的趋势。一般情况下，如果不是对性能要特殊要求，数据量特别大的场景，针对单类型的遍历已经有非常好的性能了。

从测试结果可以看出，反射实现的通用 In 函数每次执行需要进行大量的内存分配，方便的同时，也是以牺牲性能为代价的。

## 总结

本文通过一个问题引出主题，为什么 Go 中没有类似 Python 的 In 方法。我认为，一方面是实现非常简单，没有必要。除此以外，另一方面，在不同场景下，我们还需要根据实际情况分析用哪种方式实现，而不是一种固定的方式。

接着，我们介绍了 In 实现的三种方式，并分析了各自的优劣。通过性能分析测试，我们能得出大致的结论，什么方式适合什么场景，但总体还是不能说足够细致，有兴趣的朋友可以继续研究下。

## 参考

[Does Go have “if x in” construct similar to Python?](https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python)

[为什么Golang没有像Python中in一样的功能？](https://www.zhihu.com/question/328393303/answer/711287362)
