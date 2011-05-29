package main
import (
	"fmt"
	"runtime"
	"sync"
)
const n = 3000000
const nprocs = 4

func do_map(a []int64, map_udf func(* int64)) {

	block_size := n/nprocs;

	var wg sync.WaitGroup;

	for j:=0; j < nprocs; j++ {
		wg.Add(1)
		go map_block(a, j*block_size, (j+1)*block_size, &wg, map_udf)
	}

	wg.Wait()
}

func map_block(
	a [] int64,
	start int,
	end int,
	wg *sync.WaitGroup,
	map_udf func(* int64)) {

	for i:=start; i < end; i++ {
		map_udf(&a[i])
	}
	(*wg).Done()
	
}

func do_reduce(a []int64) int64{
	var total int64 = 0
	block_size := n/nprocs
	ch := make(chan int64)
	for j:=0; j < nprocs; j++ {
		go reduce_block(a, j*block_size, (j+1)*block_size, ch)
	}

	for j:=0; j < nprocs; j++ {
		total += <- ch
	}

	return total
}

func reduce_block (
	a[] int64, 
	start int,
	end int,
	ch chan int64) {
	
	var total int64 = 0
	for i:=start; i < end; i++ {
		total += a[i]
	}
	ch <- total
}

func f(x * int64) {
	*x = *x * *x
}


func main() {
	runtime.GOMAXPROCS(nprocs)
	a := make([] int64, n)
	var i int64;
	for i = 0; i < n; i++ {
		a[i] = i;
	}
	do_map(a, f)
	fmt.Printf("%v \n", do_reduce(a))
}