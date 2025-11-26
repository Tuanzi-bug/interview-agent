# （1）Go最新面经2024年

# b站外包

1. 旋转数组
1. mysql索引相关。
- B+树有什么特点？
- 为什么不用B树（查询的速度差不多，因为b+树数据都在叶子节点）。
- 非聚簇索引和聚簇索引的区别。
- 索引为什么要用id不用字符
1. Linux相关命令和场景
1. docker和k8s
# 上海莹锴网络科技

1. Kafka怎么保证消息不丢失
1. Kafka里面生产者给broker应答给生产者的时候网络断了该怎么处理？
1. 讲一下golang中的并发编程
# 梭翱信息技术

1. go语言特性（channel、map考察）
1. 讲一下waitgroup的使用
1. 知道模块化缓存吗？
1. 知道分级缓存吗？
<!-- 未处理的块类型: 34, block_id: ZlHXdrkK5of7lRxgYoQc7KQ9nNc -->

主要是一些缓存的应用场景为主。对于一些技术（redis、mq）的应用场景这部分比较欠缺

  

# 矢安科技

1. 哪些数据结构是线程不安全的
1. Map为什么是线程不安全的
1. Channel阻塞可以实现什么场景（计数，令牌桶）
1. Mysql什么时候是行锁什么时候是表锁
1. Mysql有几种错误读（脏读、幻读等等）
1. Mysql默认事务隔离级别是什么
1. 假如有个sql联了多个表还有字查询，改怎么优化
1. 你平常是怎么优化mysql的
1. Kafka为什么快？
1. Kafka怎么实现消息不丢失
1. Kafka是顺序写还是随机写
1. Go协程是怎么扩展内存的（找P要）
1. 讲一下你对docker和k8s的理解
1. 说一下某种集群的leader选举策略（举例了redis）
1. Redis中什么是主管下线，什么是客观下线。
1. 聊一聊你对GRPC的理解。
1. 为什么gprc传输比JSON快
- 少了json转为2进制
- protobuf文件中字段名用后面的数字代替，进一步减少数据量（学到了）
**场景题：**

1. 假如有大量定时任务需要在凌晨1点都准时开始执行，你会怎么做
1. 假如消息发送过多导致大量堆积怎么处理
# 华苏科技（国家电网外包）

**项目拷打（介绍项目+遇到的难题）**

1. 有缓冲channel和无缓冲channel的区别
1. 了解gin的中间件吗，讲一下你对他的了解
1. select 满足多个case的时候怎么执行的。
1. 如果有一个全局变量怎么保证并发安全。
1. CPU高问题如何解决？
1. 知道哪些设计模式？
**场景题：**

我有一个方法，用来存储一些文件资源，有多种不同的存储方式，你会怎么设计这个方法（应该是要考察泛型的使用）

# 爱可生

<!-- 未处理的块类型: 34, block_id: FrhadqMkmoHW2GxPOTRcWPdXnGc -->

面试官说是go开发，但是没有什么技术原理提问。

**介绍简历中的项目。**

1. 遇到的项目场景难题。
<!-- 未处理的块类型: 34, block_id: NkU7dUwiTo71vGx5GcjcG8fEnKh -->

（他不太想听那种用技术选型方案来解决的常规问题，吹了一下systemtap）

1. 讲一个技术栈中随便一个技术遇到的难题。
1. 平时是怎么学习的。
<!-- 未处理的块类型: 34, block_id: Yd5OdCmWAovD1rxNtERcxp7ynFb -->

整个面试几乎就没有技术性提问，一直在让我介绍项目，和遇到的问题以及我是怎么解决的，解决的思路是什么。

# 特斯拉外包笔试

**题1**

<!-- 未处理的块类型: 27, block_id: WirqdoxbWo3ECqxOYJnc4azXnH5 -->

**题2**

<!-- 未处理的块类型: 27, block_id: QnKbdPRTFogFsoxyiWwcDRbBn1g -->

<!-- 未处理的块类型: 27, block_id: BwFldigVLosHk9xT03PcMrPFnZe -->

**题3**

<!-- 未处理的块类型: 27, block_id: YBQzdbuncohTbxx3fLbco3LLnXf -->

<!-- 未处理的块类型: 27, block_id: TkgQdMKYXoNHcMxTTmocfduWn3f -->

<!-- 未处理的块类型: 34, block_id: QWxTdaL73o4U4Fxfr3Rcc3YGn8b -->

最后一道sql没写出来。但是前两题自测都对

# 成都美大

**项目拷打**

1. 讲一下mysql的索引是什么结构
1. 讲一下sql一般是怎么优化的
1. Kafka消息堆积怎么处理
1. 写一个方法的时候是传值好还是结构体好
**场景题：**

秒杀超卖怎么解决。（分布式锁+redis缓存）

# 矢安科技二面

<!-- 未处理的块类型: 34, block_id: Ta6Ad7X7qoXp1VxX6J2c6Lbdnyf -->

一面的技术leader，没有聊太多的技术话题。主要是一些团队协作沟通上的问题

1. 假如产品给你一个需求，你觉得不合适，和产品经理有冲突，你会怎么做？
1. 你平时是怎么学习的？
1. 假如给你一个活要求某个时间内快速完成，你又没学过，你会怎么做?
<!-- 未处理的块类型: 34, block_id: JbYxd3vPPoJ38wxLBGTcWZLrnyh -->

试探你是不是愿意加班

# 七云网络

## 笔试

**三道程序解答题：**

**题1**

