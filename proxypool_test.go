package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCore(t *testing.T) {

	currentTime := time.Now()

	fmt.Println(currentTime)

	oldTime := currentTime.AddDate(0, 0, -1)

	fmt.Println(oldTime.Unix())

}
