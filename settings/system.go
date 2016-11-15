package settings

var System = &system{
	Id: "system",
}

func init() {
	register("system", System)
}

type system struct {
	Id string `bson:"_id"`
}
