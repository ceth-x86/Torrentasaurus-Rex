package exchange

func (e *Exchange) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * e.PieceLength
	end = begin + e.PieceLength
	if end > e.Length {
		end = e.Length
	}
	return begin, end
}

func (e *Exchange) calculatePieceSize(index int) int {
	begin, end := e.calculateBoundsForPiece(index)
	return end - begin
}
