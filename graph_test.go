package main

import "testing"

func Test_NewGraph(t *testing.T) {
	graph, err := NewGraph()
	if err != nil {
		t.Fatalf("NewGraph failed. %s\n", err)
	}
	if graph == nil {
		t.Fatal("NewGraph returns nil graph.\n")
	}

}

func Test_AddGroups(t *testing.T) {
	graph, _ := NewGraph()
	err := graph.AddGroups(100)
	if err != nil {
		t.Fatalf("AddGroups failed. %s\n", err)
	}

}

func Test_groupSize(t *testing.T) {
	testCases := []struct {
		arg      int
		expected int
	}{
		{
			arg:      1,
			expected: 1,
		},
		{
			arg:      2,
			expected: 2,
		},
		{
			arg:      3,
			expected: 2,
		},
		{
			arg:      4,
			expected: 2,
		},
		{
			arg:      5,
			expected: 3,
		},
	}
	for _, tc := range testCases {
		if actual := groupSize(tc.arg); tc.expected != actual {
			t.Fatalf("expected %d actual %d\n", tc.expected, actual)
		}
	}

}
