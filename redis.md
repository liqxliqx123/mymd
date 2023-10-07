基本数据类型，常用键值类型

* 键值基本模型
  * Key-value
    * Key: 字符串
    * value：类型多样
      * string
        * 二进制安全数组
        * 场景：数值、文本数据、图片
      * list
        * 双向列表， 元素可重复
        * 场景： 排行榜，关注列表
      * hash value
        * 一个key对应多个值
        * 使用场景： 结构化数据
      * set
        * 无序集合， 去重
        * 场景： 共同关注，好友
      * Sorted set  （zset） 
        * 有序集合，元素去重， 依靠score排序
        * 场景： 有序排行榜，排序
      * bitmap
        * 安全的字节数组， 每个bit位表示位图的一位
        * 场景： 用户签到，主体统计



底层数据结构

* 全局哈希表
  * 保存所有的key-value
    * hash表项目 key,value,next
    * 键值对中键的hash值对哈希表大小取模
    * 哈希值取模后相同的hash值（hash冲突）使用链表连接
      * 冲突链不易过长
      * hash函数的选择
  * 使用两个hash表（ht0, ht1），实现渐进式rehash， 避免阻塞
    * 静态hash 类似数组， 初始为4个bucket
    * renhash hash表项的链表变长， 需要扩容bucket, 一般是扩容到2(bucket+1)
      * 涉及内存拷贝
      * 单线程阻塞
      * 触发条件
        * 负载因子
          * lf > 1
        * 约束条件
          * 没有rdb 切没有 aof
      * 操作时机
        * 伴随正常读写操作执行
        * 周期性操作





