package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"os"
)

func main() {
	r := regexp.MustCompile(`FVterrainString=([.X]+)&FVinsMax=(\d+)&FVinsMin=(\d+)&FVboardX=(\d+)&FVboardY=(\d+)&FVlevel=(\d+).*`)
	host := os.Getenv("host")
	user := os.Getenv("user")
	password := os.Getenv("password")
	url := fmt.Sprintf("%s?name=%s&password=%s", host, user, password)
	run := true
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	page := string(bytes)
	start := time.Now()
	for run {
		matches := r.FindStringSubmatch(page)
		level := matches[6]
		grid := matches[1]
		max, _ := strconv.Atoi(matches[2])
		min, _ := strconv.Atoi(matches[3])
		width, _ := strconv.Atoi(matches[4])
		var solution= ""
		for i := min; i <= max; i++ {
			solution = solve(grid, width, i, canD, canR)
			if len(solution) > 0 {
				break
			}
		}
		resp.Body.Close()
		solved := fmt.Sprintf("%s&path=%s", url, solution)
		elapsed := time.Since(start)
		fmt.Printf("for level %s, took %s solution %s\n", level, elapsed, solution)
		resp, _ = http.Get(solved)
		bytes, _ = ioutil.ReadAll(resp.Body)
		page = string(bytes)
		ioutil.WriteFile("/tmp/robot" + level + "_next", bytes, 0644)
		if strings.Contains(page, "boom") {
			fmt.Println("failed")
			run = false
			resp.Body.Close()
		}
	}
}

type dirStack []rune

func (d dirStack) push(v rune) dirStack {
	return append(d, v)
}

func (d dirStack) pop() (dirStack, rune) {
	l := len(d)
	return d[:l-1], d[l-1]
}

func (d dirStack) isEmpty() bool {
	return len(d) == 0
}

type point struct {
	x, y int
}

type rect struct {
	r, c int
}

func (g grid) dist() int {
	return len(g) + len(g[0]) - 2
}

type grid [][]rune

func solve(raw string, gridWidth, pathLength int, c1, c2 canMove) string {
	fullGrid := stringToGrid(gridWidth, raw)
	rects := makeRectanglesForPathLength(pathLength)
	p := ""
	for _, r := range rects {
		grids := make([]grid, 0)
		rowStep := r.r - 1
		colStep := r.c - 1
		maxStep := max(rowStep, colStep)
		numSteps := numberOfSteps(maxStep, gridWidth)
		for i := 0; i < numSteps; i++ {
			e := extractGrid(r, point{i * rowStep, i * colStep}, fullGrid)
			grids = append(grids, e)
		}
		g := condenseGrids(grids...)
		p = findPath(g, c1, c2)
		if p != "" {
			break
		}
	}
	return p
}

func makeRectanglesForPathLength(l int) []rect {
	sumRsCs := l + 2
	rects := make([]rect, 0, l+1)
	for i := 1; i < sumRsCs; i++ {
		rects = append(rects, rect{i, sumRsCs - i})
	}
	return rects
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func numberOfSteps(stepLength, gridLength int) int {
	i := gridLength / stepLength
	if gridLength % stepLength > 0 {
		i++
	}
	return i
}

func findPath(g grid, c1, c2 canMove) string {
	p := point{0, 0}
	// ds are directions, which form a path, to be returned
	ds := make(dirStack, 0)
	v := make(map[point]bool)
	var lastD rune
	for len(ds) < g.dist() {
		var d rune
		var np point
		if d, np = c1(p, g, v); d != 0 {
			ds, d, p, np, v = move(ds, d, p, np, v)
			continue
		}
		if d, np = c2(p, g, v); d != 0 {
			ds, d, p, np, v = move(ds, d, p, np, v)
			continue
		}

		// blocked, backtrack
		if !ds.isEmpty() {
			ds, lastD = ds.pop()
			if lastD == 'D' {
				p = point{p.x, p.y - 1}
				continue

			}
			if lastD == 'R' {
				p = point{p.x - 1, p.y}
				continue
			}
		} else {
			// blocked, exit
			break
		}
	}
	if len(ds) == g.dist() {
		return string(ds)
	} else {
		return ""
	}
}
func move(ds dirStack, d rune, p point, np point, v map[point]bool) (dirStack, rune, point, point, map[point]bool) {
	ds = ds.push(d)
	p = np
	v[np] = true
	return ds, d, p, np, v
}

func canR(p point, g grid, v map[point]bool) (rune, point) {
	nextX := p.x + 1
	nextP := point{nextX, p.y}
	if nextX >= len(g[0]) {
		return 0, point{}
	}
	nextCol := nextX
	if v[nextP] {
		return 0, point{}
	}
	if g[p.y][nextCol] == '.' {
		return 'R', nextP
	}
	return 0, point{}
}

type canMove func (p point, g grid, v map[point]bool) (rune, point)

func canD(p point, g grid, v map[point]bool) (rune, point) {
	nextY := p.y + 1
	nextP := point{p.x, nextY}
	if nextY >= len(g) {
		return 0, point{}
	}
	nextRow := g[nextY]
	if v[nextP] {
		return 0, point{}
	}
	if nextRow[p.x] == '.' {
		return 'D', nextP
	}
	return 0, point{}
}

func extractGrid(rect rect, o point, master grid) grid {
	var grid = make(grid, rect.r)
	for i := 0; i < rect.r; i++ {
		grid[i] = make([]rune, rect.c)
	}
	for i := 0; i < rect.r; i++ {
		for j := 0; j < rect.c; j++ {
			masterX := o.x + i
			if masterX >= len(master) {
				break
			}
			masterY := o.y + j
			if masterY >= len(master[masterX]) {
				break
			}
			grid[i][j] = master[masterX][masterY]
		}
	}
	return grid
}

func stringToGrid(length int, raw string) grid {
	var grid = make(grid, length)
	for i := range grid {
		grid[i] = make([]rune, length)
	}

	for i, v := range raw {
		r := i / length
		c := i % length
		grid[r][c] = v
	}
	return grid
}

func condense(x, y rune) rune {
	if y == '.' || y == 0 || x == y {
		return x
	}
	return y
}

func condenseGrids(grids ...grid) grid {
	for i := 1; i < len(grids); i++ {
		for j := 0; j < len(grids[0]); j++ {
			for k := 0; k < len(grids[0][j]); k++ {
				current := grids[0][j][k]
				next := grids[i][j][k]
				grids[0][j][k] = condense(current, next)
			}
		}
	}
	return grids[0]
}
