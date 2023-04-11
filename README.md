
# not pollable issue

reproduce steps:

1. open and close /dev/net/tun frequently
2. normally write and read data from fifo or socket
3. (about a few seconds later) fifo read reports error not pollable

```bash
[root@VM-1-4-centos pollable-test]# go run main.go
20:15:08 start test
20:15:08 fifo read: read /proc/self/fd/4: not pollable
```
