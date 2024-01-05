package stream_inspector

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const magic = 0

func TestFlamegraphConverter_covertTrees(t *testing.T) {
	tests := []struct {
		name  string
		trees []*Tree
		want  FlameBearer
	}{
		{
			name: "expected flame graph to be built with offset for the second level",
			trees: []*Tree{
				// 1st tree
				{
					Root: &Node{
						Name: "top_level-a", Weight: 100, Children: []*Node{
							{Name: "second_level-a", Weight: 50, Children: []*Node{
								{Name: "third_level_a", Weight: 20, Children: []*Node{
									{Name: "fourth_level_a_0", Weight: 10},
									{Name: "fourth_level_a_1", Weight: 5},
								}}},
							},
						},
					},
				},
				// 2nd tree
				{
					Root: &Node{
						Name: "top_level-b", Weight: 50, Children: []*Node{
							{Name: "second_level-b", Weight: 50, Children: []*Node{
								{Name: "third_level_b", Weight: 10, Children: []*Node{
									{Name: "fourth_level_b", Weight: 5, Children: []*Node{
										{Name: "fives_level_b", Weight: 2, Children: []*Node{
											{Name: "sixth_level_b", Weight: 1},
										}},
									}},
								}}},
							},
						},
					},
				},
			},
			want: FlameBearer{
				Units:    "bytes",
				NumTicks: 150,
				MaxSelf:  150,
				Names: []string{
					"top_level-a", "top_level-b",
					"second_level-a", "second_level-b",
					"third_level_a", "third_level_b",
					"fourth_level_a_0", "fourth_level_a_1", "fourth_level_b",
					"fives_level_b",
					"sixth_level_b",
				},
				Levels: [][]float64{
					// for each block: start_offset, end_offset, unknown_yet, index from Names slice
					// 1st level
					{ /*1st block*/ 0, 100, magic, 0 /*2nd block*/, 0, 50, magic, 1},
					// 2nd level
					{ /*1st block*/ 0, 50, magic, 2 /*2nd block*/, 50, 50, magic, 3},
					// 3rd level
					{ /*1st block*/ 0, 20, magic, 4 /*2nd block*/, 80, 10, magic, 5},
					// 4th level
					{ /*1st block*/ 0, 10, magic, 6 /*2nd block*/, 0, 5, magic, 7 /*3rd block*/, 85, 5, magic, 8},
					// 5s level
					{ /*1st block*/ 100, 2, magic, 9},
					// 6th level
					{ /*1st block*/ 100, 1, magic, 10},
				},
			},
		},
		{
			name: "expected flame graph to be built",
			trees: []*Tree{
				// 1st tree
				{
					Root: &Node{
						Name: "top_level-a", Weight: 100, Children: []*Node{
							{Name: "second_level-a", Weight: 100},
						},
					},
				},
				// 2nd tree
				{
					Root: &Node{
						Name: "top_level-b", Weight: 50, Children: []*Node{
							{Name: "second_level-b", Weight: 50},
						},
					},
				},
			},
			want: FlameBearer{
				Units:    "bytes",
				NumTicks: 150,
				MaxSelf:  150,
				Names:    []string{"top_level-a", "top_level-b", "second_level-a", "second_level-b"},
				Levels: [][]float64{
					// for each block: start_offset, end_offset, unknown_yet, index from Names slice
					// 1st level
					{ /*1st block*/ 0, 100, magic, 0 /*2nd block*/, 0, 50, magic, 1},
					// 2nd level
					{ /*1st block*/ 0, 100, magic, 2 /*2nd block*/, 0, 50, magic, 3},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FlamegraphConverter{}
			result := f.CovertTrees(tt.trees)
			require.Equal(t, tt.want, result)
		})
	}
}