# timeticker

timeticker是解决大量使用time.NewTicker会导致负载暴增影响服务器性能的包。
这是根据一个time.NewTicker和多个channel实现的

## Sample

**sample.go**
```golang
package main

import (
	"log"
	"time"

	"github.com/sevn1/timeticker"
)

func main() {
	TimerRate := time.Millisecond * 400
	t := timeticker.Init(TimerRate, 20)
	for {
		select {
		case <-t.After(TimerRate):
			// code...
			log.Println("ok")
		}
	}
}
```
2017/12/17 02:37:30 ok
2017/12/17 02:37:31 ok
2017/12/17 02:37:31 ok
2017/12/17 02:37:31 ok
2017/12/17 02:37:32 ok
2017/12/17 02:37:32 ok
2017/12/17 02:37:33 ok
2017/12/17 02:37:33 ok
2017/12/17 02:37:33 ok
```