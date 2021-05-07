# ConcurrentQueue
800个协程并发执行5万次插入+5万次删除。
mutex锁版本的队列耗时最长20s-22s，未优化的LockFreeQueue在15-16s之间，优化版在12-14s之间。CAS锁对耗时的提升大概在1/4的水平。
