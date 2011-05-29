package main
import fmt "fmt"
const n = 3000001

func do_map(a []int64, map_udf func(int64) int64) {
	var i int64;
	for i=0; i < n; i++ {
		a[i] = map_udf(a[i])
	}
}

func f(x int64) int64 {
	return x*x
}

func do_reduce(a []int64) int64{
	var total int64 = 0
	for i:=0; i < n; i++ {
		total += a[i]
	}
	return total
}

func main() {
	a := make([] int64, n)
	var i int64;
	for i = 0; i < n; i++ {
		a[i] = i;
	}
	do_map(a, f)
	fmt.Printf("%v \n", do_reduce(a))
}