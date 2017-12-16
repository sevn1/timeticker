package timeticker

import (
	"log"
	"sync"
	"time"
)

type TimeTicker struct {
	sync.Mutex                 //锁
	interval   time.Duration   //频率，每隔多长时间close一个channel
	ticker     *time.Ticker    //定时器
	quit       chan struct{}   //结束标示
	maxTimeout time.Duration   //最大间隔多长时间开始执行
	cs         []chan struct{} //channel集合
	pos        int             //指针
}

//初始化方法
//创建了一个timingwheel，精度/频率是interval，最大的超时等待时间为buckets，并生成buckets个channel
func Init(interval time.Duration, buckets int) *TimeTicker {
	w := new(TimeTicker)
	w.interval = interval
	w.quit = make(chan struct{})
	w.pos = 0
	w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))
	w.cs = make([]chan struct{}, buckets)
	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	w.ticker = time.NewTicker(interval)
	go w.Trun()
	return w
}

//停止定时器
func (w *TimeTicker) Stop() {
	close(w.quit)
}

//返回对应的channel
//根据间隔（interval的倍数），返回多少时间后close的channel，用于定时执行任务
func (w *TimeTicker) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		log.Println("timeout too much, over maxtimeout")
	}
	index := int(timeout / w.interval)
	if 0 < index {
		index--
	}
	w.Lock()
	index = (w.pos + index) % len(w.cs)
	b := w.cs[index]
	w.Unlock()
	return b
}

//定时close channel
func (w *TimeTicker) Trun() {
	for {
		select {
		case <-w.ticker.C:
			w.TonTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

//close channel并补齐对应的长度
func (w *TimeTicker) TonTicker() {
	w.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.Unlock()
	close(lastC)
}
