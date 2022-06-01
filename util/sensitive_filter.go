package util

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict("../document/sensitive_dict.txt")
	if err != nil {
		log.Println("InitFilter Fail,Err=" + err.Error())
	}
}
