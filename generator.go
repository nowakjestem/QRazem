package main

import (
   "encoding/base64"
   "fmt"
   "image"
   "image/color"
   "path"
   "strings"

   "github.com/skip2/go-qrcode"
)

// generateQR creates a raster QR code image with specified colors and size.
func generateQR(code, qrCol, bgCol string, size int) (image.Image, error) {
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
   return qr.Image(size), nil
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
   if len(svgLogo) > 0 && strings.ToLower(path.Ext(logoName)) == ".svg" {
       enc := base64.StdEncoding.EncodeToString(svgLogo)
       // reserved square for logo: 24% of QR code size
       rawSquare := float64(size) * 0.24
       // padding around logo: 2% of QR code size
       padding := float64(size) * 0.02
       // actual logo dimensions inside padding
       logoSize := rawSquare - 2*padding
       // top-left of reserved square
       sqOffset := (float64(size) - rawSquare) / 2
       // top-left of logo inside reserved square
       logoOffset := sqOffset + padding
       // clear background under reserved area
       sb.WriteString(fmt.Sprintf(
           `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`,
           sqOffset, sqOffset, rawSquare, rawSquare, bgCol))
       // embed logo image with padding
       sb.WriteString(fmt.Sprintf(
           `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" href="data:image/svg+xml;base64,%s" preserveAspectRatio="xMidYMid meet"/>`,
           logoOffset, logoOffset, logoSize, logoSize, enc))
   }
   sb.WriteString(`</svg>`)
   return []byte(sb.String()), nil
}