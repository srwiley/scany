package scany

import (
	"image"
	"math"
	"sync"

	"golang.org/x/image/math/fixed"
)

type (
	// Used to message cells via the channels
	CellM struct {
		cover, AreaOvCov, flip, x, y int
	}
	// Used as entries in the scanLinks linked list
	Cell struct {
		cover, area, x, yn int
	}
	// Fixed point line
	Line struct {
		a, b fixed.Point26_6
	}
	// CellWorker listends to cellChan
	// to store cell cover and area
	// values into the scanLines linked
	// list of cells
	CellWorker struct {
		cellChan  chan CellM
		sweepChan chan bool
		scanLinks []Cell
		lineCount int
	}
	// ScannerT is a multi-threaded version of the cl-aa
	// antialiasing algorithm. ScannerT implements the rasterx.Scanner
	// interface, so it can be used with the rasterx and oksvg packages.
	ScannerT struct {
		lineChan    chan Line
		extents     []fixed.Rectangle26_6
		cellWorkers []CellWorker
		collector   Collector
		lineWaiter  sync.WaitGroup
		cellWaiter  sync.WaitGroup
		sweepWaiter sync.WaitGroup
		firstPoint  fixed.Point26_6
		lastPoint   fixed.Point26_6
		threads     int
		height      int
		width       int
		inPath      bool
	}
)

func Include(r fixed.Rectangle26_6, p fixed.Point26_6) fixed.Rectangle26_6 {
	if p.X < r.Min.X {
		r.Min.X = p.X
	}
	if p.X > r.Max.X {
		r.Max.X = p.X
	}
	if p.Y < r.Min.Y {
		r.Min.Y = p.Y
	}
	if p.Y > r.Max.Y {
		r.Max.Y = p.Y
	}
	return r
}

func Expand(r fixed.Rectangle26_6, s fixed.Rectangle26_6) fixed.Rectangle26_6 {
	if s.Min.X < r.Min.X {
		r.Min.X = s.Min.X
	}
	if s.Max.X > r.Max.X {
		r.Max.X = s.Max.X
	}
	if s.Min.Y < r.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if s.Max.Y > r.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// LineToSegments takes lines from the s.lineChan
// and breaks them into cover and area/cover
// cell values that are sent to the cellWorkers for
// sorting and storage
func (s *ScannerT) LineToSegments(i int) {
	for line := range s.lineChan {
		s.extents[i] = Include(s.extents[i], line.a)
		dx := int(line.b.X - line.a.X)
		dy := int(line.b.Y - line.a.Y)
		if dy != 0 { // A horizontal line is ignored by the CL-AA algorithm
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
		s.lineWaiter.Done()
	}
}

// SendWestLine takes a line that runs from low y to high y and
// decreases along in the x axis. Cover and area/cover
// cell values are sent to the cellWorkers for sorting and storage.
func (s *ScannerT) SendWestLine(line Line, dx, dy, flip int) {
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
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx1, flip: flip}
			yn += 64
			xn -= 64
			cx--
			cy++
			fy1 = 0
			fx1 = 64
		case yw > yn || xw > xn: // line intersects horizontal cell wall
			fx2 := xw - xn
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		default: // line intersects vertical cell wall
			fy2 := -yn + 64 + yw
			cover := (fy2 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx1, flip: flip}
			xn -= 64
			cx--
			fx1 = 64
			fy1 = fy2
		}
	}
	if yn <= by { // Only horizontal cell wall intersections remain
		for yn <= by {
			xw := ax + (yn-ay)*dx/dy
			fx2 := xw - xn
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		}
	} else { // Only vertical cell wall intersections remain
		for xn >= bx {
			yw := ay + (xn-ax)*dy/dx
			fy2 := -yn + 64 + yw
			cover := (fy2 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx1, flip: flip}
			xn -= 64
			cx--
			fx1 = 64
			fy1 = fy2
		}
	}
	fx2 := bx & (64 - 1)
	fy2 := by & (64 - 1)
	cover := (fy2 - fy1)
	s.cellWaiter.Add(1)
	s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
}

