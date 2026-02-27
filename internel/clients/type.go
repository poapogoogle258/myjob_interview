package clients

import (
	"slices"
	"strings"
)

type JobType string

var (
	OnSite JobType = "On-Site"
	WFH    JobType = "Work-From-Home"
	Hybrid JobType = "Hybrid"
)

type Contact uint8

var (
	Email    Contact = 0
	Phone    Contact = 1
	Name     Contact = 2
	Location Contact = 3
)

type Tags map[string]struct{}

func (t Tags) Add(s string) {
	t[s] = struct{}{}
}

func (t Tags) Remove(s string) {
	delete(t, s)
}

func (t Tags) GetListString() []string {
	res := make([]string, len(t))
	i := 0
	for k := range t {
		res[i] = strings.ToLower(k)
	}
	slices.Sort(res)

	return res
}
