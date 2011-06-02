// Program which demonstrates that only values
// not references are passed through channels
package main
import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(4)
	var i = 5
	fmt.Printf("Value of i in main %v\n", i)
	ch := make(chan int)	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		v := <- ch
		v = 7
		fmt.Printf("Changed v %v\n", v)
		wg.Done()
	}()
	ch <- i
	wg.Wait()
	fmt.Printf("Value of i in main %v\n", i)
}
	