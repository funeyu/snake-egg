package models

import "testing"

func TestAllKeywords(t *testing.T) {
	res := DecideKeywords("彭玉平：陈寅恪《王观堂先生挽词并序》疏")
	println("res", res)
}
