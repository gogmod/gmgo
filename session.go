package gmgo

import (
    "github.com/sirupsen/logrus"
    mgo "gopkg.in/mgo.v2"
)

type GmgoConfig struct {
    Address     []string
    Database    string
    Username    string
    Password    string
    PoolLimit   int
}

func NewGmgoSession(cfg *GmgoConfig) *mgo.Session {
    dialInfo := &mgo.DialInfo{
        Addrs:    cfg.Address,
        Database: cfg.Database,
        Username: cfg.Username,
        Password: cfg.Password,
        Timeout:  time.Second * 30,
    }

    logrus.Infof("gmgo dial info: %+v", dialInfo)

    session, err := mgo.DialWithInfo(dialInfo)
    if err != nil {
        logrus.Errorf("connection mongodb error: %v\n",err.Error())
        panic(err)
    }
    session.SetMode(mgo.Monotonic, true)
    session.SetPoolLimit(cfg.PoolLimit)
    logrus.Info("mongodb connected!")
    return session
}


