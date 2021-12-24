package scany_test

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/srwiley/scany"
	"golang.org/x/image/math/fixed"
)

func FilePathWalkDir(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) (e error) {
		if !info.IsDir() && strings.HasSuffix(path, ".svg") {
			files = append(files, path)
		}
		return
	})
	return
}

func ReadIconSet(paths []string) (icons []*oksvg.SvgIcon) {
	for _, p := range paths {
		icon, errSvg := oksvg.ReadIcon(p, oksvg.IgnoreErrorMode)
		if errSvg == nil {
			icons = append(icons, icon)
		}
	}
	return
}

func SaveToPngFile(filePath string, m image.Image) error {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	// Create Writer from file
	b := bufio.NewWriter(f)
	// Write the image into the buffer
	err = png.Encode(b, m)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func DrawSquare(sc rasterx.Scanner) {
	sc.Start(fixed.Point26_6{X: fixed.Int26_6(64*2 + 32), Y: fixed.Int26_6(64*2 + 32)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*5 + 32), Y: fixed.Int26_6(64*2 + 32)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*5 + 32), Y: fixed.Int26_6(64*5 + 32)})
	sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*2 + 32), Y: fixed.Int26_6(64*5 + 32)})
	//sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*2 + 32), Y: fixed.Int26_6(64*2 + 32)})
	//sc.Line(fixed.Point26_6{X: fixed.Int26_6(64*4 + 20), Y: fixed.Int26_6(64*3 + 10)})
	//fmt.Println("extents", sc.GetPathExtent())
	sc.Draw()
}

func TestScanHalfCirc(t *testing.T) {
	icon, errSvg := oksvg.ReadIcon("testdata/svg/halfCirc.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		t.Error(errSvg)
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanT(1, w, h, collector)
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
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanT(1, w, h, collector)
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
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanT(10, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	icon.Draw(rasterM, 1.0)
	SaveToPngFile("testdata/TestShapes6.png", img)
}

func TestScanSquare(t *testing.T) {

	width := 7
	height := 7

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))

	collector := &scany.RGBACollector{Image: img1}

	sc := scany.NewScanT(1, width, height, collector)

	sc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	DrawSquare(sc)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			fmt.Print(img1.Pix[i*4+img1.Stride*j], "\t")
		}
		fmt.Println()
	}
}

func BenchmarkM1Spanner10(b *testing.B) {
	RunScannerT(b, 10, 1)
}

func BenchmarkM4Spanner10(b *testing.B) {
	RunScannerT(b, 10, 4)
}

func BenchmarkM8Spanner10(b *testing.B) {
	RunScannerT(b, 10, 8)
}
func BenchmarkM16Spanner10(b *testing.B) {
	RunScannerT(b, 10, 16)
}
func BenchmarkM32Spanner10(b *testing.B) {
	RunScannerT(b, 10, 32)
}

func BenchmarkM1Spanner50(b *testing.B) {
	RunScannerT(b, 50, 1)
}

func BenchmarkM4Spanner50(b *testing.B) {
	RunScannerT(b, 50, 4)
}

func BenchmarkM8Spanner50(b *testing.B) {
	RunScannerT(b, 50, 8)
}
func BenchmarkM16Spanner50(b *testing.B) {
	RunScannerT(b, 50, 16)
}

func BenchmarkM32Spanner50(b *testing.B) {
	RunScannerT(b, 50, 32)
}

func RunScannerT(b *testing.B, mult, threads int) {
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
		collector  = &scany.RGBACollector{Image: img}
		scanM      = scany.NewScanT(threads, w, h, collector)
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
	RunGradsT(b, 10, 1)
}

func BenchmarkGradsM2(b *testing.B) {
	RunGradsT(b, 10, 2)
}

func BenchmarkGradsM4(b *testing.B) {
	RunGradsT(b, 10, 4)
}

func BenchmarkGradsM8(b *testing.B) {
	RunGradsT(b, 10, 8)
}

func BenchmarkGradsM16(b *testing.B) {
	RunGradsT(b, 10, 16)
}

func BenchmarkGradsM32(b *testing.B) {
	RunGradsT(b, 10, 32)
}

func BenchmarkGradsM150(b *testing.B) {
	RunGradsT(b, 50, 1)
}

func BenchmarkGradsM250(b *testing.B) {
	RunGradsT(b, 50, 2)
}

func BenchmarkGradsM450(b *testing.B) {
	RunGradsT(b, 50, 4)
}

func BenchmarkGradsM850(b *testing.B) {
	RunGradsT(b, 50, 8)
}

func BenchmarkGradsM1650(b *testing.B) {
	RunGradsT(b, 50, 16)
}

func BenchmarkGradsM3250(b *testing.B) {
	RunGradsT(b, 50, 32)
}

func RunGradsT(b *testing.B, mult, threads int) {
	b.StopTimer()
	icon, errSvg := oksvg.ReadIcon("testdata/svg/TestShapes6.svg", oksvg.WarnErrorMode)
	if errSvg != nil {
		b.Error(errSvg)
	}
	wi, hi := int(icon.ViewBox.W), int(icon.ViewBox.H)
	w, h := wi*mult/10, hi*mult/10
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	collector := &scany.RGBACollector{Image: img}
	scanM := scany.NewScanT(threads, w, h, collector)
	rasterM := rasterx.NewDasher(w, h, scanM)
	b.StartTimer()
	for i := 0; i < 50; i++ {
		icon.SetTarget(0.0, 0.0, float64(w), float64(h))
		icon.Draw(rasterM, 1.0)
	}
	b.StopTimer()
	scanM.Close()
}
