package util

import (
	"github.com/importcjj/sensitive"
)

var Filter *sensitive.Filter

func InitFilter() error {
	Filter = sensitive.New()
	err := Filter.LoadWordDict("../document/sensitive_dict.txt")
	return err
}
