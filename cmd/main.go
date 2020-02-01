package main

import "github.com/vsbpro/mim"

func main() {
	m := &mim.MIM{}
	m.Start(":8000", "www.google.com:80")
}
