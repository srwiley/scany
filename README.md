# scany
## About

The scany package is a denovo implementation the CL-AA antigrain algorithm in both a single thread or multi-threaded structure by ScanS and ScanT, respectivly. Both implement the rasterx.Scanner interface, and therefore they can be used with the rasterx and oksvg. The single threaded ScanS is vitually as fast as the ScanFT structure implemented in github.com/swiley/scanFT. Benchmarks vary a bit from run to run but ScanS appears to be within 1 or 2 percent as fast as ScanFT. The difference is that ScanFT is under the Freetype license, since it is a direct port from the C implmentation of Freetype. while ScanS and ScanT are under the less restrictive MIT3 license, so they can be freely incoporated in commercial as well as open source projects.

If ScanT is run on a single thread it is around 50% slower than ScanFT or ScanS. However, as additional threads are added it will equal and then exceed the speed of the single thread implementations. At what point this occurs, depends on the image being rendered. SVG's with gradients, for example, particularly benefit from using multiple threads. Also note that number of threads specified can exceed the number of available CPU cores and still increase speed.

Instead of the painter interface used by Freetype to translate x y coordinates and alpha values to various image formats, scanY uses the scany.Collector interface. A collector for the image.RGBA format is provided in the scany package. Collectors for additional formats must be provided by the user.

## How to use

To use the ScanS single threaded antialiaser with oksvg do the following:

```
img := image.NewRGBA(image.Rect(0, 0, width, height))
collector := &scany.RGBACollector{Image: img}
scanner := scany.NewScanS(width, height, collector)
raster := rasterx.NewDasher(w, h, scanner)
icon.Draw(raster, 1.0) // icon is an oksvg.Icon
```

Similarly to use the ScanT multi-threaded antialiaser do:

```
threads := 6
img := image.NewRGBA(image.Rect(0, 0, width, height))
collector := &scany.RGBACollector{Image: img}
scanner := scany.NewScanT(threads, width, height, collector)
raster := rasterx.NewDasher(w, h, scanner)
icon.Draw(raster, 1.0) // icon is an oksvg.Icon
```

All threads will be started when the NewScanT method is called. These can be explicitly shut down before the struct goes out of scope by executing: 
```
scanner.Close()
```
See the test files for additional examples and benchmarks.

Thanks to [Freepik](http://www.freepik.com) from [Flaticon](https://www.flaticon.com/)
Licensed by [Creative Commons 3.0](http://creativecommons.org/licenses/by/3.0/) for the example icon, and those used in the test/landscapeIcons folder.






