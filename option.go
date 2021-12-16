package gocloc

import "regexp"

type ClocOptions struct {
	ReMatchDir *regexp.Regexp
}

func NewClocOptions() *ClocOptions {
	return &ClocOptions{}
}