<!-- 未处理的块类型: 27, block_id: IGiQdrdU0oLBaNxxiQAcqCDOn3d -->

**问1**

先输出哪个？

**题2**

<!-- 未处理的块类型: 27, block_id: MhCudltqtoA0bVxbBKIcYofKnSf -->

**问2**

输出什么？

**题3**

<!-- 未处理的块类型: 27, block_id: FgM5dlhgSo8q6fxVmhNcga2EnGf -->

**问3**

这段代码有什么问题？

**问答题：**

1. TCP和http的关系是什么？
1. 伪代码描述一下乐观锁
1. Linux怎么看磁盘占用？
1. 描述一下GC的过程？
1. SQL题：写出薪资第二高的薪资
表emp

id int

salary int

**算法：**

力扣 ： 20. 有效的括号（纸上纯手写）

## 面试

<!-- 未处理的块类型: 34, block_id: XcModOv06omowmxEA4HcE8S1ncc -->

挨个问笔试的问题。。。

口述了两个方法for循环里操作channel之类的，但是他语言组织的我实在没听懂。。。

1. 假如有一个高并发的场景，我怎么处理（不能借住其他组件，纯go程序）
<!-- 未处理的块类型: 34, block_id: GMnYd7Hhmoa5FTxenC0cU1jHnZf -->

然后mq问了两个迷一样的问题。

1. 他提到了Kafka然后问我用的什么MQ，我说Kafka就是一种mq啊。感觉面试官不是很熟悉Kafka
1. 接着他问我Kafka里的分组是怎么设置的，我以为他问的是消费者分组。结果他说是topic里的。结果他问的分组是Topic分区。。。
1. 这家就别去了。面试流程很不合理，面试官沟通起来比较费劲，也不是很专业。
# 杉岩数据

**项目拷打**

<!-- 未处理的块类型: 34, block_id: JM1rd9zXLoCSWcxxD3ocqbgAnCd -->

其中问了为什么不用普罗米修斯去监控

1. 假如你用于通知的Kafka挂了怎么办？有没有对Kafka进行监控？
1. channel的使用场景？
1. chaneel关闭之后再读和再关闭会发生什么？
1. map中的数据delete之后内存会回收吗？
1. GRPC请求和http请求有什么区别
# 腾娱

**四道基础语法题**

```plaintext
c := []int{11, 12, 13}
test(c)
log.Info("c=%v", c)
func test(s []int) {
        for i := 0; i < 10; i++ {
                s = append(s, i)
        }
}
```

1. c最后是怎么样的?
<!-- 未处理的块类型: 22, block_id: JpCvdcgRYotypVxFL4Rcwz8Fnng -->

```plaintext
func main() {
        values := []int{1,2,3,4,5,6,7,8,9}
        for _,v := range values {
                go func(){
                        println(v)
                }()
        }
}
```

1. 求输出
<!-- 未处理的块类型: 22, block_id: UzpudEG8eoP4aFxlIKacY46in8d -->

```plaintext
func main() {
    wg := sync.WaitGroup{}

    for i := 0; i < 5; i++ {
        go func(wg sync.WaitGroup, i int) {
            wg.Add(1)
            fmt.Printf("i:%d\n", i)
            wg.Done()
        }(wg, i) 
    }
    wg.Wait()
        println("exit")
}
```

1. 求输出
<!-- 未处理的块类型: 22, block_id: RZmZdTHpjoAzJFxSYEWcToxcnod -->

```plaintext
func testDefer() (err error) {
        defer func() {
                if err != nil {
                        log.Error("defer: %s", err)
                }
        }()
        log.Info("testDefer: %s", "test")
        return handle()
}

func handle()error{
        return fmt.Errorf("normal:test")
}
```

1. 求输出
<!-- 未处理的块类型: 22, block_id: UqPqd2UXTo4COixeB7McpV1Hneg -->

1. 改造他让他变得有序
```plaintext
func main() {
        values := []int{1,2,3,4,5,6,7,8,9}
        for _,v := range values {
                go func(){
                        println(v)
                }()
        }
}
```

# 杉岩二面

**项目拷打15分钟**

1. Kafka的消息丢失和消息重复消费。
1. Kafka和Rabbitmq的区别在哪？（架构、推和拉）
1. 拉的模式有什么好处（控制消费速度）
1. 使用分布式锁的过程中应用挂了？
1. 优雅启停+defer
1. 使用过期时间+自动续期
1. 对象存储和文件存储的主要区别是什么？
1. 分片上传是怎么实现的（文件合并hash一致性校验，引出文件秒传）
1. 邮箱验证码功能怎么实现的。（redis+邮箱组件）
1. jwt的格式。加密算法、内容、过期时间
1. 讲一下defer的原理
1. 讲一下map的底层结构
1. map中hash冲突怎么解决（链表、红黑树）
1. 讲一下go性能调优的案例（pprof，线程日志）
1. 通过线程日志延伸出，怎么看一个线程在线程日志里是卡在循环还是事件等待？
1. 线程日志上面会有标记。
1. 讲一下mysql的事务隔离级别？
1. 解释一下什么是可重复读？
1. 事务实现的底层原理？
1. Redis持久化机制（RDB，AOF）
1. 为什么持久化的时候是fork子进程处理
1. 讲一下docker实现容器的基本原理
1. 用过其他容器运行时吗
1. K8s有哪些组件？
1. 你们是用什么control去构建deployment的？（没听懂）
<!-- 未处理的块类型: 22, block_id: Hu0AdJP6QodlToxUP6wc1omfnvr -->

