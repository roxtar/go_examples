package main;

import (
	"fmt"
	"runtime"
	"flag"
)

const N = 30000000
var NPROCS = 4
var a [N]int64

func doInit() {
	for i :=0; i<len(a); i++ {
		a[i] = int64(i)
	}
}

func doParallelPrefixMain() {
	var ch []chan int64
	var finish chan bool
	
	//Initialize all channels
	finish = make(chan bool, NPROCS)
	ch = make([] chan int64, NPROCS)
	for i:= 0; i<NPROCS; i++ {
		ch[i] = make(chan int64, NPROCS)
	}
	
	//Launch all threads
	for i:=0; i<NPROCS; i++ {
		go threadParallelPrefix(i, ch[:], finish)
	}
	
	//Equivalent to a finish
	for i:=0; i<NPROCS; i++ {
		<- finish
	}
	fmt.Printf("%v\n", a[N-1])
	
}

func threadParallelPrefix(pid int, ch []chan int64, finish chan bool) {
	//First do a local parallel prefix computation
	var startIndex int64 = int64(pid) * N/int64(NPROCS)
	var endIndex int64 = startIndex + (N/int64(NPROCS))
	var toBeAdded int64 = 0

	for i:=startIndex+1; i<endIndex; i++ {		
		a[i] = a[i] + a[i-1]
	}
	
	for i:= pid+1; i<NPROCS; i++ {		
		ch[i] <- a[i-1]
	}
	
	for i:=0; i < pid; i++ {
		 toBeAdded += <- ch[pid]
	}
			
	for j:=startIndex; j<endIndex; j++ {
		a[j] += toBeAdded
	}
	 
	finish <- true	
}  
func printArray() {
for i:=0; i<len(a); i++ {
	fmt.Println("a[",i,"] = ",a[i]);
}
}

func main() {
	flag.IntVar(&NPROCS, "n", 1, "Number of threads")
	flag.Parse()
	runtime.GOMAXPROCS(int(NPROCS))
	doInit()
	doParallelPrefixMain()
}
