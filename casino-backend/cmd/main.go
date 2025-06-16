package main

import "jhgambling/backend/core"

func main() {
	casino := core.NewCasino()

	casino.Init()
	casino.Start()
}
