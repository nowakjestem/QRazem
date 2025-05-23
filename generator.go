package main

import (
   "encoding/base64"
   "fmt"
   "image"
   "image/color"
   "math"
   "path"
   "strings"

   "github.com/skip2/go-qrcode"
)

// generateQR creates a raster QR code image with specified colors and size.
// generateQR creates a raster QR code image with specified colors and size,
// and returns the module count (number of QR modules per side).
func generateQR(code, qrCol, bgCol string, size int) (image.Image, int, error) {
   qrColor, err := parseHexColor(qrCol)
   if err != nil {
       qrColor = color.Black
   }
   bgColor, err := parseHexColor(bgCol)
   if err != nil {
       bgColor = color.White
   }
   qr, err := qrcode.New(code, qrcode.Highest)
   if err != nil {
       return nil, err
   }
   qr.BackgroundColor = bgColor
   qr.ForegroundColor = qrColor
   // Determine module count from QR bitmap
   bitmap := qr.Bitmap()
   modules := len(bitmap)
   return qr.Image(size), modules, nil
}

// generateSVG returns an SVG representation of the QR code, optionally embedding an SVG logo.
func generateSVG(code, qrCol, bgCol string, size int, svgLogo []byte, logoName string) ([]byte, error) {
   fgColor, _ := parseHexColor(qrCol)
   bgColor, _ := parseHexColor(bgCol)
   qr, err := qrcode.New(code, qrcode.Highest)
   if err != nil {
       return nil, err
   }
   qr.BackgroundColor = bgColor
   qr.ForegroundColor = fgColor
   bitmap := qr.Bitmap()
   modules := len(bitmap)
   moduleSize := float64(size) / float64(modules)
   margin := (float64(size) - moduleSize*float64(modules)) / 2

   var sb strings.Builder
   sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
   sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">`, size, size))
   sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s"/>`, bgCol))
   for y := 0; y < modules; y++ {
       for x := 0; x < modules; x++ {
           if bitmap[y][x] {
               xPos := margin + moduleSize*float64(x)
               yPos := margin + moduleSize*float64(y)
               sb.WriteString(fmt.Sprintf(
                   `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`,
                   xPos, yPos, moduleSize, moduleSize, qrCol))
           }
       }
   }
   if len(svgLogo) > 0 {
       // Embed logo (SVG or raster) with background clearing and padding
       enc := base64.StdEncoding.EncodeToString(svgLogo)
       // reserved square for logo: 24% of QR code size
       rawSF := float64(size) * 0.24
       rawSquare := int(math.Ceil(rawSF))
       // padding: 2% of QR code size
       padF := float64(size) * 0.02
       pad := int(math.Ceil(padF))
       // inner max dimension for logo
       innerMax := rawSquare - pad*2
       // top-left of reserved square
       sqOffsetF := (float64(size) - rawSF) / 2
       sqOffset := int(math.Floor(sqOffsetF))
       // logo size and offset inside reserved square
       var logoW, logoH int
       ext := strings.ToLower(path.Ext(logoName))
       // preserve aspect ratio
       if ext == ".svg" {
           logoW = innerMax
           logoH = innerMax
       } else {
           logoW = innerMax
           logoH = innerMax
       }
       logoOffsetX := sqOffset + pad + int(math.Floor(float64(innerMax-logoW)/2))
       logoOffsetY := sqOffset + pad + int(math.Floor(float64(innerMax-logoH)/2))
       // clear background under reserved square
       sb.WriteString(fmt.Sprintf(
           `<rect x="%.0f" y="%.0f" width="%d" height="%d" fill="%s"/>`,
           sqOffsetF, sqOffsetF, rawSquare, rawSquare, bgCol))
       // determine MIME type
       mime := "image/svg+xml"
       if ext == ".png" {
           mime = "image/png"
       } else if ext == ".jpg" || ext == ".jpeg" {
           mime = "image/jpeg"
       }
       // embed logo image
       sb.WriteString(fmt.Sprintf(
           `<image x="%d" y="%d" width="%d" height="%d" href="data:%s;base64,%s" preserveAspectRatio="xMidYMid meet"/>`,
           logoOffsetX, logoOffsetY, logoW, logoH, mime, enc))
   }
   sb.WriteString(`</svg>`)
   return []byte(sb.String()), nil
}