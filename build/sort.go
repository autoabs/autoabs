package build

import (
	"sort"
)

var stateRank = map[string]int{
	"building":  1,
	"pending":   2,
	"failed":    3,
	"completed": 4,
	"inactive":  4,
}

type buildsSort struct {
	builds []*Build
}

func (b *buildsSort) Len() int {
	return len(b.builds)
}

func (b *buildsSort) Less(i, j int) bool {
	return stateRank[b.builds[i].State] < stateRank[b.builds[j].State]
}

func (b *buildsSort) Swap(i int, j int) {
	b.builds[i], b.builds[j] = b.builds[j], b.builds[i]
}

func Sort(builds []*Build) {
	srt := &buildsSort{
		builds: builds,
	}
	sort.Sort(srt)
}
