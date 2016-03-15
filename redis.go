package main
import (
    "C"
    "fmt"
    "github.com/garyburd/redigo/redis"
)

func newPool() *redis.Pool {
return &redis.Pool{
            MaxIdle: 80,
            MaxActive: 12000, // max number of connections
            Dial: func() (redis.Conn, error) {
                    c, err := redis.Dial("tcp", ":6379")
                    if err != nil {
                            panic(err.Error())
                    }
                    return c, err
            },
    } 
}

var pool = newPool()

func Output() {
        c := pool.Get()
        defer c.Close()

        test,_:=c.Do("GET", "test1")
        fmt.Println(test)
}

func main() {}
