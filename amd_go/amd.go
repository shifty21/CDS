package main

import (
	"fmt"
	"sync"
	"time"
	"runtime"
)

func md_all_pairs (dists []uint32, v uint32) {
	var wg sync.WaitGroup
	var div = v/100
	var k,i uint32;
	for k =0; k < v ; k++ {
		for i=0; i<v-div; i=i+div {
			//fmt.Printf("value range of i %d to %d \n", i,i+div)
			wg.Add(1);
			go internal_loop(dists, v, k,i,i+div,&wg)
		}
		wg.Add(1)
		go internal_loop(dists, v, k,i,v,&wg)
		wg.Wait()
	}
}


func internal_loop(dists []uint32, v uint32, k uint32,istart uint32,iend uint32, wg *sync.WaitGroup) {
	var j,i uint32;
	defer wg.Done()
	go func(){
		for i=istart;i <iend;i++{
			for j=0; j<v; j++ {
				var intermediary uint32 = dists[i*v+k] + dists[k*v+j];
				//check for overflows
				if ((intermediary >= dists[i*v+k]) &&
					(intermediary >= dists[k*v+j]) &&
					(intermediary < dists[i*v+j])){
					dists[i*v+j] = dists[i*v+k] + dists[k*v+j]
				}
		}
	}
	}()
}

func amd (dists []uint32,v uint32) {
	var i, j uint32;
	var infinity uint32 = v*v;
	var smd uint32 = 0;// sum of minimum distances
	var paths uint32 = 0;//number of paths
	var solution uint32 =0;

	for i=0;i<v;i++ {
		for j=0;j<v;j++ {
			if ((i!=j) && (dists[i*v+j] < infinity)) {
				smd += dists[i*v+j];
				paths++;
			}
		}
 	}
	solution = smd / paths;
	fmt.Printf("%d\n", solution);

}

func debug (dists []uint32, v uint32) {

	var i,j uint32;
	var infinity uint32 = v*v;

	for i=0;i<v;i++ {
		for j=0;j<v;j++ {
			if (dists[i*v+j] > infinity) {
				fmt.Printf("%7s", "inf");

			} else {
				fmt.Printf("%d", dists[i*v+j]);
			}
		}
		fmt.Print ("\n");
	}
}


func memsetRepeat(a []uint32, v uint32) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//Main program - reads input, calls FW, shows output
func main() {
//	runtime.GOMAXPROCS(1)
	//Read input
	//First line : v(number of vertices)  and e (number of edges)

//	fmt.Printf("No of available cores %d ", runtime.NumCPU())

	var v,e uint32;
	_, errv := fmt.Scanf("%d %d", &v, &e)
	if errv!=nil {
		fmt.Print("Error while reading v")
	}
	//allocates distances matrix (w/sice v*v)
	// and sets it with max distance and 0 for own vertex
	dists := make([]uint32, v*v);
	start:=makeTimestamp()
	memsetRepeat(dists, 1<<32 - 1)
	end:=makeTimestamp()
//	fmt.Printf("memsetRepeat time %d\n", end-start)
	var i uint32
	for i= 0; i <v; i++ {
		dists[i*v+i] = 0;
	}
	var source, dest ,cost uint32;
	for i=0;i<e; i++ {
		fmt.Scanf("%d %d %d", &source, &dest, &cost);
		if cost < dists[source*v+dest] {
			dists[source*v + dest] = cost;
		}
	}

	md_all_pairs(dists,v);

	amd( dists, v);

	const deb = false
	if deb {
		debug(dists,v)
	}




}
