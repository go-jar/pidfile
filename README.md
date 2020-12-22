# 思路

1. 创建 pidfile 时，如果正式的 pidfile 文件存在，且其记载的 pid 所对应的进程存在，则将新的 pid 记录在临时 pidfile 中；否则，记录在正式 pidfile 中。
2. 删除 pidfile 时，如果正式的 pidfile 中记录的 pid 是当前要删的进程的，则将其删除，并将临时 pidfile 重命名为正式 pidfile；否则，删除临时 pidfile。

# 参考
- https://github.com/goinbox/pidfile
