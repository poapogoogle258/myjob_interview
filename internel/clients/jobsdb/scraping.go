package jobsdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchJob(page int) (*JobResponse, error) {

	//https://th.jobsdb.com/th/golang-jobs/in-กรุงเทพมหานคร?page=1
	url := fmt.Sprintf(`https://th.jobsdb.com/th/golang-jobs/in-กรุงเทพมหานคร?page=%d`, page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูล: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	startKey := "window.SEEK_REDUX_DATA = "
	_, after, ok := strings.Cut(html, startKey)
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูล SEEK_REDUX_DATA:")
	}
	before, _, ok := strings.Cut(after, "};")
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูลจบไม่เจอ")
	}

	jsonString := []byte(before + "}")

	var jobData JobResponse
	if err := json.Unmarshal(jsonString, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil

}
