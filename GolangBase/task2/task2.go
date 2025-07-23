package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ✅指针
// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
// 然后在主函数中调用该函数并输出修改后的值。
func PointerFunc(p *int) {
	*p += 10
}

// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func SliceFunc(p *[]int) {
	for i := 0; i < len(*p); i++ {
		(*p)[i] *= 2
	}
}

// ✅Goroutine
// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。

// 打印奇数
func printOdd() {
	for i := 1; i <= 10; i += 2 {
		fmt.Println("奇数：", i)
	}
}

// 打印偶数
func printEven() {
	for i := 2; i <= 10; i += 2 {
		fmt.Println("偶数：", i)
	}
}

// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 定义任务类型
// 每个任务是一个无参无返回值的函数
type Task func()

// 任务调度器，接收任务列表并并发执行，统计每个任务的耗时
func RunTasks(tasks []Task) {
	var wg sync.WaitGroup
	for idx, task := range tasks {
		wg.Add(1)
		go func(i int, t Task) {
			defer wg.Done()
			start := time.Now()
			t()
			elapsed := time.Since(start)
			fmt.Printf("任务%d 执行时间: %v\n", i+1, elapsed)
		}(idx, task)
	}
	wg.Wait()
}

// ✅面向对象
// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，
// 并调用它们的 Area() 和 Perimeter() 方法。

// 定义 Shape 接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 定义 Rectangle 结构体
type Rectangle struct {
	Width, Height float64
}

// 实现 Shape 接口
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 周长
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 定义 Circle 结构体
type Circle struct {
	Radius float64
}

// 实现 Shape 接口，面积
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// 周长
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
// 组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 定义 Person 结构体
type Person struct {
	Name string
	Age  int
}

// 定义 Employee 结构体，组合 Person
type Employee struct {
	Person     // 匿名嵌入，实现组合
	EmployeeID int
}

// 为 Employee 实现 PrintInfo 方法
func (e Employee) PrintInfo() {
	fmt.Printf("姓名: %s, 年龄: %d, 工号: %d\n", e.Name, e.Age, e.EmployeeID)
}

// ✅Channel
// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func ChannelFunc(ch chan int) {
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。


// ✅锁机制
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func MutexCounterDemo() {
	var count int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println("最终计数器的值:", count)
}

// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。

func AtomicCounterDemo() {
	var count int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt32(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("最终计数器的值:", count)
}

func main() {
	// 指针1
	var v int = 9
	var p1 *int = &v
	PointerFunc(p1)
	fmt.Println(*p1)
	fmt.Println("=========")
	// 指针2
	SliceFunc(&[]int{1, 2, 3, 4, 5})

	// Goroutine1
	go printOdd()
	go printEven()
	time.Sleep(100 * time.Millisecond)

	// Goroutine2
	tasks := []Task{
		func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("任务1完成")
		},
		func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Println("任务2完成")
		},
		func() {
			time.Sleep(700 * time.Millisecond)
			fmt.Println("任务3完成")
		},
	}
	RunTasks(tasks)

	//面向对象1
	r := Rectangle{Width: 3, Height: 4}
	c := Circle{Radius: 5}

	fmt.Printf("矩形: 面积=%.2f, 周长=%.2f\n", r.Area(), r.Perimeter())
	fmt.Printf("圆形: 面积=%.2f, 周长=%.2f\n", c.Area(), c.Perimeter())
	// 面向对象2
	e := Employee{
		Person: Person{
			Name: "张三",
			Age:  28,
		},
		EmployeeID: 1001,
	}
	e.PrintInfo() // 输出：姓名: 张三, 年龄: 28, 工号: 1001

	//Channel1
	ch := make(chan int)
	ChannelFunc(ch)
	//Channel2
	ch2 := make(chan int, 100)
	ChannelFunc(ch2)
	//锁机制1
	MutexCounterDemo()
	//锁机制2
	AtomicCounterDemo()

	
}
