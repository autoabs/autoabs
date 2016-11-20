package settings

var System = &system{
	Id: "system",
}

func init() {
	register("system", System)
}

type system struct {
	Id          string   `bson:"_id"`
	TargetRepos []string `bson:"target_repos" default:"community,core,extra,multilib"`
	TargetArchs []string `bson:"target_archs" default:"any,x86_64"`
}
