package node

type BuilderSettings struct {
	NodeId      string `bson:"_id" json:"-"`
	Concurrency int    `bson:"concurrency" json:"concurrency"`
}
