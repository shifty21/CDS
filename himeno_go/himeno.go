package main 
import (
	"fmt"
	"sync"
	"log"
	"runtime/pprof"
	"flag"
	"os"
)


type Matrix struct {
	m []float64
	mnums int
	mrows int
	mcols int
	mdeps int
}


func MR_set (mat* Matrix, n int, r int, c int, d int, val float64) {
	index :=  (n)*mat.mrows*mat.mcols*mat.mdeps +
		(r)*mat.mcols*mat.mdeps + (c)*mat.mdeps + (d)
	mat.m[index] = val
}

func MR_get(mat* Matrix, n int, r int, c int, d int) (float64) {
		index :=  (n)*mat.mrows*mat.mcols*mat.mdeps +
			(r)*mat.mcols*mat.mdeps + (c)*mat.mdeps + (d)
		return mat.m[index]

}

var omega float64 = 0.8
var a,b,c,p,bnd,wrk1,wrk2 Matrix;

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	// trace.Start(os.Stderr)
	// defer trace.Stop()

	var nn int;
	var mimax,mjmax,mkmax int;
	var msize [3]int;
	var gosa float64;

	fmt.Scanf("%d", &msize[0]);
	fmt.Scanf("%d", &msize[1]);
	fmt.Scanf("%d", &msize[2]);
	fmt.Scanf("%d", &nn);

	mimax = msize[0];
	mjmax = msize[1];
	mkmax = msize[2];
	/*
   *    Initializing matrixes
   */
	new_mat(&p,1,mimax,mjmax,mkmax);
	new_mat(&bnd,1,mimax,mjmax,mkmax);
	new_mat(&wrk1,1,mimax,mjmax,mkmax);
	new_mat(&wrk2,1,mimax,mjmax,mkmax);
	new_mat(&a,4,mimax,mjmax,mkmax);
	new_mat(&b,3,mimax,mjmax,mkmax);
	new_mat(&c,3,mimax,mjmax,mkmax);
	mat_set_init(&p);
	mat_set(&bnd,0,1.0);
	mat_set(&wrk1,0,0.0);
	mat_set(&wrk2,0,0.0);
	mat_set(&a,0,1.0);
	mat_set(&a,1,1.0);
	mat_set(&a,2,1.0);
	mat_set(&a,3,1.0/6.0);
	mat_set(&b,0,0.0);
	mat_set(&b,1,0.0);
	mat_set(&b,2,0.0);
	mat_set(&c,0,1.0);
	mat_set(&c,1,1.0);
	mat_set(&c,2,1.0);
	gosa = jacobi(nn);
	fmt.Printf("%.6f\n",gosa);


}

func new_mat (mat* Matrix, vmnums int, vmrows int, vmcols int, vmdeps int) (int) {

		mat.mnums= vmnums
		mat.mrows= vmrows
		mat.mcols= vmcols
		mat.mdeps= vmdeps
		mat.m = make([]float64,vmnums*vmrows*vmcols*vmdeps)
	if mat.m !=nil {
		return 1
	}
	return 0
}

func clear_mat (mat* Matrix) {
	if (mat.m != nil) {
		mat = nil
	}
}


func mat_set(mat* Matrix, l int, val float64) {
	for i:=0 ;i < mat.mrows ; i++ {
		for j:=0 ; j< mat.mcols ; j++ {
			for k:=0; k< mat.mdeps ; k++ {
				MR_set(mat,l,i,j,k,val)
			}
		}
	}
}

func  mat_set_init(mat* Matrix) {
	for i := 0 ; i < mat.mrows; i++ {
		for j := 0; j< mat.mcols; j++ {
			for k :=0; k< mat.mdeps; k++  {
				var val = (float64)(i*i)/
					(float64)((mat.mrows -1)* (mat.mrows -1))
				MR_set(mat, 0, i, j , k,val)
			}
		}
	}
}

func jacobi(nn int) (float64) {
	var i,j,k,n,imax,jmax,kmax int;
	var gosa float64;
	imax = p.mrows-1;
	jmax = p.mcols-1;
	kmax = p.mdeps-1;
	for n=0;n<nn;n++ {
		gosa = 0.0
		var wg sync.WaitGroup
		var gosa_ch = make(chan float64)
		for i=1;i<imax;i++ {
			wg.Add(1)
			go internal_jacobi(i,jmax ,kmax ,n ,&wg ,gosa_ch)
		}
		go func() {
			wg.Wait()
			close(gosa_ch)
		}()

		for g:= range gosa_ch {
			gosa +=g
		}
		for i=1;i<imax;i++ {
			for j=1;j<jmax;j++ {
				for k=1;k<kmax;k++ {
					MR_set(&p,0,i,j,k,(MR_get(&wrk2, 0, i, j, k)))
				}
			}
		}
	}

	return gosa;
}

func internal_jacobi(i int, jmax int, kmax int, n int, wg *sync.WaitGroup, gosa_ch chan<-float64) {
	for j:=1;j<jmax;j++ {
		wg.Add(1)
		go internal_j(i ,j ,kmax ,n ,wg , gosa_ch)
	}
	(*wg).Done()
}

func internal_j(i int, j int,kmax int,n int, wg *sync.WaitGroup, gosa_ch chan<-float64) {
	var s0,ss float64
	var gosa float64
	for k:=1;k<kmax;k++{
		s0 = MR_get(&a,0,i,j,k) * MR_get(&p,0,i+1,j,  k) +
			MR_get(&a,1,i,j,k) * MR_get(&p,0,i,  j+1,k) +
						MR_get(&a,2,i,j,k) * MR_get(&p,0,i,  j,  k+1) +
						MR_get(&b,0,i,j,k) *
						(MR_get(&p,0,i+1,j+1,k) - MR_get(&p,0,i+1,j-1,k) -
						MR_get(&p,0,i-1,j+1,k) +MR_get(&p,0,i-1,j-1,k)) +
					MR_get(&b,1,i,j,k) *
						( MR_get(&p,0,i,j+1,k+1) - MR_get(&p,0,i,j-1,k+1) -
						MR_get(&p,0,i,j+1,k-1) + MR_get(&p,0,i,j-1,k-1) ) +
						MR_get(&b,2,i,j,k) *
						( MR_get(&p,0,i+1,j,k+1) - MR_get(&p,0,i-1,j,k+1) -
						MR_get(&p,0,i+1,j,k-1) + MR_get(&p,0,i-1,j,k-1) ) +
						MR_get(&c,0,i,j,k) * MR_get(&p,0,i-1,j,  k) +
						MR_get(&c,1,i,j,k) * MR_get(&p,0,i,  j-1,k) +
						MR_get(&c,2,i,j,k) * MR_get(&p,0,i,  j,  k-1) +
						MR_get(&wrk1,0,i,j,k);
					ss = (s0*MR_get(&a,3,i,j,k) - MR_get(&p,0,i,j,k))*MR_get(&bnd,0,i,j,k);
					gosa +=ss*ss;
					// fmt.Printf("Goas : %.10f \n", gosa);
		MR_set(&wrk2,0,i,j,k,(MR_get(&p,0,i,j,k) + omega*ss));
	}
	gosa_ch<-gosa

	(*wg).Done()
}
