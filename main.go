package main

import (
	"fmt"
)

const (
	PersonMustard int = iota
	PersonPlum
	PersonGreen
	PersonPeacock
	PersonScarlet
	PersonWhite

	WeaponKnife
	WeaponCandlestick
	WeaponRevolver
	WeaponRope
	WeaponPipe
	WeaponWrench

	RoomHall
	RoomLounge
	RoomDining
	RoomKitchen
	RoomBallroom
	RoomConservatory
	RoomBilliards
	RoomLibrary
	RoomStudy
)

func PersonCards() []int {
	return []int{PersonMustard, PersonPlum, PersonGreen, PersonPeacock, PersonScarlet, PersonWhite}
}

func WeaponCards() []int {
	return []int{WeaponKnife, WeaponCandlestick, WeaponRevolver, WeaponRope, WeaponPipe, WeaponWrench}
}

func RoomCards() []int {
	return []int{RoomHall, RoomLounge, RoomDining, RoomKitchen, RoomBallroom, RoomConservatory, RoomBilliards, RoomLibrary, RoomStudy}
}

func AllCards() []int {
	return append(append(PersonCards(), WeaponCards()...), RoomCards()...)
}

/*
For 14Choose5, there should be 2002 combinations:
[[0 1 2 3 4] [0 1 2 3 5] [0 1 2 3 6] [0 1 2 3 7] [0 1 2 3 8] [0 1 2 3 9] ...
[8 9 10 12 13] [8 9 11 12 13] [8 10 11 12 13] [9 10 11 12 13]]
*/

