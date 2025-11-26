# （6）Go最新面经2024年

# 深圳 道新科技 一面

**技术面**

1. 一开始就是聊短剧私域项目
<!-- 未处理的块类型: 34, block_id: XfUCdr2KeoaS9AxMwopcID9Mnof -->

主要说了消息推送服务的方案迭代，这块的难点以及怎么解决的；分布式锁解决重复上传素材的问题

1. 聊前进链区块链项目，有什么难点，怎么优化默克尔哈希计算速度的
<!-- 未处理的块类型: 34, block_id: WNLNdefzzoSPAHxPiU6csXVjnqM -->

感觉面试官以前也做过区块链

1. 等保三级你们怎么做的
<!-- 未处理的块类型: 34, block_id: Utgodhg3LoGl8oxsbIpc5uzunxY -->

回答了接口安全方面的，运维层面做IP白名单等；

1. 聊深信服的项目，有什么难点
<!-- 未处理的块类型: 34, block_id: S3vUd6vjXoyatuxHwTfcnZQAnYg -->

面试官感觉我这块的业务很简单，问我有没有接触过底层虚拟化的事情，我说我们的团队比较大，功能分的很细，这是另一个部门负责的；我补充说深信服是瀑布开发模式，一个需求可能会做一两个月

1. gorm的坑遇到过吗？遇到过零值更新问题吗？怎么强制更新
<!-- 未处理的块类型: 34, block_id: QuzFdQde4oaamAxPxgxc1Lv6nob -->

我说用map去更新的话是可以避免的，说的有点模糊，确实有段时间没用gorm了

1. 开始问go的问题，切片和数组的区别
1. 给一个函数传了切片，在函数内把一个切片扩容了，函数外的切片会有什么变化？
1. map是协程安全的吗？怎么改造
1. struct的方法的值接受者和指针接受者，使用上有什么区别
1. 如何获取一个变量的类型
<!-- 未处理的块类型: 34, block_id: RARCdcMM8oyi8zx3azgcH03CnCf -->

反射

1. type a Student 和 type a = Student的区别，能继承Student的方法吗
<!-- 未处理的块类型: 34, block_id: JBWHdMSDEoQhZSxOFFschDeMnqf -->

我说记不清了；但是争取了一下，说一般用type XXXType int来实现某个业务类型值的枚举

1. 有缓冲和无缓冲channel
1. 如何判断有缓冲的channel满了？满了之后如何让后续写入不会阻塞？
<!-- 未处理的块类型: 34, block_id: NgeCdTSSGoA4ZlxUjJIcealFnTh -->

这块不知道，我说让消费者那边去消费😂；后来查了，可以用cap和len判断，或者用select也可以

1. mysql怎么解决幻读的
<!-- 未处理的块类型: 34, block_id: ChBvde2I9of73kx5ZRdcwRJpnD8 -->

答了间隙锁，没答好

1. update会用独占锁吗
<!-- 未处理的块类型: 34, block_id: D6rLdfRJ9oyr95xkA96cCwlwnch -->

没答好

1. datetime和timestamp的区别
1. 左连接、右连接、内连接的区别
<!-- 未处理的块类型: 34, block_id: XAoVd6ApqoSJeTx6dBUcZ2Zunyh -->

有点忘了，实际开发没有特别注意这些连接

1. 用redis做过哪些功能
1. 对http协议的了解、状态码
1. 对grpc的了解、哪些传输模式、怎么写proto文件的
1. 熟悉docker和k8s吗，常用哪些命令
1. 聊家常，对加班的看法，介绍他们公司的业务
**总监面**

1. 就是聊项目难点，然后聊家常
# 深圳 屈臣氏 一面

1. 聊短剧私域项目
1. Rabbitmq和Kafka的比较
1. Redis分布式锁解决重复上传素材的问题
1. 详细说说go的GMP、GC机制
1. 为什么要使用GMP这种模型来调度，为什么要这么设计，是为了解决什么问题
1. 什么场景会出现协程饥饿/一直阻塞的情况
<!-- 未处理的块类型: 34, block_id: BKjjdiDCwoZYAixlrcyceTi2nj6 -->

用了channel、mutex，进入阻塞状态

1. 为什么Go GC要用三色标记法
没答上来，有提到其它的引用计数法，面试官说没事

1. 开始写代码题，给了一个业务场景，让我实现，包括要有读配置文件的流程，一些复杂的子流程可以用空函数：
1. 假设有一个广告投放系统，通过kafka接收APP上报上来的用户事件，事件有一个type字段，用于区分当前事件是点击事件，还是曝光事件，接收到用户事件后，如果是点击事件，需要更新点击事件的计数，并发送kafka信息到CLICK\_EVENT，如果是曝光事件，需要更新曝光事件的计数，并发送kafka信息到EXPOSURE\_EVENT，请编写代码完成该功能，需要满足以下要求：
- 需要处理重复消息
- 该系统会部署到不同的国家，中国和泰国需要该功能，韩国不需要该功能，中国和泰国对应的kafka topic不一样，程序应该能在所有国家都正常启动。
