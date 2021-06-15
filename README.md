# 并发安全的集合-GO实现

## 并发安全的有序SET

实现的以下方法有：

加粗部分意味着一写多读的变量，所以：

1）读的时候需要用atomic读取； 

2）写的时候需要用锁+atomic

### Insert

1. **找到需要插入的两侧节点A和B**，不存在直接返回
2. 锁定节点A，检查A.next != B or A.marked，如果为真，则解锁A返回step1
3. 创建新节点X
4. X.next = B, **A.next** = X
5. 解锁节点A

### Delete

1. **找到需要删除的节点B和其前置节点A**，不存在直接返回
2. 锁定节点B，检查 B.marked == true，如果为真，则解锁B然后返回step1
3. 锁定节点A，检查 A.next != B OR A.marked，如果为真，则解锁A和B，然后返回step1
4. **b.marked = true, A.next = B.next**
5. 解锁节点A和B

### Contains

1. **找到节点X，存在则返回 !X.marked**，不存在直接返回false