// difference returns the elements in `a` that aren't in `b`.
func difference[T comparable](a []T, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
func PrintRows[T any](list []T) {
	for _, row := range list {
		fmt.Println(row)
	}
}

func main() {
	allCards := AllCards()
	numAllCards := len(allCards)

	weaponsCards := WeaponCards()
	roomCards := RoomCards()
	personCards := PersonCards()

	// TODO: don't hardcode this
	ownCards := []int{RoomLibrary, RoomStudy, WeaponWrench, PersonPlum}

	ownCardAssignmentMatrix := convertCombinationAssignmentToMatrix([][]int{ownCards}, numAllCards, 99)
	PrintRows(ownCardAssignmentMatrix)

	remainingWeaponsCards := difference(weaponsCards, ownCards)
	remainingRoomCards := difference(roomCards, ownCards)
	remainingPersonCards := difference(personCards, ownCards)

	numRemainingWeaponsCards := len(remainingWeaponsCards)
	numRemainingRoomCards := len(remainingRoomCards)
	numRemainingPersonCards := len(remainingPersonCards)

	roomCardAssignmentMatrix := generateCombinationMatrix(numRemainingRoomCards, 1, 98)
	weaponCardAssignmentMatrix := generateCombinationMatrix(numRemainingWeaponsCards, 1, 98)
	personCardAssignmentMatrix := generateCombinationMatrix(numRemainingPersonCards, 1, 98)

	PrintRows(roomCardAssignmentMatrix)
	PrintRows(weaponCardAssignmentMatrix)
	PrintRows(personCardAssignmentMatrix)

	PrintRows(combinatorialProduct(roomCardAssignmentMatrix, weaponCardAssignmentMatrix))

	numPlayers := 3
	numCards := 21 - 3 - len(ownCards)

	cardsPerPlayer := numCards / numPlayers
	leftoverCards := numCards % numPlayers

	// Initialize with one row and a column for each card, set to 0
	runningCombinations := make([][]int, 1)
	runningCombinations[0] = make([]int, numCards)

	for player := 1; player <= numPlayers; player++ {
		numCardsToDraw := cardsPerPlayer
		if leftoverCards > 0 {
			numCardsToDraw++
			leftoverCards--
		}
		fmt.Printf("Player %v: %v Choose %v\n", player, numCards, numCardsToDraw)
		runningCombinations = combinatorialProduct(runningCombinations, generateCombinationMatrix(numCards, numCardsToDraw, player))

		numCards -= numCardsToDraw

	}

	runningCombinations = combinatorialProduct(ownCardAssignmentMatrix, runningCombinations)
	fmt.Println(len(ownCardAssignmentMatrix))

	fmt.Println(len(runningCombinations))
	fmt.Println(runningCombinations[0])
	fmt.Println(runningCombinations[len(runningCombinations)-1])
}

func adv(p *[]int, index int, maxIdx int) bool {
	if index > 0 && (*p)[index] == maxIdx {
		recRes := adv(p, index-1, maxIdx-1)
		(*p)[index] = (*p)[index-1] + 1
		return recRes
	} else if index >= 0 && (*p)[index] < maxIdx {
		(*p)[index] += 1
		return true
	}
	return false
}

func makeRange(min int, max int) []int {
	a := make([]int, max-min+1)

	for i := range a {
		a[i] = min + i
	}
	return a
}

func genCombinations(n int, k int) [][]int {

	numCombinations := nChooseK(n, k)

	// Allocate nChooseK number of rows, each row of length k, in contiguous memory
	// https://stackoverflow.com/questions/39804861/what-is-a-concise-way-to-create-a-2d-slice-in-go
	results := make([][]int, numCombinations)
	rows := make([]int, numCombinations*k)
	for i := 0; i < numCombinations; i++ {
		results[i] = rows[i*k : (i+1)*k]
	}

	// Initial row [1,2,3,4...]
	p := makeRange(0, k-1)
	copy(rows[0:], p[0:])

	i := 1
	for adv(&p, k-1, n-1) {
		// Copy values for each cell of each row; more direct way to do this?

		copy(rows[i*k:], p[:])
		i++
	}

	return results
}

func generateCombinationMatrix(n int, k int, assignedValue int) [][]int {
	combinations := genCombinations(n, k)
	assignedCombinations := convertCombinationAssignmentToMatrix(combinations, n, assignedValue)
	return assignedCombinations
}

func combinatorialProductAugmented[T ~int](a [][]T, b [][]T) [][]T {
	results

}

func combinatorialProduct(a [][]int, b [][]int) [][]int {
	// Define combinatorial product as:
	// Given two assignment arrays representing nChooseK
	// i.e. a := [{ 1, 0, 0, 1, 0 }, { 1, 0, 0, 0, 1} ...] and
	//      b := [{ 2, 2, 0 }, { 2, 0, 2 :}, ...}
	// produces a new matrix with length len(a) * len(b)
	// where for each row of A, there is generated all
	// possible combinations of B rows interleaved within the
	// unassigned cell values
	// i.e. [ { 1, 2, 2, 1, 0 }, { 1, 2, 0, 1, 2 } ..]
	// The following condition(s) must be met:
	// (n_a - k_a) = n_b
	// all rows of A have consistent width
	// all rows of B have consistent width

	numResultRows := len(a) * len(b)
	results := make([][]int, numResultRows)
	widthA := len(a[0])
	widthB := len(b[0])
	numColumns := max(widthA, widthB)
	rows := make([]int, numResultRows*numColumns)
	for i := 0; i < numResultRows; i++ {
		results[i] = rows[i*numColumns : (i+1)*numColumns]
	}

	resultIdx := 0
	for _, rowA := range a {
		for _, rowB := range b {
			cellAIdx := 0
			cellBIdx := 0
			for cellAIdx < widthA {
				// If the cell in A is already assigned, use the value from row A
				if rowA[cellAIdx] > 0 {
					results[resultIdx][cellAIdx] = rowA[cellAIdx]
				} else {
					// If the cell in A is not assigned, use the next value from row B
					results[resultIdx][cellAIdx] = rowB[cellBIdx]
					cellBIdx++
				}
				cellAIdx++
			}
			resultIdx++
		}
	}
	return results
}

func convertCombinationAssignmentToMatrix(mat [][]int, n int, assignedValue int) [][]int {
	numCombinations := len(mat)

	// Allocate nChooseK number of rows, each row of length k, in contiguous memory
	// https://stackoverflow.com/questions/39804861/what-is-a-concise-way-to-create-a-2d-slice-in-go
	results := make([][]int, numCombinations)
	rows := make([]int, numCombinations*n)
	for i := 0; i < numCombinations; i++ {
		results[i] = rows[i*n : (i+1)*n]
	}

	for i, row := range mat {
		for _, cell := range row {
			results[i][cell] = assignedValue
		}
	}
	return results
}

func nChooseK(n int, k int) int {
	if k > n {
		return 0
	}
	if k*2 > n {
		k = n - k
	}
	if k == 0 {
		return 1
	}

	result := n

	for i := 2; i < k+1; i++ {
		result *= n - i + 1
		result /= i
	}
	return result
}
