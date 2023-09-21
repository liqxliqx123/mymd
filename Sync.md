##  Sync

所谓实现线程安全，就是加锁

加锁的原理

* 步骤
  * 自旋 重复进行CAS操作；性能高，CPU换时间
    * 拿到锁
    * 重新尝试
      * 一直失败
        * 加入等待队列 runtime维护队列
          * 等待被唤醒 阻塞和唤醒通过runtime实现
            * 重新竞争
* 本质 状态机的变化
* 模式
  * 正常模式
    * 排队的协程和新来的协程(此刻已经在使用CPU)进行竞争
    * 能插队，非公平锁
  * 饥饿模式
    * 如果等待时间超过1ms,那么锁将会变为饥饿模式，优先被选择

### Mutex/RwMutex

一般用法是将Mutex和RWMutex和被保护的资源封装在一个结构体内，当修改时提供相应的方法

加写锁后DoubleCheck，防止加写锁以后也出现并发不安全

```go
//SafeMap ...
type SafeMap1[K comparable, V any] struct {
	m    map[K]V
	lock sync.RWMutex
}

func (s *SafeMap1[K, V]) UpdateMapDoubleCheck(k K, v1 V) (v V, b bool) {
	s.lock.RLock()
	if v, b = s.m[k]; b {
		return
	}
	s.lock.RUnlock()
	s.lock.Lock()
	defer s.lock.Unlock()
	// double check
	if v, b = s.m[k]; b {
		return
	}
	s.m[k] = v1
	v = v1
	return 
}
```

注意事项

* 锁不可重入，注意递归等重复加锁的场景
* 尽量使用defer解锁，避免panic

### Once

确保一个动作并发时最多执行一次，用来初始化资源或单例模式

方法需要使用指针对象， 否则会使用拷贝

### Pool

如果需要考虑缓存资源，比如创建好的对象，那么可以使用sync.Pool，目的是减少内存分配和减轻GC压力(CPU)

Sync.Pool会先检查自己是否有资源， 有就返回，没有就创建新的

Sync.Pool会在GC时释放缓存的资源

内存使用量不可控

补充：

* 内存分配到栈上不需要GC管理， 分配到堆上才需要
* buffer 可以看作是字节数组， 用来缓存，如拼接字符串，IO操作等
  * 三方包 bytebufferpool
    * 基于sync.Pool的封装
    * 引入校准机制，根据使用次数动态计算defaultSize和maxSize

设计细节

* 采用PMG调度模型，P代表处理器；任何绑定到P上的数据都不需要竞争， 因为P在同一时刻只有一个G在运行
* 每个P有一个poolLocal对象， 包含private和shared
* shared 是一个链表+ring buffer的结构
  * 总体是一个双向链表
  * 每个链表的节点指向一个ring buffer，后一个节点的ring buffer是前一个节点的两倍
* ring buffer的优势（数组的优势）
  * 一次性分配内存，循环利用
  * 对缓存友好

GET步骤

* 首先查找private是否可用，可用就直接返回
* 不可用就从自己的poolChain里尝试获取一个
  * 从最近创建的ring buffer开始找
* 如果找不到尝试从其他P里偷， 窃取算法。竞争小于全局共享队列
* 偷不到时去找victim
* victim中也没有则重新创建

Pool与GC

正常情况下为了控制Pool的内存消耗，需要考虑淘汰问题。但sync.Pool完全依赖于GC, 用户无法手动控制

GC的过程

* Locals 会变成victim
* victim灰被直接回收掉，如果对象再次被使用则变回locals。
  * 防止GC引起性能抖动

```go
package demo

import (
	"fmt"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			return &User{}
		},
	}
	u1 := pool.Get().(*User)
	u1.ID = 12
	u1.Name = "Tom"
	// 一通操作
	// 放回去之前要先重置掉
	u1.Reset()
	pool.Put(u1)

	u2 := pool.Get().(*User)
	fmt.Println(u2)
}

type User struct {
	ID   uint64
	Name string
}

func (u *User) Reset() {
	u.ID = 0
	u.Name = ""
}

func (u *User) ChangeName(newName string) {
	u.Name = newName
}

```



### WaitGroup

注意在goroutine外进行Add;Done减为负数会panic

实现

* noCopy主要用于告诉编译器这个货不能用来复制，比如值传递；应该使用指针 
* state1
  * 高32位 记录任务数量
  * 低32位 记录等待的goroutine数量
* state2 信号量 用于挂起或唤醒goroutine



### channel

几个要点

* 是否带缓冲
  * 非缓冲 要求发送和接收者同时存在
  * 有缓冲
    * 满了阻塞发送者
    * 空了阻塞接收者
* 谁在发
* 谁在收
* 谁来关
* 是否关闭了

用途

* 看作是队列，用于传递数据
* 利用阻塞特性控制goroutine

实现任务池

* 使用一个channel控制并发， 任务在从channel中获取值后执行， 这是一个并发队列
* 预先创建好指定个goroutine, 每个goroutine不断从任务队列中获取

goroutine泄露

* channel一直阻塞，导致goroutine资源不释放

* 读写nil的channel不会panic会一直阻塞

* 即使channel close掉了也，也会继续读出零值

  ```go
  func TestChannelClosed(t *testing.T) {
  	ch := make(chan string, 10)
  	go func() {
  		data := <-ch
  		fmt.Printf("g1 receiver %s", data)
  	}()
  
  	go func() {
  		close(ch)
  		fmt.Printf("closed")
  
  	}()
  	for {
  		<-ch
  		fmt.Println("111")
  	}
  }
  ```

  关闭了继续写会panic

  panic: send on closed channel [recovered]

内存逃逸

* 如果用channel发送指针那么必然逃逸；编译器无法确定指针数据最终是被哪个goroutine接收

实现细节

* buf 存储数据，unsafe.Pointer 指向ring buffer
* recvq 接收者队列
* sendq 发送者队列

开源实例 kratos

* 启动过程需要的考虑
  * 监听关闭信号
  * 监控server启动过程，一个失败全部退出
  * 监控server异常退出

面试

* 用channel实现一个任务池
* 用channel控制goroutine数量
* 用channel实现生产-消费模型

### 补充：

泛型

* 结构体定义

  ```go
  type SafeMap[K comparable, V any]struct{
  	values map[K]V
    lock sync.RWMutex
  }
  ```

* 方法定义

  ```go
  func(s *SafeMap[K,V]) LoadOrStore(key K, v V)(val V, b bool){}
  ```

  