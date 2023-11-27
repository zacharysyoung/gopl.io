// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const (
	_width, _height = 600, 320 // canvas size in pixels
	_cells          = 100      // number of grid cells)
	_xyrange        = 30.0     // axis ranges (-xyrange..+xyrange)

	_angle = math.Pi / 6 // angle of x, y axes (=30°)
)

var (
	width, height = _width, _height
	cells         = _cells
	// Derived
	xyscale = float64(width) / 2 / _xyrange // pixels per x or y unit
	zscale  = float64(height) * 0.4         // pixels per z unit
)

var sin30, cos30 = math.Sin(_angle), math.Cos(_angle) // sin(30°), cos(30°)

type polygon struct {
	ax, ay, bx, by, cx, cy, dx, dy float64
	z                              float64
}

var (
	fillPositive, fillNegative []string
)

func init() {
	fillPositive = []string{"#fff"}
	for i := 1; i < 16; i++ {
		fill := fmt.Sprintf("#%02x%02x%02x", 225+(2*i), 195-(i*13), 195-(i*13))
		fillPositive = append(fillPositive, fill)
	}

	fillNegative = []string{"000"}
	for i := 1; i < 16; i++ {
		fill := fmt.Sprintf("#%02x%02x%02x", 195-(i*13), 195-(i*13), 225+(2*i))
		fillNegative = append(fillNegative, fill)
	}
}

func getIntParam(q url.Values, key string, defVal int) (int, error) {
	s := q.Get(key)
	if s == "" {
		return defVal, nil
	}
	x, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return x, nil
}

func main() {
	http.HandleFunc("/surface", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var err error
		width, err = getIntParam(q, "width", _width)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "got width=%s; expected a positive int", q.Get("width"))
			return
		}
		height, err = getIntParam(q, "height", _height)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "got height=%s; expected a positive int", q.Get("height"))
			return
		}
		cells, err = getIntParam(q, "cells", _cells)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "got cells=%s; expected a positive int", q.Get("cells"))
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		surface(w)
	})

	log.Println("starting server at http://localhost:8999")
	log.Fatalln(http.ListenAndServe(":8999", nil))
}

func surface(w io.Writer) {
	polygons, minZ, maxZ := getPolygons()
	absMaxZ := maxZ
	if x := minZ * -1; x > absMaxZ {
		absMaxZ = x
	}
	writeSVG(w, polygons, absMaxZ)
}

func getPolygons() (polygons []polygon, minZ, maxZ float64) {
	updateMinMax := func(z float64) {
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
	}
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, err := corner(i+1, j)
			if err != nil {
				continue
			}
			updateMinMax(az)

			bx, by, bz, err := corner(i, j)
			if err != nil {
				continue
			}
			updateMinMax(bz)

			cx, cy, cz, err := corner(i, j+1)
			if err != nil {
				continue
			}
			updateMinMax(cz)

			dx, dy, dz, err := corner(i+1, j+1)
			if err != nil {
				continue
			}
			updateMinMax(dz)

			z := (az + bz + cz + dz) / 4

			polygons = append(polygons, polygon{ax, ay, bx, by, cx, cy, dx, dy, z})
		}
	}

	return
}

var ErrInfZ = errors.New("non-finite z, bad polygon") // Ex 3.1

func corner(i, j int) (float64, float64, float64, error) {
	// Find point (x,y) at corner of cell (i,j).
	x := _xyrange * (float64(i)/float64(cells) - 0.5)
	y := _xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z.
	z := f(x, y)
	if math.IsInf(z, 0) { // Ex 3.1
		return 0, 0, 0, ErrInfZ
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, nil
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

func writeSVG(w io.Writer, polygons []polygon, absMaxZ float64) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.1' "+
		"width='%d' height='%d'>\n", width, height)

	var fill string
	for _, p := range polygons {
		if p.z > 0 {
			fill = fillPositive[int(p.z/(absMaxZ/15))]
		} else {
			fill = fillNegative[int((p.z*-1)/(absMaxZ/15))]
		}
		fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
			p.ax, p.ay, p.bx, p.by, p.cx, p.cy, p.dx, p.dy, fill)
	}

	fmt.Fprintln(w, "</svg>")
}

//!-
