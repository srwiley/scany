package scany_test

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/srwiley/scany"
	"golang.org/x/image/math/fixed"
)

func BenchmarkGradsS10(b *testing.B) {
	RunGradsS(b, 10)
}

func BenchmarkGradsS50(b *testing.B) {
	RunGradsS(b, 50)
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
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
}

func TestScanSHalfCirc(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/halfCirc.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanS(w, h, collector)
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
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanS(w, h, collector)

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
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanS(w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/TestShapes6.png", img)
}

func DrawSliver(sc rasterx.Scanner) {
	q1 := fixed.Point26_6{X: fixed.Int26_6(64*5 + 13), Y: fixed.Int26_6(64*3 + 30)}
	p1 := fixed.Point26_6{X: fixed.Int26_6(64*4 + 54), Y: fixed.Int26_6(64*4 + 22)}
	p2 := fixed.Point26_6{X: fixed.Int26_6(64*4 + 55), Y: fixed.Int26_6(64*4 + 22)}
	q2 := fixed.Point26_6{X: fixed.Int26_6(64*4 + 60), Y: fixed.Int26_6(64*3 + 32)}
	sc.Start(q1)
	sc.Line(q2)
	sc.Line(p2)
	sc.Line(p1)
	sc.Line(q1)
	sc.Draw()
}

func TestScanSSliver(t *testing.T) {

	width := 7
	height := 7

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))

	collector := &scany.RGBACollector{Image: img1}
	sc := scany.NewScanS(width, height, collector)

	//painter := scanFT.NewRGBAPainter(img1)
	//sc := scanFT.NewScannerFT(width, height, painter)

	sc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	DrawSliver(sc)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			fmt.Print(img1.Pix[i*4+img1.Stride*j], "\t")
		}
		fmt.Println()
	}
}

func TestScanSSquare(t *testing.T) {

	width := 7
	height := 7

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))

	collector := &scany.RGBACollector{Image: img1}
	sc := scany.NewScanS(width, height, collector)

	//painter := scanFT.NewRGBAPainter(img1)
	//sc := scanFT.NewScannerFT(width, height, painter)

	sc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	DrawSquare(sc)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			fmt.Print(img1.Pix[i*4+img1.Stride*j], "\t")
		}
		fmt.Println()
	}
}
