package main

import (
    "fmt"
    "runtime"
    "sync"
)
const INF = 4294967295

func mdAllPairs(dists []uint32,v uint32,start uint32 , end uint32, k uint32, wg * sync.WaitGroup) {
    go func() {
        var j, i uint32;
        (*wg).Done()
        x := k*v
        for i = start;i < end;i++{
            A := dists[i * v + k]
            if A != INF {
                for j = 0; j < v; j++ {
                    B := dists[x+j]
                    if B != INF {
                        var intermediary= A + B;
                        C := dists[i*v+j]
                        if (intermediary < C) {
                            dists[i*v+j] = intermediary
                        }

                    }
                }
            }

        }
    }()
}

func amd(dists []uint32, v uint32) {
    var i, j uint32;
    var smd uint32 // sum of minimum distances
    var paths uint32 //number of paths
    var solution uint32

    for i = 0;i < v;i++{
        for j = 0;j < v;j++{
            A := dists[i * v + j]
            if ((i != j) && (A < INF)) {
                smd += A;
                paths++;
            }
        }
    }
    solution = smd / paths;
    fmt.Printf("%d\n", solution);

}

func debug(dists[] uint32, v uint32) {

    var i, j uint32;
    var infinity = v * v;

    for i = 0;i < v;i++{
        for j = 0;j < v;j++{
            if (dists[i * v + j] > infinity) {
                fmt.Printf("%7s", "inf");

            } else {
                fmt.Printf("%d", dists[i * v + j]);
            }
        }
        fmt.Print("\n");
    }
}


//Main program - reads input, calls FW, shows output
func main() {
    //Read input
    //First line : v(number of vertices)  and e (number of edges)
    var v, e, i,k uint32;
    var wg sync.WaitGroup

    fmt.Scanf("%d %d", & v, & e)

    dists := make([] uint32, v * v);
    for i = 0;i < v * v;i++{
        dists[i] = INF
    }
    for i = 0;i < v;i++{
        dists[i * v + i] = 0;
    }
    var source, dest, cost uint32;
    for i = 0;i < e;i++{
        fmt.Scanf("%d %d %d", & source, & dest, & cost);
        if cost < dists[source * v + dest] {
            dists[source * v + dest] = cost;
        }
    }


    var cpus = runtime.NumCPU()
    //runtime.GOMAXPROCS(cpus)
    var block = v / uint32(cpus)
    for k = 0;k < v;k++ {
        for i = 0;i <v-block;i+=block {
            wg.Add(1); //synchronize like barrier
            go mdAllPairs(dists,v,i,i+block,k, & wg)
        }
        wg.Add(1)
        go mdAllPairs(dists,v,i,v,k, & wg)
        wg.Wait()
    }

    amd(dists, v);


}
