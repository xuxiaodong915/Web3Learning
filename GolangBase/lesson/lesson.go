package main

import (
    "fmt"
	"time"
)

type A struct {
    i int
}

// 定义方法
func (a *A) add(v int) int {
    a.i += v
    return a.i
}

// 声明函数变量
var function func(int) int

func main() {
	
    go func() {
        fmt.Println("run goroutine in closure")

    }()
	
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("任务%d 执行时间: %v\n", 1, elapsed)
}	