package provider

import (
	"github.com/poapogoogle258/myjob_interview/internel/model/dao"
)

type JobProvider interface {
	GetName() string
	FetchJobs() ([]*dao.JobModel, error)
}

var registry = make(map[string]JobProvider)
var list_provider []string

func GetListProvider() []string {
	return list_provider
}

func Register(p JobProvider) {
	if _, ok := registry[p.GetName()]; !ok {
		list_provider = append(list_provider, p.GetName())
	}

	registry[p.GetName()] = p
}

func GetProvider(name string) (JobProvider, bool) {
	p, ok := registry[name]
	return p, ok
}
