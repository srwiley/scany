package scanx

import (
	"image"
	"math"

	"golang.org/x/image/math/fixed"
)

type (

	// Used as entries in the scanLinks linked list
	// Cell struct {
	// 	cover, area, x, yn int
	// }
	// // Fixed point line
	// Line struct {
	// 	a, b fixed.Point26_6
	// }
	// CellWorker listends to cellChan
	// to store cell cover and area
	// values into the scanLines linked
	// list of cells
	// CellWorker struct {
	// 	cellChan  chan CellM
	// 	sweepChan chan bool
	// 	scanLinks []Cell
	// 	lineCount int
	// }
	// ScannerT is a multi-threaded version of the cl-aa
	// antialiasing algorithm. ScannerT implements the rasterx.Scanner
	// interface, so it can be used with the rasterx and oksvg packages.
	ScannerS struct {
		extent     fixed.Rectangle26_6
		scanLinks  []Cell
		collector  Collector
		firstPoint fixed.Point26_6
		lastPoint  fixed.Point26_6
		height     int
		width      int
		inPath     bool
	}
)

// LineToSegments takes lines from the s.lineChan
// and breaks them into cover and area/cover
// cell values that are sent to the cellWorkers for
// sorting and storage
func (s *ScannerS) LineToSegments(line Line) {
	s.extent.Add(line.a)
	dy := int(line.b.Y - line.a.Y)
	if dy != 0 { // A horizontal line is ignored by the CL-AA algorithm
		dx := int(line.b.X - line.a.X)
		switch {
		case dx == 0:
			s.SendVerticalLine(line)
		case dx > 0 && dy > 0:
			s.SendEastLine(line, dx, dy, 1)
		case dx < 0 && dy < 0:
			line.a, line.b = line.b, line.a
			s.SendEastLine(line, dx, dy, -1)
		case dx < 0 && dy > 0:
			s.SendWestLine(line, dx, dy, 1)
		default: // dx > 0 && dy < 0
			line.a, line.b = line.b, line.a
			s.SendWestLine(line, dx, dy, -1)
		}
	}
}

// SendWestLine takes a line that runs from low y to high y and
// increases negatively in the x axis. Cover and area/cover
// cell values are sent to the cellWorkers for sorting and storage.
func (s *ScannerS) SendWestLine(line Line, dx, dy, flip int) {
	ax := int(line.a.X)
	ay := int(line.a.Y)
	bx := int(line.b.X)
	by := int(line.b.Y)

	cx := ax >> 6
	cy := ay >> 6

	fy1 := ay & (64 - 1)
	fx1 := ax & (64 - 1)

	yn := 64 + cy<<6
	xn := cx << 6

	for yn <= by && xn >= bx {
		xw := ax + (yn-ay)*dx/dy
		yw := ay + (xn-ax)*dy/dx
		switch {
		case xw == xn && yn == yw: // corner intersection
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx1))
			yn += 64
			xn -= 64
			cx--
			cy++
			fy1 = 0
			fx1 = 64
		case yw > yn || xw > xn: // line intersects horizontal cell wall
			fx2 := xw & (64 - 1)
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		default: // line intersects vertical cell wall
			fy2 := yw & (64 - 1)
			cover := (fy2 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx1))
			xn -= 64
			cx--
			fx1 = 64
			fy1 = fy2
		}
	}
	if yn <= by { // Only horizontal cell wall intersections remain
		for yn <= by {
			xw := ax + (yn-ay)*dx/dy
			fx2 := xw & (64 - 1)
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		}
	} else { // Only vertical cell wall intersections remain
		for xn >= bx {
			yw := ay + (xn-ax)*dy/dx
			fy2 := yw & (64 - 1)
			cover := (fy2 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx1))
			xn -= 64
			cx--
			fx1 = 64
			fy1 = fy2
		}
	}
	fx2 := bx & (64 - 1)
	fy2 := by & (64 - 1)
	cover := (fy2 - fy1) * flip
	s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
}

