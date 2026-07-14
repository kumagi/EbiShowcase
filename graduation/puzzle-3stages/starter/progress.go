package starter

func NextStage(stage int, solved bool) int {
	if solved && stage < 3 {
		return stage + 1
	}
	return stage
}
