package main

import (
	"chain"
	"strings"
	"testing"
)

var _ chain.Step = &trueIfDocContains{}

type trueIfDocContains struct {
	contains  string
	targetKey string
	doc       string
}

// Do implements chain.Step.
func (t *trueIfDocContains) Do(actions *chain.Actions) ([]chain.Step, error) {
	val := strings.Contains(t.doc, t.contains)
	actions.Set(t.targetKey, val)
	return nil, nil
}

// Inputs implements chain.Step.
func (t *trueIfDocContains) Inputs() []chain.Input {
	return []chain.Input{
		chain.I("document", &t.doc),
	}
}

func TestStep(t *testing.T) {
	steps := []chain.Step{
		chain.NewSetStep("constant_true", true),
		&trueIfDocContains{contains: "doc", targetKey: "doc_contains_substr"},
		chain.NewConditionalStep(
			"doc_contains_substr",
			[]chain.Step{
				chain.NewSetStep("has_doc", true),
			},
			[]chain.Step{
				chain.NewSetStep("has_doc", false),
			},
		),
	}

	ctx1 := chain.NewContext(map[string]any{"document": "This is a document"})
	if err := chain.Run(steps, ctx1); err != nil {
		t.Fatal(err)
	}
	hasDoc1, err := chain.Get[bool](ctx1, "has_doc")
	if err != nil {
		t.Fatal(err)
	} else if !hasDoc1 {
		t.Fatal("incorrect prediction for doc 1")
	}

	ctx2 := chain.NewContext(map[string]any{"document": "This is not"})
	if err := chain.Run(steps, ctx2); err != nil {
		t.Fatal(err)
	}
	hasDoc2, err := chain.Get[bool](ctx2, "has_doc")
	if err != nil {
		t.Fatal(err)
	} else if hasDoc2 {
		t.Fatal("incorrect prediction for doc 1")
	}
}
