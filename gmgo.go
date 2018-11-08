package gmgo

import (
	"time"

	"github.com/gogmod/pool"
	"github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type GmgoSrc struct {
	*mgo.Session
}

var (
	GmgoPool     pool.Pool
	session      *mgo.Session
	err          error
	MGO_CONN_STR string
	MGO_CONN_CAP int
)

type Config struct {
	GmgoConnStr      string
	GmgoConnCap      int
	GmgoConnGcSecond int64
}

//初始化mgo
func InitGmgo(config Config) {
	var connGcSecond = time.Duration(config.GmgoConnGcSecond) * 1e9
	MGO_CONN_STR = config.GmgoConnStr
	MGO_CONN_CAP = config.GmgoConnCap
	GmgoPool = pool.ClassicPool(
		config.GmgoConnCap, config.GmgoConnCap/5,
		func() (pool.Src, error) {
			if err != nil || session.Ping() != nil {
				if session != nil {
					session.Close()
				}
				Refresh()
			}
			return &GmgoSrc{session.Clone()}, err
		},
		connGcSecond)
}

//Refresh ....
func Refresh() {
	session, err = mgo.Dial(MGO_CONN_STR)
	if err != nil {
		logrus.Errorf("Gmgo: %v\n", err)
	} else if err = session.Ping(); err != nil {
		logrus.Errorf("Gmgo: %v\n", err)
	} else {
		session.SetPoolLimit(MGO_CONN_CAP)
	}
}

// 判断资源是否可用
func (self *GmgoSrc) IsUsable() bool {
	if self.Session == nil || self.Session.Ping() != nil {
		return false
	}
	return true
}

//重置方法
func (self *GmgoSrc) Reset() {}

//释放方法
func (self *GmgoSrc) Release() {
	if self.Session == nil {
		return
	}
	self.Session.Close()
}

//调用资源池中的资源
func Call(fn func(pool.Src) error) error {
	return GmgoPool.Call(fn)
}

//销毁资源池
func Release() {
	GmgoPool.Release()
}

//返回当前剩余资源数量
func Len() int {
	return GmgoPool.Len()
}

// DatabaseNames... 获取所有database
func DatabaseNames() (names []string, err error) {
	err = GmgoPool.Call(func(src pool.Src) error {
		names, err = src.(*GmgoSrc).DatabaseNames()
		return err
	})
	return
}

// CollectionNames... 获取所有集合
func CollectionNames(dbname string) (names []string, err error) {
	GmgoPool.Call(func(src pool.Src) error {
		names, err = src.(*GmgoSrc).DB(dbname).CollectionNames()
		return err
	})
	return
}
