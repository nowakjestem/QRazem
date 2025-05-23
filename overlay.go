package main

import (
   "bytes"
   "image"
   "image/color"
   "image/draw"
   _ "image/gif"
   _ "image/png"
   _ "image/jpeg"
   "math"

   "github.com/nfnt/resize"
   "github.com/srwiley/oksvg"
   "github.com/srwiley/rasterx"
)

// overlaySVG draws an SVG logo centered on the base QR code image.
func overlaySVG(base image.Image, svgData []byte, scale float64, bgCol color.Color) image.Image {
   w, h := base.Bounds().Dx(), base.Bounds().Dy()
   const paddingRate = 0.02
   // Compute reserved square size (scale) and padding (percentage of width)
   rawSF := float64(w) * scale
   rawSquare := int(math.Ceil(rawSF))
   pad := int(math.Ceil(float64(w) * paddingRate))
   innerMax := rawSquare - pad*2

   icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
   if err != nil {
       return base
   }
   origW := icon.ViewBox.W
   origH := icon.ViewBox.H
   var logoW, logoH int
   if origW > 0 && origH > 0 {
       ratio := origH / origW
       if ratio <= 1 {
           logoW = innerMax
           logoH = int(math.Round(float64(innerMax) * ratio))
       } else {
           logoH = innerMax
           logoW = int(math.Round(float64(innerMax) / ratio))
       }
   } else {
       logoW, logoH = innerMax, innerMax
   }
   icon.SetTarget(0, 0, float64(logoW), float64(logoH))
   rgba := image.NewRGBA(image.Rect(0, 0, logoW, logoH))
   scanner := rasterx.NewScannerGV(logoW, logoH, rgba, rgba.Bounds())
   raster := rasterx.NewDasher(logoW, logoH, scanner)
   icon.Draw(raster, 1.0)

   dst := image.NewRGBA(base.Bounds())
   draw.Draw(dst, base.Bounds(), base, image.Point{}, draw.Over)
   // Determine reserved square position and clear background
   sqOffsetX := int(math.Floor((float64(w) - rawSF) / 2))
   sqOffsetY := int(math.Floor((float64(h) - rawSF) / 2))
   draw.Draw(dst,
       image.Rect(sqOffsetX, sqOffsetY, sqOffsetX+rawSquare, sqOffsetY+rawSquare),
       &image.Uniform{bgCol}, image.Point{}, draw.Src)
   // Center logo inside reserved area with padding
   offsetX := sqOffsetX + pad + int(math.Floor(float64(innerMax-logoW)/2))
   offsetY := sqOffsetY + pad + int(math.Floor(float64(innerMax-logoH)/2))
   draw.Draw(dst,
       image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH),
       rgba, image.Point{}, draw.Over)
   return dst
}

// overlayRaster draws a raster logo centered on the base QR code image.
func overlayRaster(base image.Image, imgData []byte, scale float64, bgCol color.Color) image.Image {
   w, h := base.Bounds().Dx(), base.Bounds().Dy()
   const paddingRate = 0.02
   // Compute padding and reserved square size, rounding up for full coverage
   padF := float64(w) * paddingRate
   pad := int(math.Ceil(padF))
   rawSF := float64(w) * scale
   rawSquare := int(math.Ceil(rawSF))
   innerMax := rawSquare - pad*2

   img, _, err := image.Decode(bytes.NewReader(imgData))
   if err != nil {
       return base
   }
   orig := img.Bounds()
   origW, origH := orig.Dx(), orig.Dy()
   var targetW, targetH uint
   if origW > 0 && origH > 0 {
       ratio := float64(origH) / float64(origW)
       if ratio <= 1 {
           targetW = uint(innerMax)
           targetH = uint(float64(innerMax) * ratio)
       } else {
           targetH = uint(innerMax)
           targetW = uint(float64(innerMax) / ratio)
       }
   } else {
       targetW, targetH = uint(innerMax), uint(innerMax)
   }
   scaled := resize.Resize(targetW, targetH, img, resize.Lanczos3)

   dst := image.NewRGBA(base.Bounds())
   draw.Draw(dst, base.Bounds(), base, image.Point{}, draw.Over)
   // Clear background under logo (reserved square)
   sqOffsetX := int(math.Floor((float64(w) - rawSF) / 2))
   sqOffsetY := int(math.Floor((float64(h) - rawSF) / 2))
   draw.Draw(dst,
       image.Rect(sqOffsetX, sqOffsetY, sqOffsetX+rawSquare, sqOffsetY+rawSquare),
       &image.Uniform{bgCol}, image.Point{}, draw.Src)
   // Draw logo centered
   logoW, logoH := scaled.Bounds().Dx(), scaled.Bounds().Dy()
   offsetX := int(math.Floor((float64(w) - float64(logoW)) / 2))
   offsetY := int(math.Floor((float64(h) - float64(logoH)) / 2))
   draw.Draw(dst,
       image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH),
       scaled, image.Point{}, draw.Over)
   return dst
}