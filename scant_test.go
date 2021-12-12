package scanx_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/srwiley/scanFT"
	"github.com/srwiley/scanx"
	"golang.org/x/image/math/fixed"
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

func DrawSquare(sc *scanx.ScannerT) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*4 + 22), Y: fixed.Int26_6(64*3 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*18 + 22), Y: fixed.Int26_6(64*3 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*18 + 22), Y: fixed.Int26_6(64*15 + 15)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*4 + 20), Y: fixed.Int26_6(64*15 + 15)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*4 + 20), Y: fixed.Int26_6(64*3 + 10)})
	sc.Draw()
}

func DrawTriangle(sc *scanx.ScannerT) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*10 + 10)})

	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*4 + 22), Y: fixed.Int26_6(64*10 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})

	sc.Draw()
}

func DrawTriangle2(sc *scanx.ScannerT) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*18 + 22), Y: fixed.Int26_6(64*10 + 10)})

	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*10 + 10)})
	//sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})

	sc.Draw()
}

func DrawProblematic(sc *scanx.ScannerT) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*224 + 63), Y: fixed.Int26_6(64*235 + 22)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*225 + 48), Y: fixed.Int26_6(64*234 + 46)})

	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*226 + 34), Y: fixed.Int26_6(64*234 + 06)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*227 + 21), Y: fixed.Int26_6(64*233 + 31)})

	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*284 + 42), Y: fixed.Int26_6(64*233 + 31)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*285 + 29), Y: fixed.Int26_6(64*234 + 06)})

	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*286 + 15), Y: fixed.Int26_6(64*234 + 46)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*287 + 0), Y: fixed.Int26_6(64*235 + 22)})
	sc.Draw()
}

func DrawDiamond(sc *scanx.ScannerT) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*18 + 22), Y: fixed.Int26_6(64*10 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*18 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*4 + 22), Y: fixed.Int26_6(64*10 + 10)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*10 + 22), Y: fixed.Int26_6(64*4 + 10)})
	sc.Draw()
}

func TestScanHalfCirc(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/halfCirc.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanT(1, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/halfCirc.png", img)
}

func TestScanIcon(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/landscapeIcons/sea.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanT(1, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/sea.png", img)
}

func TestGrads(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// collector := &scanx.RGBACollector{Image: img}
	// scanM := scanx.NewScan(1, w, h, collector)
	// rasterM := rasterx.NewDasher(w, h, scanM)

	painter := scanFT.NewRGBAPainter(img)
	scannerFT := scanFT.NewScannerFT(w, h, painter)
	rasterM := rasterx.NewDasher(w, h, scannerFT)

	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/TestShapes6.png", img)
}

func TestScanToImage(t *testing.T) {

	width := 30
	height := 30

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))

	collector := &scanx.RGBACollector{Image: img1}

	sc := scanx.NewScanT(1, width, height, collector)
	sc.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	DrawSquare(sc)

	sc.Close()
	SaveToPngFile("testdata/square.png", img1)
}

// func BenchmarkM1Spanner10(b *testing.B) {
// 	RunMScanner(b, 10, 1)
// }

// func BenchmarkM4Spanner10(b *testing.B) {
// 	RunMScanner(b, 10, 4)
// }

// func BenchmarkM8Spanner10(b *testing.B) {
// 	RunMScanner(b, 10, 8)
// }
// func BenchmarkM16Spanner10(b *testing.B) {
// 	RunMScanner(b, 10, 16)
// }
// func BenchmarkM32Spanner10(b *testing.B) {
// 	RunMScanner(b, 10, 32)
// }

// func BenchmarkM1Spanner50(b *testing.B) {
// 	RunMScanner(b, 50, 1)
// }

// func BenchmarkM4Spanner50(b *testing.B) {
// 	RunMScanner(b, 50, 4)
// }

// func BenchmarkM8Spanner50(b *testing.B) {
// 	RunMScanner(b, 50, 8)
// }
// func BenchmarkM16Spanner50(b *testing.B) {
// 	RunMScanner(b, 50, 16)
// }

// func BenchmarkM32Spanner50(b *testing.B) {
// 	RunMScanner(b, 50, 32)
// }

func RunMScanner(b *testing.B, mult, threads int) {
	b.StopTimer()
	beachIconNames, err := FilePathWalkDir("testdata/svg/landscapeIcons")
	if err != nil {
		b.Log("cannot walk file path testdata/svg/landscapeIcons")
		b.FailNow()
	}
	var (
		beachIcons = ReadIconSet(beachIconNames)
		wi, hi     = int(beachIcons[0].ViewBox.W), int(beachIcons[0].ViewBox.H)
		w, h       = wi * mult / 10, hi * mult / 10
		bounds     = image.Rect(0, 0, w, h)
		img        = image.NewRGBA(bounds)
		collector  = &scanx.RGBACollector{Image: img}
		scanM      = scanx.NewScanT(threads, w, h, collector)
		rasterM    = rasterx.NewDasher(w, h, scanM)
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, ic := range beachIcons {
			ic.SetTarget(0.0, 0.0, float64(bounds.Max.X), float64(bounds.Max.Y))
			ic.Draw(rasterM, 1.0)
		}
	}
	b.StopTimer()
	scanM.Close()
}

func BenchmarkGradsM1(b *testing.B) {
	RunGradsM(b, 10, 1)
}

