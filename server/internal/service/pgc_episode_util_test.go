package service

import (
	"testing"

	"interastral-peace.com/alnitak/internal/domain/dto"
)

func TestCountEpisodeReqWithBoundVideo(t *testing.T) {
	t.Parallel()
	ep := []dto.EpisodeReq{
		{EpisodeNumber: 1, VID: 0},
		{EpisodeNumber: 2, VID: 10},
		{EpisodeNumber: 3, VID: 10},
	}
	if got := countEpisodeReqWithBoundVideo(ep); got != 2 {
		t.Fatalf("bound count: got %d want 2", got)
	}
}

func TestUniquePositiveVidsFromEpisodeReq(t *testing.T) {
	t.Parallel()
	ep := []dto.EpisodeReq{
		{VID: 0},
		{VID: 5},
		{VID: 5},
		{VID: 7},
	}
	got := uniquePositiveVidsFromEpisodeReq(ep)
	if len(got) != 2 || got[0] != 5 || got[1] != 7 {
		t.Fatalf("unique vids: got %#v want [5 7]", got)
	}
}
