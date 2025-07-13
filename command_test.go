package main

import "testing"

func TestCommandMap(t *testing.T) {

	var cfg *config
	err := commandMap(cfg)
	if err != nil {
		t.Log(err)
	}

}
