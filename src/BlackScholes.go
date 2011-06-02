package main;

import (
	"fmt"
	"runtime"
	"math"
	"io/ioutil"
	"strconv"
	"strings"
	"flag"
	);

const inv_sqrt_2xPI=0.39894228040143270286;
const NUM_RUNS = 100;
const N=65536;

	
type OptionData struct {
        s float64;          // spot price
        strike float64;     // strike price
        r float64;          // risk-free interest rate
        divq float64;       // dividend rate
        v float64;          // volatility
        t float64;          // time to maturity or option expiration in years 
                           //     (1yr = 1.0, 6mos = 0.5, 3mos = 0.25, ..., etc)  
        OptionType uint8;   // Option type.  "P"=PUT, "C"=CALL
        divs float64;       // dividend vals (not used in this test)
        DGrefval float64;   // DerivaGem Reference Value
} 

var data [N]OptionData;

var prices [N]float64;
var numOptions int = N;

var otype [N]bool;
var sptprice [N]float64;
var strike [N]float64;
var rate [N]float64;
var volatility [N]float64;
var otime [N]float64;
var numError int;
var nThreads int;

func CNDF ( InputX float64 ) float64 {
    var sign int32;
    var OutputX float64;
    var xInput float64;
    var xNPrimeofX float64;
    var expValues float64;
    var xK2 float64;
    var xK2_2 float64;
    var xK2_3 float64;
    var xK2_4 float64;
    var xK2_5 float64;
    var xLocal float64;
    var xLocal_1 float64;
    var xLocal_2 float64;
    var xLocal_3 float64;

    // Check for negative value of InputX
    if (InputX < 0.0) {
        InputX = -InputX;
        sign = 1;
    } else 
        sign = 0;

    xInput = InputX;
    // Compute NPrimeX term common to both four & six decimal accuracy calcs
    expValues =   math.Exp(-0.5 * InputX * InputX);
    xNPrimeofX = expValues;
    xNPrimeofX = xNPrimeofX * inv_sqrt_2xPI;

    xK2 = 0.2316419 * xInput;
    xK2 = 1.0 + xK2;
    xK2 = 1.0 / xK2;
    xK2_2 = xK2 * xK2;
    xK2_3 = xK2_2 * xK2;
    xK2_4 = xK2_3 * xK2;
    xK2_5 = xK2_4 * xK2;
    
    xLocal_1 = xK2 * 0.319381530;
    xLocal_2 = xK2_2 * (-0.356563782);
    xLocal_3 = xK2_3 * 1.781477937;
    xLocal_2 = xLocal_2 + xLocal_3;
    xLocal_3 = xK2_4 * (-1.821255978);
    xLocal_2 = xLocal_2 + xLocal_3;
    xLocal_3 = xK2_5 * 1.330274429;
    xLocal_2 = xLocal_2 + xLocal_3;

    xLocal_1 = xLocal_2 + xLocal_1;
    xLocal   = xLocal_1 * xNPrimeofX;
    xLocal   = 1.0 - xLocal;

    OutputX  = xLocal;
    
    if (sign > 0) {
        OutputX = 1.0 - OutputX;
    }
    return OutputX;
} 

