package gmgo

import mgo "gopkg.in/mgo.v2"

type Pool interface {

    Get() *mgo.Session

    Put(*mgo.Session)

    Used() int

    Close()
}


