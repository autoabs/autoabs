package node

type BuilderStats struct {
	Active int `bson:"active" json:"active"`
}

type BuilderSettings struct {
	NodeId      string `bson:"_id" json:"-"`
	Concurrency int    `bson:"concurrency" json:"concurrency"`
}
