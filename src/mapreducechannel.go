// Another version of map reduce which doesn't use
// any arrays. Everything is done through channels.
// 
// Is much more elegant than pmapreduce.go, but 
// also much slower.

package main
import (
	"fmt"
	"runtime"
)

func f(x int64) int64{
	return x*x
}


func do_map(in chan int64, out chan int64) {	
	for i:= range in{
		out <- f(i)
	}
	close(out)
}

func do_reduce(in chan int64, out chan int64){
	var total int64 = 0
	for i:= range in{
		total += i
	}	
	out <- total
	close(out)
}

func main() {

	const n = 3000000
	const nprocs = 4

	runtime.GOMAXPROCS(nprocs)

	map_in := make(chan int64)
	map_out := make(chan int64)
	reduce_out := make(chan int64)

	go do_map(map_in, map_out)
	go do_reduce(map_out, reduce_out)

	for i := 0; i < n; i++ {
		map_in <- int64(i)
	}
	close(map_in)
	fmt.Printf("%v \n", <- reduce_out)	
}