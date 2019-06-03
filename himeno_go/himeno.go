package main 
import (
	"fmt"
	"log"
	"runtime/pprof"
	"flag"
	"os"
	"runtime"
	"sync"
)


type Matrix struct {
	m []float32
	mnums int
	mrows int
	mcols int
	mdeps int
}


func MR_set (mat* Matrix, n int, r int, c int, d int, val float32) {
	index :=  (n)*mat.mrows*mat.mcols*mat.mdeps +
		(r)*mat.mcols*mat.mdeps + (c)*mat.mdeps + (d)
	mat.m[index] = val
}

func MR_get(mat* Matrix, n int, r int, c int, d int) (float32) {
		index :=  (n)*mat.mrows*mat.mcols*mat.mdeps +
			(r)*mat.mcols*mat.mdeps + (c)*mat.mdeps + (d)
		return mat.m[index]

}

var omega float32 = 0.8
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


	var nn int;
	var mimax,mjmax,mkmax int;
	var msize [3]int;
	var gosa float32;

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
	//fmt.Printf("value of p: mrows : %d, mcols: %d, mdeps: : %d\n", p.mnums,p.mcols,p.mdeps);
	//fmt.Printf("value of bnd: mrows : %d, mcols: %d, mdeps: : %d\n", a.mnums,a.mcols,a.mdeps);
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

	/*
   *    Start measuring
   */
	gosa = jacobi(nn,&a,&b,&c,&p,&bnd,&wrk1,&wrk2);

	fmt.Printf("%.6f\n",gosa);

	/*
   *   Matrix free
   */
	//clear_mat(&p);
	//clear_mat(&bnd);
	//clear_mat(&wrk1);
	//clear_mat(&wrk2);
	//clear_mat(&a);
	//clear_mat(&b);
	//clear_mat(&c);

	if *memprofile != "" {
		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}

func new_mat (mat* Matrix, vmnums int, vmrows int, vmcols int, vmdeps int) (int) {

		mat.mnums= vmnums
		mat.mrows= vmrows
		mat.mcols= vmcols
		mat.mdeps= vmdeps
		mat.m = make([]float32,vmnums*vmrows*vmcols*vmdeps)
	if mat.m !=nil {
		//fmt.Printf("value of matrix: mrows : %d, mcols: %d, mdeps: : %d\n", mat.mnums,mat.mcols,mat.mdeps);
		return 1
	}
	return 0
}

func clear_mat (mat* Matrix) {
	if (mat.m != nil) {
		mat = nil
	}
}


func mat_set(mat* Matrix, l int, val float32) {
	for i:=0 ;i < mat.mrows ; i++ {
		for j:=0 ; j< mat.mcols ; j++ {
			for k:=0; k< mat.mdeps ; k++ {
				MR_set(mat,l,i,j,k,val)
			}
		}
	}
}

func  mat_set_init(mat* Matrix) {
	//fmt.Printf("value of p: mrows : %d, mcols: %d, mdeps: : %d", mat.mrows,mat.mcols,mat.mdeps);
	for i := 0 ; i < mat.mrows; i++ {
		for j := 0; j< mat.mcols; j++ {
			for k :=0; k< mat.mdeps; k++  {
				var val = (float32)(i*i)/
					(float32)((mat.mrows -1)* (mat.mrows -1))
				MR_set(mat, 0, i, j , k,val)
				//fmt.Printf("initvalue set for p %f",MR_get(mat, 0, i, j, k));
			}
		}
	}
}


func jacobi(nn int, a* Matrix, b* Matrix,c* Matrix, p* Matrix, bnd* Matrix, wrk1* Matrix, wrk2* Matrix) (float32) {
	var i,j,k,n,imax,jmax,kmax int;
	var gosa float32;
	imax = p.mrows-1;
	jmax = p.mcols-1;
	kmax = p.mdeps-1;
	//fmt.Printf("inside jacobi imax : %d, jmax : %d, kmax: %d \n",imax,jmax,kmax);
	//200
	for n=0;n<nn;n++ {
		gosa = 0.0
		//256
		for i=1;i<imax;i++ {
			//256
			for j=1;j<jmax;j++ {
				//512
				var wg sync.WaitGroup
				var temp_gosa = make(chan float32)
				var div = 100
				for k=1;k<kmax -div ;k = k + div{
					wg.Add(1)
					go internal_k(temp_gosa,i,j,k, k+50)
							}
				wg.Add(1)
				go internal_k(temp_gosa, i, j, k,kmax)
				gosa = gosa + <-temp_gosa
				fmt.Printf("no of goroutines %d \n", runtime.NumGoroutine())
				fmt.Printf("value of gosa %f \n", gosa)
				wg.Wait()
			}
		}

		for i=1;i<imax;i++ {
			for j=1;j<jmax;j++ {
				for k=1;k<kmax;k++ {
					MR_set(p,0,i,j,k,(MR_get(wrk2, 0, i, j, k)))
				}

			}

		}
fmt.Printf("gosa value after %d is %f\n", n,gosa)
	}
return gosa;
}


func internal_k (temp_gosa chan float32, i int, j int, kstart int, kmax int) {
	//fmt.Printf("thread launched %d \n", k)
	var ss, s0 float32;
	for k:=kstart; k < kmax; k++ {
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
		temp_gosa <- ss*ss;
		//fmt.Printf("Goas : %.10f \n", gosa);
		MR_set(&wrk2,0,i,j,k,(MR_get(&p,0,i,j,k) + omega*ss));


	}

}