// SendEastLine takes a line that runs from low y to high y and
// increases postively in the x axis. Cover and area/cover
// cell values that are sent to the cellWorkers for
// sorting and storage.
func (s *ScannerS) SendEastLine(line Line, dx, dy, flip int) {
	ax := int(line.a.X)
	ay := int(line.a.Y)
	bx := int(line.b.X)
	by := int(line.b.Y)

	cx := ax >> 6
	cy := ay >> 6

	fy1 := ay & (64 - 1)
	fx1 := ax & (64 - 1)

	yn := 64 + cy<<6
	xn := 64 + cx<<6

	for yn <= by && xn <= bx {
		xw := ax + (yn-ay)*dx/dy
		yw := ay + (xn-ax)*dy/dx
		switch {
		case xw == xn && yn == yw: // corner intersection
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(64+fx1))
			yn += 64
			xn += 64
			cx++
			cy++
			fy1 = 0
			fx1 = 0
		case yw > yn || xw < xn: // line intersects horizontal cell wall
			fx2 := xw & (64 - 1)
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		default: // line intersects vertical cell wall
			fy2 := yw & (64 - 1)
			cover := (fy2 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(64+fx1))
			xn += 64
			cx++
			fx1 = 0
			fy1 = fy2
		}
	}
	if yn <= by { // Only horizontal cell wall intersections remain
		for yn <= by {
			xw := ax + (yn-ay)*dx/dy
			fx2 := xw & (64 - 1)
			cover := (64 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		}
	} else { // Only vertical cell wall intersections remain
		for xn <= bx {
			yw := ay + (xn-ax)*dy/dx
			fy2 := yw & (64 - 1)
			cover := (fy2 - fy1) * flip
			s.SaveCell(cx, cy, cover, cover*(64+fx1))
			xn += 64
			cx++
			fx1 = 0
			fy1 = fy2
		}
	}
	fx2 := bx & (64 - 1)
	fy2 := by & (64 - 1)
	cover := (fy2 - fy1) * flip
	s.SaveCell(cx, cy, cover, cover*(fx2+fx1))
}

// SendVerticalLine takes a vertical line that runs in either direction
// and sends Cover and area/cover cell values to the cellWorkers for
// sorting and storage.
func (s *ScannerS) SendVerticalLine(line Line) {
	y1 := int(line.a.Y) >> 6
	y2 := int(line.b.Y) >> 6
	x1 := int(line.a.X) >> 6
	x1f2 := (int(line.a.X) - (x1 << 6)) * 2
	if y1 == y2 {
		cover := int(line.b.Y - line.a.Y)
		s.SaveCell(x1, y1, cover, cover*x1f2)
		return
	}
	y1f := int(line.a.Y) - (y1 << 6)
	y2f := int(line.b.Y) - (y2 << 6)
	flip := 1
	if y2 < y1 {
		y1, y2 = y2, y1
		y1f, y2f = y2f, y1f
		flip = -1
	}
	cover := (64 - y1f) * flip
	s.SaveCell(x1, y1, cover, cover*x1f2)
	for y := y1 + 1; y < y2; y++ {
		cover := flip << 6
		//s.cellWaiter.Add(1)
		//s.SaveCell(CellM{x: x1, y: y, cover: 64, AreaOvCov: x1f * 2, flip: flip})
		s.SaveCell(x1, y, cover, cover*x1f2)
	}
	cover = y2f * flip
	s.SaveCell(x1, y2, cover, cover*x1f2)
}

// CellSaver threadIndex determines the cell worker
// this thread acts on. One CellSaver per
// threadIndex is instantiated, so updates to scanLinks
// slice will not conflict
func (s *ScannerS) SaveCell(x, y, cover, area int) {

	// No cover or off the top or bottom so can be ignored
	if cover == 0 || y < 0 || y >= s.height {
		return
	}
	// Pin any segments going out of the side bounds to the edges.
	// Scan line cover sums should still be zero.
	if x < 0 {
		x = 0
		area = cover << 6 // area further in gets full area for the cover
	} else if x >= s.width {
		x = s.width - 1
		area = cover << 6 // area further in gets full area for the cover
	}

	//store := s.scanLinks
	ic := y // s.threads // Find the offset of the link list header sentinel
	var icPrev int
	cc := s.scanLinks[ic]
	// The algorithm expects v.x >= 0
	// as enforced above so, the sentinel x value of -1 is always less
	// than x. The icPrev value is only set if the loop fires
	// at least once, but that is ensured where it is used
	// in the default switch case bellow.
	for cc.x < x && cc.yn != -1 {
		icPrev = ic
		ic = cc.yn
		cc = s.scanLinks[ic]
	}
	switch {
	case cc.x == x:
		// Cell exists, just add area and cover
		s.scanLinks[ic].area += area
		s.scanLinks[ic].cover += cover
	case cc.yn == -1 && cc.x < x: // Add new cell to end of list, yn = -1 indicates the cell is terminal.
		s.scanLinks[ic].yn = len(s.scanLinks)
		s.scanLinks = append(s.scanLinks, Cell{x: x, yn: -1, area: area, cover: cover})
	default: //cc.x > v.x thus insert new cell into list between cc and previous cell
		s.scanLinks[icPrev].yn = len(s.scanLinks)
		s.scanLinks = append(s.scanLinks, Cell{x: x, yn: ic, area: area, cover: cover})
	}
}

// SetBounds set the boundaries in which the scanner
// is allowed to draw. Negative values are excluded
func (s *ScannerS) SetBounds(height, width int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	s.width = width
	s.height = height

	s.scanLinks = s.scanLinks[:0]

	for j := 0; j < s.height; j++ {
		s.scanLinks = append(s.scanLinks, Cell{x: -1, yn: -1})
	}
}

