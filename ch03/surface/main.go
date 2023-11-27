// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"errors"
	"fmt"
	"math"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	type point struct {
		ax, ay, bx, by, cx, cy, dx, dy float64
		z                              float64
	}
	var (
		maxZ, minZ float64
		points     []point
	)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, _, err := corner(i+1, j)
			if err != nil {
				continue
			}
			bx, by, _, err := corner(i, j)
			if err != nil {
				continue
			}
			cx, cy, _, err := corner(i, j+1)
			if err != nil {
				continue
			}
			dx, dy, z, err := corner(i+1, j+1)
			if err != nil {
				continue
			}
			if z < minZ {
				minZ = z
			}
			if z > maxZ {
				maxZ = z
			}
			points = append(points, point{ax, ay, bx, by, cx, cy, dx, dy, z})
		}
	}

	fillPos := []string{"#dddddd"}
	for i := 1; i < 16; i++ {
		fill := fmt.Sprintf("#%02x%02x%02x", 225+(2*i), 195-(i*13), 195-(i*13))
		fillPos = append(fillPos, fill)
	}
	fillNeg := []string{"#dddddd"}
	for i := 1; i < 16; i++ {
		fill := fmt.Sprintf("#%02x%02x%02x", 195-(i*13), 195-(i*13), 225+(2*i))
		fillNeg = append(fillNeg, fill)
	}

	maxDelta := maxZ
	if x := minZ * -1; x > maxDelta {
		maxDelta = x
	}

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	var (
		i    int
		fill string
	)
	for _, p := range points {
		if p.z > 0 {
			i = int(p.z / (maxDelta / 15))
			fill = fillPos[i]
		} else {
			i = int((p.z * -1) / (maxDelta / 15))
			fill = fillNeg[i]
		}
		fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
			p.ax, p.ay, p.bx, p.by, p.cx, p.cy, p.dx, p.dy, fill)
	}

	fmt.Println("</svg>")
}

var ErrInfZ = errors.New("non-finite z, bad polygon") // Ex 3.1

func corner(i, j int) (float64, float64, float64, error) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	if math.IsInf(z, 0) { // Ex 3.1
		return 0, 0, 0, ErrInfZ
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, nil
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
