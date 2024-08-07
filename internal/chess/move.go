package chess

const (
	NORTH     int = -WIDTH
	SOUTH     int = WIDTH
	NORTH2    int = 2 * -WIDTH
	SOUTH2    int = 2 * WIDTH
	EAST      int = 1
	WEST      int = -1
	NORTHEAST int = -WIDTH + 1
	SOUTHEAST int = WIDTH + 1
	NORTHWEST int = -WIDTH - 1
	SOUTHWEST int = WIDTH - 1
	KNIGHT0   int = -2*WIDTH - 1 // Knights moves are numbeded starting at (up 2, left 1) and proceeds clockwise
	KNIGHT1   int = -2*WIDTH + 1
	KNIGHT2   int = -WIDTH + 2
	KNIGHT3   int = WIDTH + 2
	KNIGHT4   int = 2*WIDTH + 1
	KNIGHT5   int = 2*WIDTH - 1
	KNIGHT6   int = WIDTH - 2
	KNIGHT7   int = -WIDTH - 2
	INVALID   int = NUM_SQUARES + 1
)

// move returns the direction and step size from start to end
// move assumes a board of WIDTH represented as a 1D array
//
// direction is represented by the difference in indices of start and dest.
// step size is the number of steps it takes in direction to reach dest.
//
// for example, with the representation below and starting at P:
// if dest = -14 => return -7, 2
// if dest = 5 => return 1, 5
// if dest = -17 => return 0, 1   // Knight move
// if dest = 21 => return 65, 0  // invalid move
//
//	8
//	7
//	6 -18 -17 -16 -15 -14 -13 -12 -11
//	5 -10 -9  -8  -7  -6  -5  -4  -3
//	4 -2  -1   P  +1  +2  +3  +4  +5
//	3 +6  +7  +8  +9  +10 +11 +12 +13
//	2 +14 +15 +16 +17 +18 +19 +20 +21
//	1
//	   a   b   c   d   e   f   g   h
func move(start *Square, dest *Square) (int, int) {
	startFile := start.file()
	destFile := dest.file()
	startRank := start.rank()
	destRank := dest.rank()

	if startFile == destFile && start.index > dest.index {
		return NORTH, startRank - destRank
	} else if startFile == destFile && start.index < dest.index {
		return SOUTH, destRank - startRank
	} else if startRank == destRank && start.index > dest.index {
		return WEST, startFile - destFile
	} else if startRank == destRank && start.index < dest.index {
		return EAST, destFile - startFile
	}

	rankDiff := startRank - destRank
	fileDiff := startFile - destFile
	if rankDiff < 0 && fileDiff < 0 && rankDiff == fileDiff {
		return SOUTHEAST, -fileDiff
	} else if rankDiff < 0 && fileDiff > 0 && -rankDiff == fileDiff {
		return SOUTHWEST, fileDiff
	} else if rankDiff > 0 && fileDiff > 0 && rankDiff == fileDiff {
		return NORTHWEST, fileDiff
	} else if rankDiff > 0 && fileDiff < 0 && -rankDiff == fileDiff {
		return NORTHEAST, -fileDiff
	}

	if rankDiff == 2 && fileDiff == 1 {
		return KNIGHT0, 1
	} else if rankDiff == 2 && fileDiff == -1 {
		return KNIGHT1, 1
	} else if rankDiff == 1 && fileDiff == -2 {
		return KNIGHT2, 1
	} else if rankDiff == -1 && fileDiff == -2 {
		return KNIGHT3, 1
	} else if rankDiff == -2 && fileDiff == -1 {
		return KNIGHT4, 1
	} else if rankDiff == -2 && fileDiff == 1 {
		return KNIGHT5, 1
	} else if rankDiff == -1 && fileDiff == 2 {
		return KNIGHT6, 1
	} else if rankDiff == 1 && fileDiff == 2 {
		return KNIGHT7, 1
	}

	return INVALID, 0
}
