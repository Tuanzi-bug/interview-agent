# Go语言的map删除一个key，内存是否会立即释放？为什么？写源码验证一下。

**Go语言的map删除一个key，内存是否会立即释放？为什么？写源码验证一下。**

在Go语言中，删除一个map中的key并不会立即释放相关的内存。这是因为Go语言的垃圾回收机制是基于标记-清除算法的，并不会即时回收被删除的map元素所占用的内存。

当删除map中的一个key时，该key对应的值会被标记为不可达，但并不会立即从内存中移除。Go的垃圾回收器会定期运行，通过标记所有可达的对象，然后清除所有不可达的对象来回收内存。**这意味着被删除的key所占用的内存会在垃圾回收器运行时被释放。**

我写了验证代码，但是并没有支撑起我的想法：

<!-- 未处理的块类型: 31, block_id: doxcnTLyIHaHEO9OJGaULP9oVTf -->

<!-- 未处理的块类型: 32, block_id: doxcnpbjJcUNDGeAihcEljmJ99b -->

package main

import (
        "fmt"
        "runtime"
)

func main() {
        m := make(map[int]int)

        for i := 0; i < 1000000; i++ {
                m[i] = i
        }

        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)
        heapAlloc := memStats.HeapAlloc
        fmt.Println("Before deletion:", heapAlloc)

        delete(m, 0)

        var memStatsAfter runtime.MemStats
        runtime.ReadMemStats(&memStatsAfter)
        heapAllocAfter := memStatsAfter.HeapAlloc
        fmt.Println("After deletion:",heapAllocAfter)
}

运行结果如下:

<!-- 未处理的块类型: 27, block_id: doxcnc6hYWhaHMj7aAOX9QSKVud -->

在上述代码中，我们创建了一个包含1000000个元素的map，并在删除第一个元素后打印了HeapAlloc的值。

我的猜想是发现删除元素后HeapAlloc的值并没有立即减少，而是在垃圾回收器运行后才会释放相关的内存。

但是实际的运行结果并不是这样，因为Go语言的垃圾回收器是自动管理内存的，程序员无需手动释放内存。我没有办法确定GC什么时候回收的内存。

**关于这个问题大家有什么想法，欢迎交流讨论。**

## **延伸一下：Go语言中runtime.MemStats.HeapAlloc的作用是什么？**

在Go语言中，runtime.MemStats.HeapAlloc是runtime包中的一个结构体MemStats的字段，用于获取当前程序堆上已分配的内存大小。

具体来说，HeapAlloc表示当前程序堆上已分配的内存字节数。堆是用于动态分配内存的区域，其中包含了程序运行时动态分配的对象。HeapAlloc的值可以用于监控和诊断程序的内存使用情况。

通过访问runtime.MemStats.HeapAlloc，我们可以获取当前程序在堆上已分配的内存大小。这对于了解程序的内存占用情况、进行性能优化和内存泄漏排查等都非常有用。

需要注意的是，runtime.MemStats结构体中还有其他字段，如HeapSys表示堆的总大小，HeapIdle表示空闲的堆内存大小，HeapReleased表示已释放的堆内存大小等，这些字段可以提供更全面的内存使用信息。

以下是一个简单的示例，演示了如何使用runtime.MemStats.HeapAlloc获取当前程序堆上已分配的内存大小：

<!-- 未处理的块类型: 31, block_id: doxcnu4GupvjGs4CU3o5H8Gnz5f -->

<!-- 未处理的块类型: 32, block_id: doxcnRxAc4yrdHW0EtZhghFjBfe -->

package main

import (
        "fmt"
        "runtime"
)

func main() {
        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)

        fmt.Println("HeapAlloc:", memStats.HeapAlloc)
}

在上述代码中，我们通过runtime.ReadMemStats函数获取当前的内存统计信息，并打印了HeapAlloc的值，即当前程序堆上已分配的内存大小。