func BlkSchlsEqEuroNoDiv( sptprice float64,
                            strike float64, rate float64, volatility float64,
                            time float64,  otype bool,  timet float64 ) float64 {
    var OptionPrice float64;
    // local private working variables for the calculation
//    var xStockPrice float64;
//    var xStrikePrice float64;
    var xRiskFreeRate float64;
    var xVolatility float64;
    var xTime float64;
    var xSqrtTime float64;

    var logValues float64;
    var xLogTerm float64;
    var xD1 float64; 
    var xD2 float64;
    var xPowerTerm float64;
    var xDen float64;
    var d1 float64;
    var d2 float64;
    var FutureValueX float64;
    var NofXd1 float64;
    var NofXd2 float64;
    var NegNofXd1 float64;
    var NegNofXd2 float64;    
    
//    xStockPrice := sptprice;
//    xStrikePrice := strike;
    xRiskFreeRate = rate;
    xVolatility = volatility;

    xTime = time;
    xSqrtTime = math.Sqrt(xTime);

    logValues = math.Log( sptprice / strike );
        
    xLogTerm = logValues;
        
    
    xPowerTerm = xVolatility * xVolatility;
    xPowerTerm = xPowerTerm * 0.5;
        
    xD1 = xRiskFreeRate + xPowerTerm;
    xD1 = xD1 * xTime;
    xD1 = xD1 + xLogTerm;

    xDen = xVolatility * xSqrtTime;
    xD1 = xD1 / xDen;
    xD2 = xD1 -  xDen;

    d1 = xD1;
    d2 = xD2;
    
    NofXd1 = CNDF( d1 );
    NofXd2 = CNDF( d2 );

    FutureValueX = strike * ( math.Exp( -(rate)*(time) ) );        
    if (!otype) {            
        OptionPrice = (sptprice * NofXd1) - (FutureValueX * NofXd2);
    } else { 
        NegNofXd1 = (1.0 - NofXd1);
        NegNofXd2 = (1.0 - NofXd2);
        OptionPrice = (FutureValueX * NegNofXd2) - (sptprice * NegNofXd1);
    }
    
    return OptionPrice;
}

func bs_thread(tid int, finish chan bool) {
    var price float64;
//    var priceDelta float64;
    start := tid * (numOptions / nThreads);
    end := start + (numOptions / nThreads);

    for j:=0; j<NUM_RUNS; j++ {
        for i:=start; i<end; i++ {
            price = BlkSchlsEqEuroNoDiv( sptprice[i], strike[i],
                                         rate[i], volatility[i], otime[i], 
                                         otype[i], 0);
            prices[i] = price;
            priceDelta := data[i].DGrefval - price;
            if( math.Fabs(priceDelta) >= 1e-4 ){
                fmt.Println("Error on %d. Computed=%.5f, Ref=%.5f, Delta=%.5f\n",i, price, data[i].DGrefval, priceDelta);
                numError = numError+1;
            }
        }
 	}
 	
 	finish <- true;
 }


func main () {
    var NTHREADS int
    flag.IntVar(&NTHREADS, "n", 1, "Number of threads")
	flag.Parse()
    fmt.Println("PARSEC Benchmark Suite\n");
    //Read input data from file 
    inputFile := "in_64K.txt"
    fp, _ := ioutil.ReadFile(inputFile);
    fields := strings.Fields(string(fp));
    k := 0;
    numOptions = N;
    nThreads = NTHREADS;
    runtime.GOMAXPROCS(int(nThreads));
   
    for i:= 1; i< len(data); i+=9 {
	    data[k].s, _ = strconv.Atof64(fields[i]);
	    data[k].strike, _ = strconv.Atof64(fields[i+1]);
	    data[k].r, _ =   strconv.Atof64(fields[i+2]);
	    data[k].divq, _ = strconv.Atof64(fields[i+3]);
	    data[k].v, _ = strconv.Atof64(fields[i+4]);
	    data[k].t, _ = strconv.Atof64(fields[i+5]);
	    data[k].OptionType =  fields[i+6][0];
	    data[k].divs, _ = strconv.Atof64(fields[i+7]);    
	    data[k].DGrefval, _ = strconv.Atof64(fields[i+8]);
	    k++;    	
    }
    
    fmt.Printf("Num of Options: %d\n", numOptions);
    fmt.Printf("Num of Runs: %d\n", NUM_RUNS);
        
     for i:=0; i<numOptions; i++ {
        otype[i]      = data[i].OptionType == 'P';
        sptprice[i]   = data[i].s;
        strike[i]     = data[i].strike;
        rate[i]       = data[i].r;
        volatility[i] = data[i].v;    
        otime[i]      = data[i].t;
    }
    var finish chan bool;
    finish =make(chan bool, NTHREADS)
    //Launch the appropriate number of threads
    for i:= 0; i<nThreads; i++ {
    	go bs_thread(i, finish);
    }
    
    for i:=0; i<nThreads; i++ {
    	<- finish;
    }
    
    


}
