package main

import (
   "bytes"
   "image"

   "github.com/srwiley/oksvg"
   "github.com/srwiley/rasterx"
)

// rasterizeSVG converts SVG data into a raster image of given size.
func rasterizeSVG(svgData []byte, size int) (image.Image, error) {
   icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
   if err != nil {
       return nil, err
   }
   // Set rendering target to full canvas
   icon.SetTarget(0, 0, float64(size), float64(size))
   // Prepare RGBA canvas
   rgba := image.NewRGBA(image.Rect(0, 0, size, size))
   // Rasterize SVG into RGBA
   scanner := rasterx.NewScannerGV(size, size, rgba, rgba.Bounds())
   raster := rasterx.NewDasher(size, size, scanner)
   icon.Draw(raster, 1.0)
   return rgba, nil
}