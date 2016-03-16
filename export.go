package main

/*
#include <stdio.h>

typedef struct {
    char *p;
    int l;
} S;

extern S GoFn(void);

static void F() {
    S s;
    s = GoFn();
    printf("%p %d\n", s.p, s.l);
}
*/
import "C"

//export GoFn
func GoFn() C.S {
    return C.S{p:(*C.char)(C.malloc(1)), l:1}
}

func main() {
    C.F()
}
