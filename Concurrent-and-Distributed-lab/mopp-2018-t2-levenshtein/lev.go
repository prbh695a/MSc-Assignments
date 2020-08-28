package main

import (
	"fmt"
	"sync"

)

func Distance(s, t string,wg *sync.WaitGroup) int {
	var (
		n = len(s)
		m = len(t)
	)

	switch {
	case n == 0:
		return m
	case m == 0:
		return n
	}

	d := buildMatrix(n, m)
	defer wg.Done()
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {

			cost := 0

			if s[i-1] != t[j-1] {
				cost = 1
			}

			d[j][i] = minimum(d[j-1][i]+1, d[j][i-1]+1, d[j-1][i-1]+cost)
		}
	}
	return d[m][n]
}

// helper function to build a matrix
func buildMatrix(n, m int) [][]int {
	matrix := make([][]int, m+1)

	for i := 0; i <= m; i++ {
		matrix[i] = append(matrix[i], make([]int, n+1)...)
		matrix[i][0] = i
	}

	for i := 0; i <= n; i++ {
		matrix[0][i] = i
	}

	return matrix
}

// helper function to determine the minimum among 3 integers
func minimum(a, b, c int) int {
	switch {
	case a < b && a < c:
		return a
	case b < a && b < c:
		return b
	}

	return c
}

func main() {
	// create two example strings, calculate the edit distance between them
  //fmt.Print("Enter text: ")
	var wg sync.WaitGroup

  var str1 string
  var str2 string
  fmt.Scanln(&str1)
  fmt.Scanln(&str2)
	wg.Add(1);
	fmt.Println(Distance(str1, str2,&wg))
}
