package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RequestData struct {
	Contents         Content          `json:"contents"`
	SafetySettings   SafetySettings   `json:"safety_settings"`
	GenerationConfig GenerationConfig `json:"generation_config"`
}

type GenerationConfig struct {
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"topP"`
	TopK        float32 `json:"topK"`
}

type SafetySettings struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type Content struct {
	Role  string `json:"role"`
	Parts Parts  `json:"parts"`
}
type Parts struct {
	Text string `json:"text"`
}

func main() {

	raw := RequestData{
		Contents: Content{
			Role: "user",
			Parts: Parts{
				Text: "Give me a recipe for banana bread.",
			},
		},
		SafetySettings: SafetySettings{
			Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
			Threshold: "BLOCK_LOW_AND_ABOVE",
		},
		GenerationConfig: GenerationConfig{
			Temperature: 0.2,
			TopP:        0.8,
			TopK:        40,
		},
	}

	reqData, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://us-central1-aiplatform.googleapis.com/v1/projects/prompt-lab-383408/locations/us-central1/publishers/google/models/gemini-pro:streamGenerateContent?alt=sse",
		bytes.NewReader(reqData),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer ")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		if len(scanner.Text()) != 0 {
			foo := new(Response)
			data := strings.Split(scanner.Text(), " ")
			datas := data[0][0:4]
			data[0] = "\"" + datas + "\":"
			dataFlow := strings.Join(data, " ")
			// log.Println(dataFlow)
			if err := json.Unmarshal([]byte("{"+dataFlow+"}"), foo); err != nil {
				log.Panic(err)
			}
			println(foo.Data.Candidates[0].Content.Parts[0].Text)
		}
	}
}

type Response struct {
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	Candidates []CandidatesResponse `json:"candidates"`
}

type CandidatesResponse struct {
	Content ContentResponse `json:"content,omitempty"`
}
type ContentResponse struct {
	Parts []TextResponse `json:"parts,omitempty"`
}
type TextResponse struct {
	Text string `json:"text,omitempty"`
}

// data -> candidates -> content -> parts -> text
