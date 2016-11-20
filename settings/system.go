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

func (s *system) HasArch(arch string) bool {
	for _, arc := range s.TargetArchs {
		if arc == arch {
			return true
		}
	}

	return false
}

func (s *system) HasRepo(repo string) bool {
	for _, rep := range s.TargetRepos {
		if rep == repo {
			return true
		}
	}

	return false
}
