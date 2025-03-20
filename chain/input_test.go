package chain

import "testing"

type inputTestCase struct {
	id            string
	vals          map[string]any
	targetName    string
	targetVar     inputTestCaseTargetInter
	expectedVal   any
	expectedError bool
}

type inputTestCaseTargetInter interface {
	input(key string) Input
	getVal() any
}

type inputTestCaseTargetVal[T any] struct {
	val T
}

func (i *inputTestCaseTargetVal[T]) input(key string) Input {
	return I(key, &i.val)
}

func (i *inputTestCaseTargetVal[T]) getVal() any {
	return i.val
}

var inputTestCases = []inputTestCase{
	{
		id:          "string",
		vals:        map[string]any{"target": "text", "other": 5},
		targetName:  "target",
		targetVar:   &inputTestCaseTargetVal[string]{},
		expectedVal: "text",
	},
	{
		id:          "int",
		vals:        map[string]any{"target": 2, "other": 5},
		targetName:  "target",
		targetVar:   &inputTestCaseTargetVal[int]{},
		expectedVal: 2,
	},
	{
		id:            "wrong_type",
		vals:          map[string]any{"target": 2, "other": 5},
		targetName:    "target",
		targetVar:     &inputTestCaseTargetVal[string]{},
		expectedError: true,
	},
	{
		id:            "wrong_key",
		vals:          map[string]any{"not_target": "text", "other": 5},
		targetName:    "target",
		targetVar:     &inputTestCaseTargetVal[string]{},
		expectedError: true,
	},
}

func TestInput(t *testing.T) {
	for _, testCase := range inputTestCases {
		t.Run(testCase.id, func(t *testing.T) {
			i := testCase.targetVar.input(testCase.targetName)
			ctx := NewContext(testCase.vals)
			err := extractInput(i, ctx)
			if !testCase.expectedError {
				if err != nil {
					t.Fatalf("failed to extract input: %v", err)
				}
				if testCase.targetVar.getVal() != testCase.expectedVal {
					t.Fatalf("observed %v and expected %v were not equal", testCase.targetVar.getVal(), testCase.expectedVal)
				}
			} else {
				if err == nil {
					t.Fatalf("was expecting error but none was present")
				}
			}
		})
	}
}
