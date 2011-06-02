// Shared memory demo
package main
import (
	"fmt"
	"runtime")
var v = 45
func main() {
	runtime.GOMAXPROCS(4)	
	go func() {
		for {
			v = 34
			fmt.Printf("thread 1\n");
		}
	}()

	go func() {
		for {
			v = 66
			fmt.Printf("thread 2\n");
		}
	}()

	for {
		fmt.Printf("%v\n", v)
	}
}
		