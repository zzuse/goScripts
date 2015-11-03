package main

import (
    "C"
    "fmt"
)

//export PrintInt
func PrintInt(x int) {
    fmt.Println(x)
}

func main() {}
