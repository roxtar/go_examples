package main
import (
	"fmt"
	"runtime"
	"flag"
)
const n = 3000000
var nprocs = 1
var nthreads = 10
func do_map(a []int64) {

	block_size := n/nthreads;
	ch := make(chan bool, nprocs)	

	for j:=0; j < nthreads; j++ {
		go map_block(a[j*block_size : (j+1)*block_size], ch)
	}

	// We simulate a "finish" by draining the channel

	for j:=0; j < nthreads; j++ {
	 	<- ch
	}
	close(ch)
}

func map_block(
	a [] int64,
	ch chan bool,
	) {

	for i, v := range a {
		a[i] = f(v)
	}
	ch <- true
}

func do_reduce(a []int64) int64{
	var total int64 = 0
	block_size := n/nthreads
	ch := make(chan int64)
	for j:=0; j < nthreads; j++ {
		go reduce_block(a[j*block_size:(j+1)*block_size], ch)
	}

	for j:=0; j < nthreads; j++ {
		total += <- ch
	}
	close(ch)

	return total
}

func reduce_block (
	a[] int64, 
	ch chan int64) {
	
	var total int64 = 0
	for _, v := range a{
		total += v
	}
	ch <- total
}

func f(x  int64) int64{
	return x * x
}


func main() {
	flag.IntVar(&nprocs, "n", 1, "Number of processors")
	flag.IntVar(&nthreads, "t", 10, "Number of threadas")
	flag.Parse()
	runtime.GOMAXPROCS(nprocs)
	a := make([] int64, n)
	var i int64;
	for i = 0; i < n; i++ {
		a[i] = i;
	}
	do_map(a)
	fmt.Printf("%v \n", do_reduce(a))
}