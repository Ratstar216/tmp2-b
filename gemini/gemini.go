package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "asia-northeast1"
	modelName = "gemini-1.5-flash-002"
	projectID = "term6-tomotaka-hoshina"
)

func Translate(content string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return ""
	}

	prompt := genai.Text(fmt.Sprintf(
		`次のツイートを「にゃ]を使って猫っぽくしてください。 %s`,
		content))

	// Gemini APIを呼び出す
	resp, err := client.GenerativeModel(modelName).GenerateContent(ctx, prompt)
	if err != nil {
		return ""
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return ""
	}

	// レスポンスを print
	fmt.Println(string(jsonData))

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return ""
	}

	parts, ok := result["Candidates"].([]interface{})[0].(map[string]interface{})["Content"].(map[string]interface{})["Parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return ""
	}

	output := parts[0].(string)

	return output
}