// SendEastLine takes a line that runs from low y to high y and
// increases along in the x axis. Cover and area/cover
// cell values that are sent to the cellWorkers for
// sorting and storage.
func (s *ScannerT) SendEastLine(line Line, dx, dy, flip int) {
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
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: 64 + fx1, flip: flip}
			yn += 64
			xn += 64
			cx++
			cy++
			fy1 = 0
			fx1 = 0
		case yw > yn || xw < xn: // line intersects horizontal cell wall
			fx2 := xw - xn + 64
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		default: // line intersects vertical cell wall
			fy2 := -yn + 64 + yw
			cover := (fy2 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: 64 + fx1, flip: flip}
			xn += 64
			cx++
			fx1 = 0
			fy1 = fy2
		}
	}
	if yn <= by { // Only horizontal cell wall intersections remain
		for yn <= by {
			xw := ax + (yn-ay)*dx/dy
			fx2 := xw - xn + 64
			cover := (64 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
			yn += 64
			cy++
			fx1 = fx2
			fy1 = 0
		}
	} else { // Only vertical cell wall intersections remain
		for xn <= bx {
			yw := ay + (xn-ax)*dy/dx
			fy2 := -yn + 64 + yw
			cover := (fy2 - fy1)
			s.cellWaiter.Add(1)
			s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: 64 + fx1, flip: flip}
			xn += 64
			cx++
			fx1 = 0
			fy1 = fy2
		}
	}
	fx2 := bx & (64 - 1)
	fy2 := by & (64 - 1)
	cover := (fy2 - fy1)
	s.cellWaiter.Add(1)
	s.cellWorkers[cy%s.threads].cellChan <- CellM{x: cx, y: cy, cover: cover, AreaOvCov: fx2 + fx1, flip: flip}
}

