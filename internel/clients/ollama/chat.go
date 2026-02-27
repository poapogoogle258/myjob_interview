package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

type jobDescriptor interface {
	GetJobDescription() string
}

func GetSkillsRequestFromContent(jd jobDescriptor) ([]string, error) {

	ctx := context.Background()
	prompt := fmt.Sprintf(`
	Summarize the required skills from this job description, listing them item by item using the format { skills : [ skill1, skill2, skill3, ...]} 
	If it has an abbreviation like "mainSkill (subSkill1, subSkill2)," separate it only subSkill into a new skill .
	job description : 
	%s
	`, jd.GetJobDescription())

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	steam := false
	req := &api.GenerateRequest{
		Model:  "scb10x/typhoon2.5-qwen3-4b:latest", // Specify the model to use
		Stream: &steam,
		Format: json.RawMessage(`"json"`),
		Prompt: prompt,
	}

	var responseBuilder strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		responseBuilder.WriteString(resp.Response)
		return nil
	}

	// Use the Chat function to get a response
	if err := client.Generate(ctx, req, respFunc); err != nil {
		return nil, err
	}

	data := &struct {
		Skill []string `json:"skills"`
	}{}

	if err := json.Unmarshal([]byte(responseBuilder.String()), data); err != nil {
		return nil, err

	}

	return data.Skill, nil

}
