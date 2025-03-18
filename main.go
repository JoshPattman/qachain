package main

import (
	"chain"
	"fmt"
	"os"
)

var openAIKey = os.Getenv("OPENAI_KEY")

func main() {
	client := NewOpenAIClient(openAIKey, "gpt-4o-mini")

	steps := []chain.Step{
		NewLLMStep(
			client,
			[]LLMQuestion{
				{"company_name", "What is the name of the company in question?", LLMText},
				{"how_many_backers", "How many companies/customers trust this company? Respond with an integer and no extra text.", LLMInt},
			},
		),
		chain.NewConditionalStep(
			func(ctx *chain.Context) (bool, error) {
				n, err := chain.Get[int](ctx, "how_many_backers")
				if err != nil {
					return false, err
				}
				return n > 500, nil
			},
			[]chain.Step{
				NewLLMStep(client, []LLMQuestion{
					{"catch_phrase", "What is the companies catch-phrase / slogan?", LLMText},
				}),
			},
			[]chain.Step{
				chain.NewSetStep("catch_phrase", "not_relevant"),
			},
		),
	}

	ctx := chain.NewContext(text)
	chain.Run(steps, ctx)

	for key, val := range ctx.Values() {
		fmt.Println(key, ":", val)
	}
}

var text = `Legal-Grade™ AI

Wherever Computer Meets Contract

Built on a proprietary legal Large Language Model, Luminance's end-to-end AI platform enhances every touchpoint a business has with its contracts, from generation to negotiation and post-execution analysis.

Developed by world-leading AI experts, validated by leading lawyers, and trusted by 700+ organisations, Luminance brings Legal-Grade™ AI to the whole enterprise.
`
