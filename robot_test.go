package main

import (
	"fmt"
	"testing"
)

func TestFindCandidatePaths(t *testing.T) {
	grid := "......X...X.X.......X.....X..X.....X.X....X..X...X..X....X....X.....XX...X.XXX....X..........X...X.X.X.XX.X.....X.X..XXX...X..X.............X....XX.X........X..X....X..X...X......X..XX....X.....XXX.XX.....X.X.....XXX.......X.......XX...XX...........XX......XXX....X..X..X..XX..X.X..X....X.XXX......X.X...X.....XXXX.........X....XXX..XX....X..........X....XX.........X....X..X...X..............X....X.X...XX.X.XX.....XX.X.XX.X.X....X...X..XX.X.XX...X....X...XX.........XX..X.......XXX...........X......XX.XX........X........X..X.......X..X..X........X.XX..X..X......X......XXX...X...XX.X..........X.X.....XX.......XX.........X....X.XXX....X.....X.........X......X...X............X..XXX..X..X..XX........X.......X................X.....XX.............X...X.X..X...........X.X.....XX....XX...X........X...X..X..X..XX........X...............X.XXXXX.X.X.......X......X.XX......X......X....X.XXXX......X.......X..........X.X.X..XXX....X....X.X.X.....X.X.X.......XX....X.X.X.XX...X..X....X.....X.XX...............X.......X..X.........XX....XX.......X.X........X...X....XX.XXX.X..X......X.......X..X....X..X.X...XX.......XX.XXX.X.X..X....X.X.....X....XXXX.X....X..X......X..X..........X.....X..X...XX.......XX.XX.....X....X........X...X............X...X.X.........XX....X........X......X..X..X...X.......X.XX.....XXX.X..XX....X.X.X...XXX.X.........X.........X..X......X.X.XX.X...XXX.XX.X.....X....X...X..XXX..X....X....X....X.XX...XXX.........XX.X..X.............X....X..X...XX..XX...XXXXX.......XX...XX..........X.....XXX..XX.X.............XXX.......X..X..X..X...X.X....X.X...XX.X.X..XXX..X...XX........X.....X....XX.......X.X...X..X...XX.....XXX.........XX.....X.X.....XXX..X........XX......X....X.X.........X...X...X..X..X..X.XX...XX....X......X.X....X..X.XX..XX........"
	gw := 42
	actual := ""
	for i := 2; i <= 42; i++ {
		actual = solve(grid, gw, i, canD, canR)
		if len(actual) > 0 {
			break
		}
	}
	expect := "RDDRRRDRDRRRDRRDRDDDRDDRR"
	if expect != actual {
		t.Errorf("expected: %v, actual: %v", expect, actual)
	}
}

func TestLevel5(t *testing.T) {
	g := ".X.X....XX.....X....XX..."
	actual := solve(g, 5, 2, canD, canR)
	expect := "DR"
	if expect != actual {
		t.Errorf("expected: %v, actual: %v", expect, actual)
	}
}

func TestLevel10(t *testing.T) {
	g := ".X.X.X...X.............XX..X......XX.....X...XX.XX.........X.....X.....X.X........X....X..X.XX....X."
	actual := solve(g, 10, 4, canD, canR)
	expect := ""
	if expect != actual {
		t.Errorf("expected: %v, actual: %v", expect, actual)
	}
	actual = solve(g, 10, 5, canD, canR)
	expect = "DDDRR"
	if expect != actual {
		t.Errorf("expected: %v, actual: %v", expect, actual)
	}
	actual = solve(g, 10, 6, canD, canR)
	expect = "DDRRDD"
	if expect != actual {
		t.Errorf("expected: %v, actual: %v", expect, actual)
	}
}

func Test10Extraction(t *testing.T) {
	fullGrid := stringToGrid(10, ".X.X.X...X.............XX..X......XX.....X...XX.XX.........X.....X.....X.X........X....X..X.XX....X.")
	e := extractGrid(rect{5, 1}, point{8, 0}, fullGrid)
	fmt.Printf("e: %v\n", e)
}

