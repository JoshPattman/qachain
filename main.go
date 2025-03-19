package main

import (
	"chain"
	"fmt"
	"os"
)

var openAIKey = os.Getenv("OPENAI_KEY")

var _ chain.Step = &trueIfGTE{}

type trueIfGTE struct {
	srcKey    string
	tarKey    string
	threshold int
	srvVal    int
}

// Do implements chain.Step.
func (t *trueIfGTE) Do(actions *chain.Actions) ([]chain.Step, error) {
	val := false
	if t.srvVal > t.threshold {
		val = true
	}
	actions.Set(t.tarKey, val)
	return nil, nil
}

// Inputs implements chain.Step.
func (t *trueIfGTE) Inputs() []chain.Input {
	return []chain.Input{
		chain.I(t.srcKey, &t.srvVal),
	}
}

func main() {
	client := NewOpenAIClient(openAIKey, "gpt-4o-mini")

	steps := []chain.Step{
		NewLLMStep(
			client,
			[]LLMQuestion{
				{"company_name", "What is the name of the company in question?", LLMText, "Unnamed Company"},
				{"how_many_backers", "How many companies/customers trust this company? Respond with an integer and no extra text.", LLMInt, 0},
			},
		),
		&trueIfGTE{srcKey: "how_many_backers", tarKey: "enough_backers", threshold: 500},
		chain.NewConditionalStep(
			"enough_backers",
			[]chain.Step{
				NewLLMStep(client, []LLMQuestion{
					{"catch_phrase", "What is the companies catch-phrase / slogan?", LLMText, "No Slogan"},
				}),
			},
			[]chain.Step{
				chain.NewSetStep("catch_phrase", "not_relevant"),
			},
		),
	}

	ctx := chain.NewContext(map[string]any{"document": text})
	err := chain.Run(steps, ctx)
	if err != nil {
		panic(err)
	}

	for key, val := range ctx.Values() {
		fmt.Println(key, ":", val)
	}
}

var text = `Legal-Grade™ AI

Wherever Computer Meets Contract

Built on a proprietary legal Large Language Model, Luminance's end-to-end AI platform enhances every touchpoint a business has with its contracts, from generation to negotiation and post-execution analysis.

Developed by world-leading AI experts, validated by leading lawyers, and trusted by 700+ organisations, Luminance brings Legal-Grade™ AI to the whole enterprise.
`
