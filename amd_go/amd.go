package main

import (
	"fmt"
)

func md_all_pairs (dists []uint32, v uint32) {

	var k,i,j uint32;
	for k =0; k < v ; k++ {
		for i=0; i<v; i++ {
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
	}
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
	fmt.Printf("solution -- %d\n", solution);

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
//Main program - reads input, calls FW, shows output
func main() {
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
	fmt.Printf("len=%d cap=%d %v\n", len(dists), cap(dists), dists)
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
