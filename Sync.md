##  Sync

所谓实现线程安全，就是加锁

加锁的原理

* 步骤
  * 自旋 重复进行CAS操作；性能高，CPU换时间
    * 拿到锁
    * 重新尝试
      * 一直失败
        * 加入等待队列 runtime维护队列
          * 等待被唤醒
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

如果需要考虑缓存资源，比如创建好的对象，那么可以使用sync.Pool，目的是减少内存分配和减轻GC压力

Sync.Pool会先检查自己是否有资源， 有就返回，没有就创建新的

Sync.Pool会在GC时释放缓存的资源

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



WaitGroup

补充：

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

  