func TestGridOverlay(t *testing.T) {
	raw := "012345678"
	length := 3
	actual := stringToGrid(length, raw)
	expected := [][]rune{{'0', '1', '2'}, {'3', '4', '5'}, {'6', '7', '8'}}
	assertGridsEqual(t, expected, actual)
}

func TestCondensingGrid(t *testing.T) {
	raw := "...................................."
	length := 6
	g := stringToGrid(length, raw)
	g[2][5] = 'X'
	expected := grid{{'.', 'X', '.'}, {'.', '.', '.'}}
	grids := make([]grid, 0)
	rowStep := 1
	colStep := 2
	maxStep := max(rowStep, colStep)
	numSteps := numberOfSteps(maxStep, length)
	rect := rect{rowStep + 1, colStep + 1}
	for i := 0; i < numSteps; i++ {
		e := extractGrid(rect, point{i * rowStep, i * colStep}, g)
		grids = append(grids, e)
	}
	actual := condenseGrids(grids...)
	assertGridsEqual(t, expected, actual)
}

func TestNumberOfSteps(t *testing.T) {
	n := numberOfSteps(4, 10)
	if n != 3 {
		t.Errorf("len exp: %d len act %d", 3, n)
	}
}

func TestCreateRectangles(t *testing.T) {
	length := 3
	act := makeRectanglesForPathLength(length)
	exp := []rect{{1, 4}, {2, 3}, {3, 2}, {4, 1}}
	if len(act) != len(exp) {
		t.Errorf("len exp: %d len act %d", len(exp), len(act))
	}
	for i, r := range exp {
		if r != act[i] {
			t.Errorf("exp: %v act %v", r, act[i])
		}
	}

	length = 2
	act = makeRectanglesForPathLength(length)
	exp = []rect{{1, 3}, {2, 2}, {3, 1}}
	if len(act) != len(exp) {
		t.Errorf("len exp: %d len act %d", len(exp), len(act))
	}
	for i, r := range exp {
		if r != act[i] {
			t.Errorf("exp: %v act %v", r, act[i])
		}
	}
}

func TestCanDirections(t *testing.T) {
	var v map[point]bool
	v = map[point]bool{
		point{0, 0}: false,
	}
	g := grid{{'.', 'X', '.'}, {'.', '.', 'X'}}
	if r, _ := canD(point{0, 0}, g, v); r == 0 {
		t.Error("D should be available")
	}
	if r, _ := canD(point{0, 1}, g, v); r > 0 {
		t.Error("D should not be available as at bottom")
	}
	if r, _ := canD(point{0, 2}, g, v); r > 0 {
		t.Error("D should not be available as X at 1,2")
	}
	if r, _ := canR(point{0, 1}, g, v); r == 0 {
		t.Error("R should be available")
	}
	if r, _ := canR(point{2, 0}, g, v); r > 0 {
		t.Error("R should be not be available as at edge")
	}
	if r, _ := canR(point{0, 0}, g, v); r > 0 {
		t.Error("R should not be available as X at 1,0")
	}
	v = map[point]bool{
		point{0, 1}: true,
		point{1, 1}: true,
	}
	if r, _ := canD(point{0, 0}, g, v); r > 0 {
		t.Error("D should be not be available, visited")
	}
	if r, _ := canR(point{0, 1}, g, v); r > 0 {
		t.Error("R should not be available, visited")
	}
}

func TestFindPath(t *testing.T) {
	raw := "...."
	length := 2
	g := stringToGrid(length, raw)
	fmt.Printf("g %d\n", g.dist())

	act := findPath(g, canD, canR)
	exp := "DR"
	if act != exp {
		t.Errorf("expected %v actual %v", exp, act)
	}
	act = findPath(g, canR, canD)
	exp = "RD"
	if act != exp {
		t.Errorf("expected %v actual %v", exp, act)
	}
}

func assertGridsEqual(t *testing.T, expected, actual grid) {
	for i, r := range expected {
		for j, exp := range r {
			act := actual[i][j]
			if exp != act {
				t.Errorf("expected %#U actual %#U", exp, act)
			}
		}
	}
}