func BenchmarkGradsM2(b *testing.B) {
	RunGradsM(b, 10, 2)
}

func BenchmarkGradsM4(b *testing.B) {
	RunGradsM(b, 10, 4)
}

func BenchmarkGradsM8(b *testing.B) {
	RunGradsM(b, 10, 8)
}

func BenchmarkGradsM16(b *testing.B) {
	RunGradsM(b, 10, 16)
}

func BenchmarkGradsM32(b *testing.B) {
	RunGradsM(b, 10, 32)
}

func BenchmarkGradsS10(b *testing.B) {
	RunGradsS(b, 10)
}

func BenchmarkGradsFT10(b *testing.B) {
	RunGradsFT(b, 10)
}

func BenchmarkGradsM150(b *testing.B) {
	RunGradsM(b, 50, 1)
}

func BenchmarkGradsM250(b *testing.B) {
	RunGradsM(b, 50, 2)
}

func BenchmarkGradsM450(b *testing.B) {
	RunGradsM(b, 50, 4)
}

func BenchmarkGradsM850(b *testing.B) {
	RunGradsM(b, 50, 8)
}

func BenchmarkGradsM1650(b *testing.B) {
	RunGradsM(b, 50, 16)
}

func BenchmarkGradsM3250(b *testing.B) {
	RunGradsM(b, 50, 32)
}

func BenchmarkGradsMI3250(b *testing.B) {
	b.StopTimer()
	mult := 50
	threads := 320
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		b.Error(errSvg)
	}
	wi, hi := int(icon.ViewBox.W), int(icon.ViewBox.H)
	w, h := wi*mult/10, hi*mult/10
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanT(threads, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
	b.StopTimer()
	scanM.Close()
}

func BenchmarkGradsFT50(b *testing.B) {
	RunGradsFT(b, 50)
}

func BenchmarkGradsS50(b *testing.B) {
	RunGradsS(b, 50)
}

func RunGradsM(b *testing.B, mult, threads int) {
	b.StopTimer()
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		b.Error(errSvg)
	}
	wi, hi := int(icon.ViewBox.W), int(icon.ViewBox.H)
	w, h := wi*mult/10, hi*mult/10
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanT(threads, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
	b.StopTimer()
	scanM.Close()
}

func RunGradsS(b *testing.B, mult int) {
	b.StopTimer()
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		b.Error(errSvg)
	}
	wi, hi := int(icon.ViewBox.W), int(icon.ViewBox.H)
	w, h := wi*mult/10, hi*mult/10
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scanx.RGBACollector{Image: img}
	scanM := scanx.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
}

func RunGradsFT(b *testing.B, mult int) {
	b.StopTimer()
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		b.Error(errSvg)
	}
	wi, hi := int(icon.ViewBox.W), int(icon.ViewBox.H)
	w, h := wi*mult/10, hi*mult/10
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	painter := scanFT.NewRGBAPainter(img)
	scannerFT := scanFT.NewScannerFT(w, h, painter)
	rasterM := rasterx.NewDasher(w, h, scannerFT)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
}

// func TestScanTriangle(t *testing.T) {
// 	pad := [20][20]int{}

// 	sc := scanx.NewScan(4, 20, 20, collector)

// 	DrawTriangle(sc)
// 	sc.Clear()

// 	for i := range pad {
// 		for j := range pad[0] {
// 			fmt.Print("\t", pad[i][j])
// 		}
// 		fmt.Println()
// 	}
// 	sc.Close()
// }

// Scan a simple square pattern
// func TestScanSquare(t *testing.T) {
// 	pad := [20][20]int{}

// 	sc := scanx.NewScan(4, 20, 20, collector)

// 	DrawSquare(sc)
// 	sc.Clear()

// 	for i := range pad {
// 		for j := range pad[0] {
// 			fmt.Print("\t", pad[i][j])
// 		}
// 		fmt.Println()
// 	}
// 	sc.Close()

//}

// func TestScanDiamond(t *testing.T) {
// 	pad := [20][20]int{}

// 	sc := scanx.NewScan(4, 20, 20, func(x, y int, a int16) {
// 		pad[y][x] = int(a) >> 6
// 	}, nil)

// 	DrawDiamond(sc)
// 	sc.Clear()

// 	for i := range pad {
// 		for j := range pad[0] {
// 			fmt.Print("\t", pad[i][j])
// 		}
// 		fmt.Println()
// 	}
// 	sc.Close()

// }

// // Scan a simple square pattern twice to make sure the clear function is working
// func TestScanTRep(t *testing.T) {
// 	pad1 := [20][20]int{}

// 	sc := scanx.NewScan(4, 20, 20, func(x, y int, a int16) {
// 		pad1[y][x] = int(a) >> 6
// 	}, nil)

// 	DrawSquare(sc)
// 	sc.Clear()

// 	pad2 := [20][20]int{}
// 	af2 := func(x, y int, a int16) {
// 		pad2[x][y] = int(a) >> 6
// 	}
// 	sc.AlphaFunc = af2
// 	//Square is a direct repeat of the lines above
// 	DrawSquare(sc)

// 	//compare the results
// 	for i := range pad1 {
// 		for j := range pad1[i] {
// 			if pad1[i][j] != pad2[i][j] {
// 				t.Error("pad1 and pad 2 mismatch", i, j, pad1[i][j], pad2[i][j])
// 			}
// 		}
// 	}
// 	sc.Close()

// }
