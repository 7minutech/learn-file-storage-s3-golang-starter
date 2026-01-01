package main

import (
	"testing"
)

func TestCalculateAspectRatio(t *testing.T) {
	cases := []struct {
		height, width int
		expected      string
	}{
		{1280, 720, "16:9"},
		{1920, 1080, "16:9"},
		{3840, 2160, "16:9"},
		{720, 1280, "9:16"},
		{1080, 1920, "9:16"},
		{640, 480, "other"},
	}

	for _, c := range cases {
		actual := calculateAspectRatio(c.height, c.width)

		if c.expected != actual {
			t.Errorf("calculateAspectRatio(%d, %d) = %s; want %s", c.height, c.width, actual, c.expected)
		}
	}
}
