package relkw

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	mapKw, _ := GetRelKw("facebook ads")
	fmt.Printf("len mapKw: %v\n", len(mapKw))
}
