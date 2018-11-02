package gmgo

import (
    "sync"
    "time"
)

type Shrimp interface{
    New() Shrimp
    Close()
    Clean()
    IsValid() bool
}

type Pool struct{
    Cap int
    Src map[Shrimp]bool
    Shrimp
    sync.Mutex
}

func NewPool(shrimp Shrimp, size ...int) *Pool {
    if len(size) == 0 {
        size = append(size,1024)
    }
    return &Pool{
            Cap: size[0],
            Src: make(map[Shrimp]bool),
            Shrimp: shrimp,
        }
}

type Default struct{}

func (Default) New() Shrimp {
    return nil
}

func (Default) Close() {}
func (Default) Clean() {}
func (Default) IsValid() bool {
    return true
}

func (self *Pool) Free(m ...Shrimp){
    for i, count := 0, len(m); i < count; i++ {
        m[i].Clean()
        self.Src[m[i]] = false
    }
}

func (self *Pool) Remove(m ...Shrimp){
    for _, c := range m {
        c.Close()
        delete(self.Src, c)
    }
}

func (self *Pool) Reset(){
    for k, _ := range self.Src {
        k.Close()
        delete(self.Src,k)
    }
}

func (self *Pool) increment(){
    if len(self.Src) < self.Cap {
        self.Src[self.Shrimp.New()] = false
    }
}

func (self *Pool) GetOne() Shrimp {
    self.Mutex.Lock()
    defer self.Mutex.Unlock()

    for {
        for k, v := range self.Src{
            if v {
                continue
            }
            if !k.IsValid() {
                self.Remove(k)
                continue
            }
            self.Src[k] = true
            return k
        }
        if len(self.Src) <= self.Cap {
            self.increment()
        }else{
            time.Sleep(5e8)
        }
    }
    return nil
}

