package gmgo

import mgo "gopkg.in/mgo.v2"

type CollectionManager interface{}

type Collection struct {
	collection *mgo.Collection
}
