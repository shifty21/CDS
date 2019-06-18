package main

import (
	"fmt"
	"sync"
	"time"
	"runtime"
	// "os"
	// "log"
	// "flag"
	// "runtime/pprof"
	// "runtime/trace"

)

func md_all_pairs (dists []uint32, v uint32) {
	var wg sync.WaitGroup
	var div uint32 = v/uint32(runtime.NumCPU())
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
	(*wg).Done()
	go func(){
		for i=istart;i <iend;i++{
			//pre calculating indexes 1.04m to 22.29s
			ivk := i*v+k
			temp_dists := dists[ivk]
			for j=0; j<v; j++ {
				// var intermediary uint32 = intermediary1 + intermediary2;
				kvj := k*v+j
				ivj := i*v+j
				temp_dists2 := dists[kvj]
				var intermediary uint32 = temp_dists + temp_dists2;
				//check for overflows
				if ((intermediary >= dists[ivk]) &&
					(intermediary >= dists[kvj]) &&
					(intermediary < dists[ivj])){
					dists[ivj] = dists[ivk] + dists[kvj]
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
// var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
// var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")
//Main program - reads input, calls FW, shows output
func main() {
	// trace.Start(os.Stderr)
	// defer trace.Stop()

	// flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create("cpu.prof")
	// 	if err != nil {
	// 		log.Fatal("could not create CPU profile: ", err)
	// 	}
	// 	defer f.Close()
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		log.Fatal("could not start CPU profile: ", err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	//Read input
	//First line : v(number of vertices)  and e (number of edges)
	var v,e uint32;
	_, errv := fmt.Scanf("%d %d", &v, &e)
	if errv!=nil {
		fmt.Print("Error while reading v")
	}
	//allocates distances matrix (w/sice v*v)
	// and sets it with max distance and 0 for own vertex
	dists := make([]uint32, v*v);
	memsetRepeat(dists, 1<<32 - 1)
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

	// ... rest of the program ...

	// if *memprofile != "" {
	// 	f, err := os.Create("mem.prof")
	// 	if err != nil {
	// 		log.Fatal("could not create memory profile: ", err)
	// 	}
	// 	defer f.Close()
	// 	runtime.GC() // get up-to-date statistics
	// 	if err := pprof.WriteHeapProfile(f); err != nil {
	// 		log.Fatal("could not write memory profile: ", err)
	// 	}
	// }



}
