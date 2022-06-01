package sensitive_filter

import (
	"github.com/importcjj/sensitive"
)

var Filter *sensitive.Filter

func InitFilter() error {
	Filter = sensitive.New()
	err := Filter.LoadWordDict("sensitive_dict.txt")
	return err
}
