package scanx_test

import (
	"image"
	"testing"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/srwiley/scanx"
)

// Run line to seg tests
// func TestLineSegs(t *testing.T) {
// 	sendCell := func(a, b, c, d int, com string) {
// 		fmt.Println("cx,cy,ar,cv", a, b, c, d, com)
// 	}

// 	// fy := scanx.FindY(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 12)},
// 	// 	fixed.Point26_6{X: fixed.Int26_6(64*4 + 50), Y: fixed.Int26_6(64*4 + 40)}, 64*4)
// 	// fmt.Println("fy", fy, fy&(64-1))

// 	// fyt := scanx.FindY(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 12)},
// 	// 	fixed.Point26_6{X: fixed.Int26_6(64*4 + 50), Y: fixed.Int26_6(64*4 + 40)}, 64*4+50)
// 	// fmt.Println("fyt", fyt, fyt&(64-1))
// 	fmt.Println("East")
// 	scanx.SendEastLine(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 12)},
// 		fixed.Point26_6{X: fixed.Int26_6(64*4 + 50), Y: fixed.Int26_6(64*4 + 40)}, 1, sendCell)

// 	fmt.Println()
// 	fmt.Println("West")

// 	scanx.SendWestLine(fixed.Point26_6{X: fixed.Int26_6(64*4 + 50), Y: fixed.Int26_6(64 + 12)},
// 		fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64*4 + 40)}, 1, sendCell)

// 	fmt.Println()
// 	fmt.Println("East c")
// 	scanx.SendEastLine(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 32)},
// 		fixed.Point26_6{X: fixed.Int26_6(64*3 + 50), Y: fixed.Int26_6(64*3 + 50)}, 1, sendCell)

// 	fmt.Println()
// 	fmt.Println("West c")

// 	scanx.SendWestLine(fixed.Point26_6{X: fixed.Int26_6(64*3 + 32), Y: fixed.Int26_6(64 + 32)},
// 		fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64*3 + 32)}, 1, sendCell)
// 	// scanx.SendWestLine(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 12)},
// 	// 	fixed.Point26_6{X: fixed.Int26_6(64*2 + 50), Y: fixed.Int26_6(64*2 + 40)}, 1, sendCell)

// 	fmt.Println()
// 	fmt.Println("East")

// 	scanx.SendEastLine(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 12)},
// 		fixed.Point26_6{X: fixed.Int26_6(64*2 + 50), Y: fixed.Int26_6(64*2 + 40)}, -1, sendCell)

// 	fmt.Println("West")

// 	scanx.SendWestLine(fixed.Point26_6{X: fixed.Int26_6(64*2 + 50), Y: fixed.Int26_6(64 + 12)},
// 		fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64*2 + 40)}, -1, sendCell)

// 	// scanx.LineToSegs(fixed.Point26_6{X: fixed.Int26_6(64 + 32), Y: fixed.Int26_6(64 + 32)},
// 	// 	fixed.Point26_6{X: fixed.Int26_6(64 + 50), Y: fixed.Int26_6(64 + 50)}, sendCell)

// 	// fmt.Println()

// 	// scanx.LineToSegs(fixed.Point26_6{X: fixed.Int26_6(64 + 25), Y: fixed.Int26_6(64 + 32)},
// 	// 	fixed.Point26_6{X: fixed.Int26_6(64*3 + 50), Y: fixed.Int26_6(64 + 50)}, sendCell)

// }

func TestScanSHalfCirc(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/halfCirc.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/halfCirc.png", img)
}

func TestScanSIcon(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/landscapeIcons/sea.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/sea.png", img)
}

func TestGradsS(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)

	// painter := scanFT.NewRGBAPainter(img)
	// scannerFT := scanFT.NewScannerFT(w, h, painter)
	// rasterM := rasterx.NewDasher(w, h, scannerFT)

	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/TestShapes6.png", img)
}
