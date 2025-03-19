package main

import (
	"chain"
	"fmt"
	"os"
	"steps"
)

var openAIKey = os.Getenv("OPENAI_KEY")

func main() {
	client := steps.NewOpenAIClient(openAIKey, "gpt-4o-mini")

	steps := []chain.Step{
		steps.NewLLMStep(
			client,
			[]steps.LLMQuestion{
				{ID: "company_name", Question: "What is the name of the company in question?", Type: steps.LLMText, Default: "Unnamed Company"},
				{ID: "how_many_backers", Question: "How many companies/customers trust this company? Respond with an integer and no extra text.", Type: steps.LLMInt, Default: 0},
			},
		),
		steps.NewGTEStep("how_many_backers", "enough_backers", 500),
		steps.NewConditionalStep(
			"enough_backers",
			[]chain.Step{
				steps.NewLLMStep(client, []steps.LLMQuestion{
					{ID: "catch_phrase", Question: "What is the companies catch-phrase / slogan?", Type: steps.LLMText, Default: "No Slogan"},
				}),
			},
			[]chain.Step{
				steps.NewSetStep("catch_phrase", "not_relevant"),
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