// NewScanS returns a single-threaded implementation of the cl-aa antialiasing algorithm
// ScannerS implements the rasterx.Scannner interface for use with the rasterx and oksvg packages.
// An object implementing the Collector interface must be provided, which will convert x,y, and alpha
// values to the target format such as an image.RGBA. A collector for image.RGBA, RGBACollector, is
// defined in this package.
func NewScanS(width, height int, collector Collector) (s *ScannerS) {
	s = &ScannerS{height: height, width: width, collector: collector}
	for j := 0; j < height; j++ {
		// Cell.x = -1 means it is a sentinel and cell.yn = -1 means the list is empty.
		s.scanLinks = append(s.scanLinks, Cell{x: -1, yn: -1})
	}
	s.extent.Max = fixed.Point26_6{X: fixed.Int26_6(-math.MaxInt32), Y: fixed.Int26_6(-math.MaxInt32)}
	s.extent.Min = fixed.Point26_6{X: fixed.Int26_6(math.MaxInt32), Y: fixed.Int26_6(math.MaxInt32)}
	return
}

//// Functions below implement the scanner interface defined in github.com/srwiley/rasterx/fill.go
// Scanner interface {
//     Start(a fixed.Point26_6)
//     Line(b fixed.Point26_6)
//     Draw()
//     GetPathExtent() fixed.Rectangle26_6
//     SetBounds(w, h int)
//     SetColor(color interface{})
//     SetWinding(useNonZeroWinding bool)
//     Clear()

//     // SetClip sets an optional clipping rectangle to restrict rendering
//     // only to that region -- if size is 0 then ignored (set to image.ZR
//     // to clear)
//     SetClip(rect image.Rectangle)
// }

// Start initiates a new path. If a path is already in
// progress it will automatically close.
func (s *ScannerS) Start(a fixed.Point26_6) {
	if s.inPath {
		if s.firstPoint != s.lastPoint {
			s.Line(s.firstPoint) // close the last path
		}
	}
	s.firstPoint = a
	s.lastPoint = a
	s.inPath = true
}

// Line adds a straight line segment to the path
func (s *ScannerS) Line(b fixed.Point26_6) {
	s.LineToSegments(Line{a: s.lastPoint, b: b})
	s.lastPoint = b
}

// Draw finishes the path if it is open
// and then sweeps the accumulated area and cover
// cell values to the collector
func (s *ScannerS) Draw() {
	if s.inPath {
		if s.firstPoint != s.lastPoint {
			s.Line(s.firstPoint) // Close the last path
		}
		s.inPath = false
	}
	s.sweep()
}

// Clear reinitializes the cell linked lists and
// the path extents to make it ready for new paths
func (s *ScannerS) Clear() {

	s.scanLinks = s.scanLinks[:s.height]
	for j := range s.scanLinks {
		s.scanLinks[j].yn = -1
	}
	// Set max/min sentinel values for extent rects
	s.extent.Max = fixed.Point26_6{X: fixed.Int26_6(-math.MaxInt32), Y: fixed.Int26_6(-math.MaxInt32)}
	s.extent.Min = fixed.Point26_6{X: fixed.Int26_6(math.MaxInt32), Y: fixed.Int26_6(math.MaxInt32)}

}

// GetPathExtent returns the bounaries of the current path
func (s *ScannerS) GetPathExtent() fixed.Rectangle26_6 {
	return s.extent
}

//SetWinding does nothing for now
func (s *ScannerS) SetWinding(useNonZeroWinding bool) {
}

//SetColor sends either a rasterx.ColorFunc or
// color.Color value to the collector
func (s *ScannerS) SetColor(color interface{}) {
	s.collector.SetColor(color)
}

// SetClip does nothing for now
func (s *ScannerS) SetClip(rect image.Rectangle) {

}

func (s *ScannerS) sweep() {
	for i := 0; i < s.height; i++ {
		ic := s.scanLinks[i].yn
		flip := 1
		if ic != -1 {
			cc := s.scanLinks[ic]
			cover := cc.cover
			lastX := cc.x
			val := cover<<6 - cc.area>>2
			if val < 0 { // first val in each line should be all pos or negative, but
				// to avoid repeat code, just testing per line for now
				flip = -1
			}
			s.collector.Sweeper(lastX, i, 1, int16(val*flip))
			ic = cc.yn
			for ic != -1 {
				cc := s.scanLinks[ic]
				// fill in gap in cells
				if cc.x-lastX-1 > 0 { // Fill gap
					s.collector.Sweeper(lastX+1, i, (cc.x - lastX - 1), int16((cover<<6)*flip))
				}
				cover += cc.cover
				lastX = cc.x
				s.collector.Sweeper(lastX, i, 1, int16((cover<<6-cc.area>>2)*flip))
				ic = cc.yn
			}
		}
	}
}