package jobsdb_test

import (
	"fmt"
	"testing"

	"github.com/poapogoogle258/myjob_interview/internel/clients/jobsdb"
)

func TestFetchJob(t *testing.T) {
	data, _ := jobsdb.FetchJob(1)
	fmt.Println(data)
}
