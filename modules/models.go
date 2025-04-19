package modules

type Input struct {
	Move  string
	IsFEN bool
}

type CalculatedMove struct {
	Move              string
	ScoreInCP         string // M1, -M2, 1.5, -3.21
	ContinuationMoves []string
}
