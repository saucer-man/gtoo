package port

// port tcp scan

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

func loop(inport chan int, startport, endport int) {
	for i := startport; i < endport; i++ {
		inport <- i
	}
	close(inport)
}

type ScanSafeCount struct {
	// 结构体
	count int
	mux   sync.Mutex
}

var scanCount ScanSafeCount

func scanner(inport int, outport chan int, host *net.IPAddr, endport int) {
	// 扫描函数

	in := inport // 定义要扫描的端口号
	// fmt.Printf(" %d ", in) // 输出扫描的端口
	addr := fmt.Sprintf("%s:%d", host.IP.String(), in)        // 类似（ip,port）  + ""
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second) //建立tcp连接
	if err != nil {
		// tcp连接失败
		// fmt.Printf("\n *************( %d 不可以 )*****************\n", in)
		// fmt.Printf("\n%v\n", err)
		outport <- 0
	} else {
		// tcp连接成功
		outport <- in // 将端口写入outport信号
		fmt.Printf("\n *************( %d 可以 )*****************\n", in)
		conn.Close()
	}

	// 线程锁
	scanCount.mux.Lock()
	scanCount.count = scanCount.count - 1
	if scanCount.count <= 0 {
		close(outport)
	}
	scanCount.mux.Unlock()

}

func PortScan(host *net.IPAddr, StartPort int, EndPort int) {
	// 设置最大可使用的cpu核数
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 定义变量
	inport := make(chan int) // 信号变量，类似python中的queue
	outport := make(chan int)
	collect := []int{} // 定义一个切片变量，类似python中的list
	s_time := time.Now().Unix()
	fmt.Println("扫描开始：") // 获取当前时间
	// 定义scanCount变量为ScanSafeCount结构体，即计算扫描的端口数量
	scanCount = ScanSafeCount{count: (EndPort - StartPort)}
	fmt.Printf("扫描 %v：%d----------%d\n", host, StartPort, EndPort)
	go loop(inport, StartPort, EndPort) // 执行loop函数将端口写入input信号
	for v := range inport {
		// 开始循环input
		go scanner(v, outport, host, EndPort)
	}
	// 输出结果
	for port := range outport {
		if port != 0 {
			collect = append(collect, port)
		}
	}

	fmt.Println("--")
	fmt.Println(collect)
	e_time := time.Now().Unix()
	fmt.Println("扫描时间:", e_time-s_time)
}
