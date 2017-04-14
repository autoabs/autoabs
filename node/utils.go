package node

import (
	"github.com/autoabs/autoabs/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, nodeId string) (nde *Node, err error) {
	nde = &Node{}
	coll := db.Nodes()

	err = coll.FindOneId(nodeId, nde)
	if err != nil {
		return
	}

	err = nde.LoadSettings(db)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (nodes []*Node, err error) {
	nodes = []*Node{}
	coll := db.Nodes()

	cursor := coll.Find(&bson.M{}).Iter()

	nde := &Node{}
	for cursor.Next(nde) {
		err = nde.LoadSettings(db)
		if err != nil {
			return
		}

		nodes = append(nodes, nde)
		nde = &Node{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
