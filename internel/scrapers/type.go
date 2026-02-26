package scrapers

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

type Contact string

var (
	Email    Contact = "email"
	Phone    Contact = "phone"
	Name     Contact = "name"
	Location Contact = "location"
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
