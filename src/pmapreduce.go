package main
import (
	"fmt"
	"runtime"
	"flag"
)
const n = 3000000
var nprocs = 1
func do_map(a []int64) {

	block_size := n/nprocs;
	ch := make(chan bool, nprocs)	

	for j:=0; j < nprocs; j++ {
		go map_block(a, j*block_size, (j+1)*block_size, ch)
	}

	// We simulate a "finish" by draining the channel

	for j:=0; j < nprocs; j++ {
	 	<- ch
	}
	close(ch)
}

func map_block(
	a [] int64,
	start int,
	end int,
	ch chan bool,
	) {

	for i:=start; i < end; i++ {
		a[i] = f(a[i])
	}
	ch <- true
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
	close(ch)

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

func f(x  int64) int64{
	return x * x
}


func main() {
	flag.IntVar(&nprocs, "n", 1, "Number of threads")
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