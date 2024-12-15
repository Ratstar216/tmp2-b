package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "asia-northeast1"        // モデルのリージョン
	modelName = "gemini-1.5-flash-002"   // 使用するGeminiモデル名
	projectID = "term6-tomotaka-hoshina" // GCPプロジェクトID
)

func Translate(content string) string {
	// Geminiクライアントの作成
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return ""
	}

	// プロンプトを作成（信頼度スコアの取得）
	prompt := genai.Text(fmt.Sprintf(
		`次のツイートを「にゃ]を使って猫っぽくしてください。 %s`,
		content))

	// Gemini APIを呼び出す
	resp, err := client.GenerativeModel(modelName).GenerateContent(ctx, prompt)
	if err != nil {
		return ""
	}

	// レスポンスの確認
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return ""
	}

	// レスポンスを print
	fmt.Println(string(jsonData))

	// JSON から "Parts" の値を取得する
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return ""
	}

	// "Parts" フィールドの取り出し
	parts, ok := result["Candidates"].([]interface{})[0].(map[string]interface{})["Content"].(map[string]interface{})["Parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return ""
	}

	output := parts[0].(string)

	return output
}
