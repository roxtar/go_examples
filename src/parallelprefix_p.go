package main;

import (
	"fmt"
	"runtime"
	"flag"
)

const N int64 = 30000000
var NPROCS int64 = 0
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
	ch = make([]chan int64, NPROCS)
	finish = make(chan bool, NPROCS)
	for i:= int64(0); i<NPROCS; i++ {
		ch[i] = make(chan int64, NPROCS)
	}
	
	//Launch all threads
	for i:=int64(0); i<NPROCS; i++ {
		go threadParallelPrefix(i, ch, finish)
	}
	
	//Equivalent to a reduce
	for i:=int64(0); i<NPROCS; i++ {
		<- finish
	}
	
}

func threadParallelPrefix(pid int64, ch []chan int64, finish chan bool) {
	//First do a local parallel prefix computation
	var startIndex int64 = pid*N/NPROCS
	var endIndex int64 = startIndex + (N/NPROCS)
	var localSum int64 = a[startIndex]
	for i:=startIndex+1; i<endIndex; i++ {
		localSum += a[i]
		a[i] = localSum
	}
	
	for i:= pid+1; i<NPROCS; i++ {
		
		ch[i] <- localSum
	}
	
	var toBeAdded int64 = 0
	for i := int64(0); i<pid; i++ {
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

func init() {
	flag.Int64Var(&NPROCS, "n", 1, "Number of threads");
	flag.Parse();
}

func main() {
	runtime.GOMAXPROCS(int(NPROCS))
	doInit()
	doParallelPrefixMain()
	fmt.Println(a[len(a)-1])
//	printArray();
// 	fmt.Println("Hello world!!")
}