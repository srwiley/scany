package scany

import (
	"image"
	"image/color"

	"github.com/srwiley/rasterx"
)

const (
	m        = 1<<16 - 1
	q uint32 = 0xFF00
)

// Collector translates the line seweep data to the target format.
// For example a 16 bit alpha only image would use a different
// collector than an RGBA image.
type Collector interface {
	Sweeper(lastX, scanLine, len int, cover int16)
	SetColor(clr interface{})
}

// RGBA collector implements the collector interface
// for image.RGBA images.
type RGBACollector struct {
	Color     color.RGBA
	ColorFunc rasterx.ColorFunc
	Image     *image.RGBA
}

// Sweeper sweeps the accumulated alpha value into
// the collector's RGBA image using Duff-Porter color
// composition equations. If the collector's ColorFunc is
// not nil it uses that otherwise it uses the collectors color.
// len is the number of steps in the y direction the sweep
// extends.
func (r *RGBACollector) Sweeper(x, y, len int, alpha int16) {

	offset := x*4 + r.Image.Stride*y
	ma := uint32(alpha) << 4
	//fmt.Println("ma", ma, alpha)
	if ma > 65535 { // Kind of frustating to have to do this, but
		// cant think of better work around.
		ma = 65535
	}
	if r.ColorFunc == nil {
		// Duff-Porter color composition
		rma := uint32(r.Color.R) * ma
		gma := uint32(r.Color.G) * ma
		bma := uint32(r.Color.B) * ma
		ama := uint32(r.Color.A) * ma
		a := m - (ama / (m >> 8))
		for i := offset; i < offset+len*4; i += 4 {
			r.Image.Pix[i+0] = uint8((uint32(r.Image.Pix[i+0])*a + rma) / q)
			r.Image.Pix[i+1] = uint8((uint32(r.Image.Pix[i+1])*a + gma) / q)
			r.Image.Pix[i+2] = uint8((uint32(r.Image.Pix[i+2])*a + bma) / q)
			r.Image.Pix[i+3] = uint8((uint32(r.Image.Pix[i+3])*a + ama) / q)
		}
	} else {
		for i := offset; i < offset+len*4; i += 4 {
			rcr, rcg, rcb, rca := r.ColorFunc(x, y).RGBA()
			//fmt.Println("col", r.ColorFunc(x, y))
			x++
			dr := uint32(r.Image.Pix[i+0])
			dg := uint32(r.Image.Pix[i+1])
			db := uint32(r.Image.Pix[i+2])
			da := uint32(r.Image.Pix[i+3])
			a := (m - (rca * ma / m)) * 0x101
			//fmt.Println("cr", (dr*a+rcr*ma)/m>>8)
			//fmt.Println("cg", (dg*a+rcg*ma)/m>>8)
			//fmt.Println("cg", (db*a+rcb*ma)/m>>8)
			//fmt.Println("cca", (da*a+rca*ma)/m>>8)
			r.Image.Pix[i+0] = uint8((dr*a + rcr*ma) / m >> 8)
			r.Image.Pix[i+1] = uint8((dg*a + rcg*ma) / m >> 8)
			r.Image.Pix[i+2] = uint8((db*a + rcb*ma) / m >> 8)
			r.Image.Pix[i+3] = uint8((da*a + rca*ma) / m >> 8)
			//fmt.Println("res", r.Image.Pix[i+0], r.Image.Pix[i+1], r.Image.Pix[i+2], r.Image.Pix[i+3])
		}
	}
}

// SetColor accepts either a color.Color or
// a rastserx.ColorFunc
func (r *RGBACollector) SetColor(clr interface{}) {
	switch c := clr.(type) {
	case rasterx.ColorFunc:
		r.ColorFunc = c
	case color.RGBA:
		r.Color = c
		r.ColorFunc = nil
	case color.Color:
		rd, g, b, a := c.RGBA()
		r.Color = color.RGBA{
			R: uint8(rd >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8)}
		r.ColorFunc = nil
	default:
		r.Color = color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0}
		r.ColorFunc = nil
	}
}

// Sweep receives a boolean trigger from the
// sweepChan to sweep the assigned
// thread value and associated scan lines to the collector
// the first value is used to determine if the path is an
// "innie" or an "outtie".
func (s *ScannerT) Sweep(thread int) {
	for range s.cellWorkers[thread].sweepChan {
		store := s.cellWorkers[thread].scanLinks
		// Iterate over each linked list start sentinel
		for i := 0; i < s.cellWorkers[thread].lineCount; i++ {
			// Sweep scan line
			ic := store[i].yn
			if ic != -1 {
				scanLine := i*s.threads + thread
				flip := 1
				cc := store[ic]
				cover := cc.cover
				lastX := cc.x
				val := 64*cover - cc.area/2
				if val < 0 { // all vals in each line should be all pos or negative, but
					// to avoid repeat code, just testing per line for now
					flip = -1
				}
				s.collector.Sweeper(lastX, scanLine, 1, int16(val*flip))
				ic = cc.yn
				for ic != -1 {
					cc := store[ic]
					// fill in gap in cells
					//fmt.Print(cc)
					if cc.x-lastX-1 > 0 { // Fill gap
						s.collector.Sweeper(lastX+1, scanLine, (cc.x - lastX - 1), int16((64*cover)*flip))
					}
					cover += cc.cover
					lastX = cc.x
					s.collector.Sweeper(lastX, scanLine, 1, int16((64*cover-cc.area/2)*flip))
					ic = cc.yn
				}
			}
		}
		s.sweepWaiter.Done()
	}
}
