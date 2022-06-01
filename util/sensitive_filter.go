package util

import (
	"github.com/importcjj/sensitive"
)

const NetWordDicUrl = "https://raw.githubusercontent.com/importcjj/sensitive/master/dict/dict.txt"

var Filter *sensitive.Filter

func InitFilter() error {
	Filter = sensitive.New()
	err := Filter.LoadNetWordDict(NetWordDicUrl)
	return err
}
