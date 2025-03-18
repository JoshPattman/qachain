package main

import (
	"chain"
	"strings"
	"testing"
)

func TestStep(t *testing.T) {
	steps := []chain.Step{
		chain.NewSetStep("constant_true", true),
		chain.NewConditionalStep(
			func(ctx *chain.Context) (bool, error) {
				return strings.Contains(ctx.Document(), "doc"), nil
			},
			[]chain.Step{
				chain.NewSetStep("has_doc", true),
			},
			[]chain.Step{
				chain.NewSetStep("has_doc", false),
			},
		),
	}

	ctx1 := chain.NewContext("This is a document")
	if err := chain.Run(steps, ctx1); err != nil {
		t.Fatal(err)
	}
	hasDoc1, err := chain.Get[bool](ctx1, "has_doc")
	if err != nil {
		t.Fatal(err)
	} else if !hasDoc1 {
		t.Fatal("incorrect prediction for doc 1")
	}

	ctx2 := chain.NewContext("This is not")
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
