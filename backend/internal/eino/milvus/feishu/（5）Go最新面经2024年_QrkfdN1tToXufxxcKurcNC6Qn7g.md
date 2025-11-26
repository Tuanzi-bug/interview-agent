# （5）Go最新面经2024年

# 北京-度小满

自我介绍

1、介绍项目具体业务

2、项目碰到的技术难点

<!-- 未处理的块类型: 34, block_id: HhiLdntUQobIyExxNL0cKH5inLd -->

微服务第一次接触没有底，go-zero、looklook学习，looklook架构用到项目上，问题：浪费资源、替代方案

3、日志收集用来干什么？

<!-- 未处理的块类型: 34, block_id: IYoEdWP2mok1NSxnzaocGxlcnfb -->

微服务链路跟踪

4、链路最长有多少个服务？

<!-- 未处理的块类型: 34, block_id: SA6WdTmcSon5SjxXVGLcOY4onbe -->

3个

5、日志收集的意义在哪里，服务很少，没有必要做trace这些东西呀？

<!-- 未处理的块类型: 34, block_id: QxNOdueKBo4PPPxFBQZc0vv3n9e -->

项目后面需要引入AI，还有加入其他功能服务，提前准备好

6、redis持久化机制

<!-- 未处理的块类型: 34, block_id: RH4RdbOrjoW01BxKV8bc8vMGnNc -->

AOF和RDB；两个持久化的特点介绍

7、重新介绍一下AOF写入的三种方式

8、了解过redis集群吗？

<!-- 未处理的块类型: 34, block_id: Q0qvdJRgZomJfixvlMEcPiWAnbc -->

没有，运维人员做的工作个人觉得没必要看

9、redis从客户端执行命令到最终命令的响应，这中间经历了那些过程，能大概描述一下吗

10、redis速度为什么那么快？

<!-- 未处理的块类型: 34, block_id: RZyTdngsaobFFWxX4D2c4tYxnqh -->

基于内存：极高的读写速度，特别对于简单的存取操作，执行时间非常短，主要耗时在于网络IO

单线程：必要上下文切换，锁竞争  

IO多路服用：如epoll,能够在一个线程中高效地处理多个客户端连接  >

高效的数据结构：全局哈希、压缩表、跳跃表

11、使用过那些redis的数据类型？

<!-- 未处理的块类型: 34, block_id: S2JfdtYzZoPhZpxbTy6ctiiwnqf -->

redis的常用数据类型，主要使用string  

12、redis的过期清理策略

<!-- 未处理的块类型: 34, block_id: BGO6dkwHKori30xl89tc1KSOnf7 -->

惰性删除和定期删除：详细介绍...  加强

13、说一下innodb数据存储的结构?

<!-- 未处理的块类型: 34, block_id: ZzmRdE8pMovjbTxjVW2cTeG7nMc -->

B+树，B+树的特性

14、MySQL怎么做异常恢复的，MySQL挂了，重新启动的时候怎么做异常恢复？

<!-- 未处理的块类型: 34, block_id: PqmJdJmTvoRRaIxrzQocG1fQnMb -->

redo-log日志做恢复，redo-log主要是记录写操作，通过里面的写操作记录恢复  两次提交的状态   加强

15、聚簇索引和非聚簇索引的区别？

16、建索引的时候有哪些需要注意的点？

17、你常用Linux命令有哪些？

mv、cd、ls、vim、ps 

18、平时使用kafka的时候，出现消息堆积一般是怎么处理的？

19、TCP和UDP的区别？

20、TCP怎么处理拥塞控制的？

21、算法：删除链表倒数第N个节点

# 广州-没有提供公司名称

<!-- 未处理的块类型: 34, block_id: PN3gdyFmpoZf1MxI8M5cYSmOnNf -->

主要问Linux运维，排查问题的多

介绍一下Linux，IO多路复用

Linux的文件描述符

常用的Linux指令、vim指令

1. MySQL的数据结构？
<!-- 未处理的块类型: 34, block_id: DcPmduKLlotF0gxzv8VcTVeUnZf -->

B+树，回答它的特点

1. B+树和B树的区别？
1. TCP和UDP的区别
1. 三次握手和四次挥手的过程
# 杭州-玩心不止玩网络科技有限公司

**项目问答**

1. 介绍一下锁？
1. CAS是什么？
1. 自旋的意义是什么？
1. golang怎么判断对象是分配到堆上还是栈上？
1. 发生内存泄漏怎么排查？
1. 介绍一下GMP模型？
1. 介绍一下GC？
1. 如果A对象和B对象相互引用，会被GC吗？为什么？
1. 假设需要请求第三方接口，而第三方接口不太稳定，你会怎么设计？
1. MySQL的数据结构是什么？
1. B+树和B树的区别？
1. Redis的IO复用
