package main

import (
	"chain"
	"fmt"
	"strconv"
	"strings"
)

type LLMType uint8

const (
	LLMText LLMType = iota
	LLMInt
	LLMFloat
)

type LLMQuestion struct {
	ID       string
	Question string
	Type     LLMType
}

var _ chain.Step = &llmStep{}

type llmStep struct {
	client    Client
	questions []LLMQuestion
}

func NewLLMStep(client Client, qs []LLMQuestion) chain.Step {
	return &llmStep{
		client:    client,
		questions: qs,
	}
}

func parse(text string, t LLMType) (any, error) {
	switch t {
	case LLMText:
		return text, nil
	case LLMFloat:
		return strconv.ParseFloat(text, 64)
	case LLMInt:
		return strconv.Atoi(text)
	default:
		return nil, fmt.Errorf("unrecognised type")
	}
}

// Do implements chain.Step.
func (l *llmStep) Do(ctx *chain.Context) ([]chain.Step, error) {
	qLines := make([]string, 0)
	for qi, q := range l.questions {
		line := fmt.Sprintf("%d: %s", qi, q.Question)
		qLines = append(qLines, line)
	}

	systemPrompt := "You are a helpful and accurate QA bot. You will answer the users questions perfectly. Respond with each answer on a new line, in format '<question_number>: <answer>', for example '5: true.'"
	questionsText := strings.Join(qLines, "\n")
	userPrompt := fmt.Sprintf("========== DOCUMENT ==========\n%s========== QUESTIONS ==========\n%s", ctx.Document(), questionsText)

	resp, _, err := l.client.GetLLMResponse(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}
	respLines := strings.Split(resp, "\n")

	for _, line := range respLines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		id, err := strconv.Atoi(strings.Trim(parts[0], " \r\n\t"))
		if err != nil {
			continue
		}
		if id < 0 || id >= len(l.questions) {
			continue
		}
		answerText := strings.Trim(parts[1], " \r\n\t")
		answer, err := parse(answerText, l.questions[id].Type)
		if err != nil {
			continue
		}
		if err := chain.Set(ctx, l.questions[id].ID, answer); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
