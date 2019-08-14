package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSortedRegions(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "successfully sorted region list",
			want: []string{"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-south1", "asia-southeast1", "australia-southeast1", "europe-north1", "europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "global", "northamerica-northeast1", "southamerica-east1", "us-central1", "us-east1", "us-east4", "us-west1", "us-west2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSortedRegions()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("getSortedRegions returned unexpected output (-got +want):\n%s", diff)
			}
		})
	}
}
