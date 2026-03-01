package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

func GetSkillsRequestFromContent(jd string) (string, error) {

	ctx := context.Background()
	prompt := fmt.Sprintf(`
	Summarize the required skills from this job description, listing them item by item order by priority using the json format 
	{ 
		languages  : [ language1, language2, ...],
		frameworks : [ framework1, framework2, ... ],
		tools      : [ tool1, tool2, ... ],
		databases  : [ database1, database2, ... ],
		hardSkills : [ skill1, skill2, ... ],
		softSkills : [ skill1, skill2, ... ]
	} 
	If it has an abbreviation like "mainSkill (subSkill1, subSkill2)," separate it only subSkill into a new skill .
	job description : 
	%s
	`, jd)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", err
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
		return "", err
	}

	return responseBuilder.String(), nil

}
