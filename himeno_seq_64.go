package main 
import (
	"fmt"
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
func main() {

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
	// fmt.Printf("value of p: mrows : %d, mcols: %d, mdeps: : %d\n", p.mnums,p.mcols,p.mdeps);
	// fmt.Printf("value of bnd: mrows : %d, mcols: %d, mdeps: : %d\n", a.mnums,a.mcols,a.mdeps);
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
	// gosa++ 
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
}

func new_mat (mat* Matrix, vmnums int, vmrows int, vmcols int, vmdeps int) (int) {

		mat.mnums= vmnums
		mat.mrows= vmrows
		mat.mcols= vmcols
		mat.mdeps= vmdeps
		mat.m = make([]float64,vmnums*vmrows*vmcols*vmdeps)
	
	if mat.m !=nil {
		// fmt.Printf("value of matrix: mrows : %d, mcols: %d, mdeps: : %d\n", mat.mnums,mat.mcols,mat.mdeps);
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
	// fmt.Printf("value of p: mrows : %d, mcols: %d, mdeps: : %d", mat.mrows,mat.mcols,mat.mdeps);
	for i := 0 ; i < mat.mrows; i++ {
		for j := 0; j< mat.mcols; j++ {
			for k :=0; k< mat.mdeps; k++  {
				var val = (float64)(i*i)/
					(float64)((mat.mrows -1)* (mat.mrows -1))
				MR_set(mat, 0, i, j , k,val)
				//fmt.Printf("initvalue set for p %f",MR_get(mat, 0, i, j, k));
			}
		}
	}
}


func jacobi(nn int, a* Matrix, b* Matrix,c* Matrix, p* Matrix, bnd* Matrix, wrk1* Matrix, wrk2* Matrix) (float64) {

	var i,j,k,n,imax,jmax,kmax int;
	var gosa,s0,ss float64;
	imax = p.mrows-1;
	jmax = p.mcols-1;
	kmax = p.mdeps-1;
	// fmt.Printf("inside jacobi imax : %d, jmax : %d, kmax: %d \n",imax,jmax,kmax);
	for n=0;n<nn;n++ {
		gosa = 0.0
		for i=1;i<imax;i++ {
			for j=1;j<jmax;j++ {
				for k=1;k<kmax;k++ {
					s0 = MR_get(a,0,i,j,k) * MR_get(p,0,i+1,j,  k) +
						MR_get(a,1,i,j,k) * MR_get(p,0,i,  j+1,k) +
						MR_get(a,2,i,j,k) * MR_get(p,0,i,  j,  k+1) +
						MR_get(b,0,i,j,k) *
						(MR_get(p,0,i+1,j+1,k) - MR_get(p,0,i+1,j-1,k) -
						MR_get(p,0,i-1,j+1,k) +MR_get(p,0,i-1,j-1,k)) +
						MR_get(b,1,i,j,k) *
						( MR_get(p,0,i,j+1,k+1) - MR_get(p,0,i,j-1,k+1) -
						MR_get(p,0,i,j+1,k-1) + MR_get(p,0,i,j-1,k-1) ) +
						MR_get(b,2,i,j,k) *
						( MR_get(p,0,i+1,j,k+1) - MR_get(p,0,i-1,j,k+1) -
						MR_get(p,0,i+1,j,k-1) + MR_get(p,0,i-1,j,k-1) ) +
						MR_get(c,0,i,j,k) * MR_get(p,0,i-1,j,  k) +
						MR_get(c,1,i,j,k) * MR_get(p,0,i,  j-1,k) +
						MR_get(c,2,i,j,k) * MR_get(p,0,i,  j,  k-1) +
						MR_get(wrk1,0,i,j,k);

					ss = (s0*MR_get(a,3,i,j,k) - MR_get(p,0,i,j,k))*MR_get(bnd,0,i,j,k);

					gosa += ss*ss;
					//fmt.Printf("Goas : %.10f \n", gosa);
					MR_set(wrk2,0,i,j,k,(MR_get(p,0,i,j,k) + omega*ss));

				}
			}
		}
		fmt.Printf("%.15f\n", gosa)
		for i=1;i<imax;i++ {
			for j=1;j<jmax;j++ {
				for k=1;k<kmax;k++ {
					MR_set(p,0,i,j,k,(MR_get(wrk2, 0, i, j, k)))
				}
	}
}
}
return gosa;
}
