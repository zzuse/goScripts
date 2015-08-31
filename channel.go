package main

import (
    "fmt"
)

func sum(a []int, c chan int) {
    total := 0
    fmt.Println("sum")
    for _, v := range a {
        fmt.Println(v)
        total += v
    }
    c <- total
}

func main() {
    a := make([]int, 100, 1000)
    //a[8] = 10
    for i := 0; i < 100; i++ {
        a[i] = i
    }
    c := make(chan int)

    go sum(a[:len(a)/2], c)
    go sum(a[len(a)/2:], c)
    x, y := <-c, <-c
    fmt.Println(x, y, x+y)
    //fmt.Print(a)
}