// SendVerticalLine takes a vertical line that runs in either direction
// and sends Cover and area/cover cell values to the cellWorkers for
// sorting and storage.
func (s *ScannerT) SendVerticalLine(line Line) {
	y1 := int(line.a.Y) >> 6
	y2 := int(line.b.Y) >> 6
	x1 := int(line.a.X) >> 6
	x1f2 := (int(line.a.X) - (x1 << 6)) * 2
	if y1 == y2 {
		cover := int(line.b.Y - line.a.Y)
		s.cellWaiter.Add(1)
		s.cellWorkers[y1%s.threads].cellChan <- CellM{x: x1, y: y1, cover: cover, AreaOvCov: x1f2, flip: 1}
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
	cover := 64 - y1f
	s.cellWaiter.Add(1)
	s.cellWorkers[y1%s.threads].cellChan <- CellM{x: x1, y: y1, cover: cover, AreaOvCov: x1f2, flip: flip}
	for y := y1 + 1; y < y2; y++ {
		s.cellWaiter.Add(1)
		s.cellWorkers[y%s.threads].cellChan <- CellM{x: x1, y: y, cover: 64, AreaOvCov: x1f2, flip: flip}
	}
	cover = y2f
	s.cellWaiter.Add(1)
	s.cellWorkers[y2%s.threads].cellChan <- CellM{x: x1, y: y2, cover: cover, AreaOvCov: x1f2, flip: flip}
}

// CellSaver threadIndex determines the cell worker
// this thread acts on. One CellSaver per
// threadIndex is instantiated, so updates to scanLinks
// slice will not conflict
func (s *ScannerT) CellSaver(threadIndex int) {

	for v := range s.cellWorkers[threadIndex].cellChan {
		// No cover or off the top or bottom so can be ignored

		if v.cover == 0 || v.y < 0 || v.y >= s.height {
			s.cellWaiter.Done()
			continue
		}
		v.cover *= v.flip
		// Pin any segments going out of the side bounds to the edges.
		// Scan line cover sums should still be zero.
		if v.x < 0 {
			v.x = 0
			v.AreaOvCov = v.cover << 6 // area further in gets full area for the cover
		} else if v.x >= s.width {
			v.x = s.width - 1
			v.AreaOvCov = v.cover << 6 // area further in gets full area for the cover
		} else {
			//Now calc the true  area, replacing area/cover value
			v.AreaOvCov = v.AreaOvCov * v.cover
		}

		store := s.cellWorkers[threadIndex].scanLinks
		ic := v.y / s.threads // Find the offset of the link list header sentinel
		var icPrev int
		cc := store[ic]
		// The algorithm expects v.x >= 0
		// as enforced above so, the sentinel x value of -1 is always less
		// than x. The icPrev value is only set if the loop fires
		// at least once, but that is ensured where it is used
		// in the default switch case bellow.
		for cc.x < v.x && cc.yn != -1 {
			icPrev = ic
			ic = cc.yn
			cc = store[ic]
		}
		switch {
		case cc.x == v.x:
			// Cell exists, just add area and cover
			store[ic].area += v.AreaOvCov
			store[ic].cover += v.cover
		case cc.yn == -1 && cc.x < v.x: // Add new cell to end of list, yn = -1 indicates the cell is terminal.
			store[ic].yn = len(store)
			s.cellWorkers[threadIndex].scanLinks = append(store, Cell{x: v.x, yn: -1, area: v.AreaOvCov, cover: v.cover})
		default: //cc.x > v.x thus insert new cell into list between cc and previous cell
			store[icPrev].yn = len(store)
			s.cellWorkers[threadIndex].scanLinks = append(store, Cell{x: v.x, yn: ic, area: v.AreaOvCov, cover: v.cover})
		}
		s.cellWaiter.Done()
	}
}

// SetBounds set the boundaries in which the scanner
// is allowed to draw. Negative values are excluded
func (s *ScannerT) SetBounds(height, width int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	s.width = width
	s.height = height
	for i := 0; i < s.threads; i++ {
		lineCount := height / s.threads
		if i < height%s.threads {
			lineCount++
		}
		s.cellWorkers[i].scanLinks = s.cellWorkers[i].scanLinks[:0]
		s.cellWorkers[i].lineCount = lineCount
		for j := 0; j < lineCount; j++ {
			// Cell.x = -1 means it is a sentinel and cell.yn = -1 means the list is empty.
			s.cellWorkers[i].scanLinks = append(s.cellWorkers[i].scanLinks, Cell{x: -1, yn: -1})
		}
	}
}

// NewScanT returns a multi-threaded implementation of the cl-aa antialiasing algorithm
// ScannerT implements the rasterx.Scannner interface for use with the rasterx and oksvg packages.
// An object implementing the Collector interface must be provided, which will convert x,y, and alpha
// values to the target format such as an image.RGBA. A collector for image.RGBA, RGBACollector, is
// defined in this package.
func NewScanT(threads, width, height int, collector Collector) (s *ScannerT) {
	s = &ScannerT{height: height, width: width, threads: threads, collector: collector,
		extents: make([]fixed.Rectangle26_6, threads), lineChan: make(chan Line, threads*64),
		cellWorkers: make([]CellWorker, threads)}
	for i := 0; i < threads; i++ {
		// Create one linked list start sentinels/place holder for each scan line serviced by this thread.
		// Some threads have one more scan line assigned than the others depending on the divide remainder.
		lineCount := height / threads
		if i < height%threads {
			lineCount++
		}
		s.cellWorkers[i].lineCount = lineCount
		for j := 0; j < lineCount; j++ {
			// Cell.x = -1 means it is a sentinel and cell.yn = -1 means the list is empty.
			s.cellWorkers[i].scanLinks = append(s.cellWorkers[i].scanLinks, Cell{x: -1, yn: -1})
		}
		// Set extents rect to sentinel values
		s.extents[i].Max = fixed.Point26_6{X: fixed.Int26_6(-math.MaxInt32), Y: fixed.Int26_6(-math.MaxInt32)}
		s.extents[i].Min = fixed.Point26_6{X: fixed.Int26_6(math.MaxInt32), Y: fixed.Int26_6(math.MaxInt32)}
		// Each cellWorker holds a chan to receive Cell values in CellSaver
		// and a sweepChan to trigger scanline sweeping in SweepLines.
		s.cellWorkers[i].cellChan = make(chan CellM, 64)
		s.cellWorkers[i].sweepChan = make(chan bool, 64)

		go s.CellSaver(i)      // Listens to cellWorkers[i].cellChan
		go s.Sweep(i)          // Listens to cellWorkers[i].sweepChan
		go s.LineToSegments(i) // Listens to lineChan and sends to cellChans
		// Which cellChan gets sent the cell area and coverage increment is
		// determined by y/s.threads. This way each cell linked list store
		// is generated without conflict from another thread.
	}
	return
}

// Close shuts down the channels associated with the ScannerT
func (s *ScannerT) Close() {
	close(s.lineChan)
	for i := 0; i < s.threads; i++ {
		close(s.cellWorkers[i].cellChan)
		close(s.cellWorkers[i].sweepChan)
	}
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
func (s *ScannerT) Start(a fixed.Point26_6) {
	// if s.inPath {
	// 	if s.firstPoint != s.lastPoint {
	// 		s.Line(s.firstPoint) // close the last path
	// 	}
	// }
	s.firstPoint = a
	s.lastPoint = a
	s.inPath = true
}

// Line adds a straight line segment to the path
func (s *ScannerT) Line(b fixed.Point26_6) {
	s.lineWaiter.Add(1)
	s.lineChan <- Line{a: s.lastPoint, b: b}
	s.lastPoint = b
}

// Draw finishes the path if it is open
// and then sweeps the accumulated area and cover
// cell values to the collector
func (s *ScannerT) Draw() {
	if s.inPath {
		if s.firstPoint != s.lastPoint {
			s.Line(s.firstPoint) // Close the last path
		}
		s.inPath = false
	}
	s.lineWaiter.Wait()
	s.cellWaiter.Wait()
	s.sweepWaiter.Add(s.threads)
	for i := 0; i < s.threads; i++ {
		s.cellWorkers[i].sweepChan <- true
	}
	s.sweepWaiter.Wait()
}

// Clear reinitializes the cell linked lists and
// the path extents to make it ready for new paths
func (s *ScannerT) Clear() {
	// In case Clear is called before Draw completes (should not happen),
	// wait on threads to finish
	s.lineWaiter.Wait()
	s.cellWaiter.Wait()
	s.sweepWaiter.Wait()
	for i := 0; i < s.threads; i++ {
		// Set all the scan linked lists back to empty
		s.cellWorkers[i].scanLinks = s.cellWorkers[i].scanLinks[:s.cellWorkers[i].lineCount]
		for j := range s.cellWorkers[i].scanLinks {
			s.cellWorkers[i].scanLinks[j].yn = -1
		}
		// Set max/min sentinel values for extent rects
		s.extents[i].Max = fixed.Point26_6{X: fixed.Int26_6(-math.MaxInt32), Y: fixed.Int26_6(-math.MaxInt32)}
		s.extents[i].Min = fixed.Point26_6{X: fixed.Int26_6(math.MaxInt32), Y: fixed.Int26_6(math.MaxInt32)}
	}
}

// GetPathExtent returns the bounaries of the current path
func (s *ScannerT) GetPathExtent() fixed.Rectangle26_6 {
	s.lineWaiter.Wait() // These have to finish before the extent can be calculated
	s.cellWaiter.Wait()
	maxRect := s.extents[0]
	for i := 1; i < s.threads; i++ {
		maxRect = Expand(maxRect, s.extents[i])
	}
	return maxRect
}

//SetWinding does nothing for now
func (s *ScannerT) SetWinding(useNonZeroWinding bool) {
}

//SetColor sends either a rasterx.ColorFunc or
// color.Color value to the collector
func (s *ScannerT) SetColor(color interface{}) {
	s.collector.SetColor(color)
}

// SetClip does nothing for now
func (s *ScannerT) SetClip(rect image.Rectangle) {

}
