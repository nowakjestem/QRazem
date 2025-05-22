package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "image"
   "image/color"
   "image/draw"
   "image/png"
   _ "image/gif"
   _ "image/jpeg"
   "io"
   "log"
   "net/http"
   "os"
   "path"
   "strconv"
   "strings"
   "encoding/base64"
   "image/jpeg"

   "github.com/nfnt/resize"
   "github.com/skip2/go-qrcode"
   "github.com/srwiley/oksvg"
   "github.com/srwiley/rasterx"
)

// QRRequest - struktura do odczytu JSON z frontend
type QRRequest struct {
   Text    string `json:"text"`
   QRColor string `json:"qr_color"`
   BgColor string `json:"bg_color"`
   // Download format: svg, png, jpg
   Format  string `json:"format"`
   // Image size (px)
   Size    int    `json:"size"`
}

// Build SVG representation of the QR code with optional embedded logo
func generateSVG(code, qrCol, bgCol string, size int, svgLogo []byte, logoName string) ([]byte, error) {
   // Create QR bitmap
   qrColor, _ := parseHexColor(qrCol)
   bgColor, _ := parseHexColor(bgCol)
   qr, err := qrcode.New(code, qrcode.Highest)
   if err != nil {
       return nil, err
   }
   qr.BackgroundColor = bgColor
   qr.ForegroundColor = qrColor
   bitmap := qr.Bitmap()
   modules := len(bitmap)
   // Compute module pixel size and margin
   moduleSize := float64(size) / float64(modules)
   margin := (float64(size) - moduleSize*float64(modules)) / 2
   // Build SVG
   var sb strings.Builder
   sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
   sb.WriteString(fmt.Sprintf("<svg width=\"%d\" height=\"%d\" xmlns=\"http://www.w3.org/2000/svg\">", size, size))
   // background
   sb.WriteString(fmt.Sprintf("<rect width=\"100%%\" height=\"100%%\" fill=\"%s\"/>", bgCol))
   // modules
   fgHex := qrCol
   for y := 0; y < modules; y++ {
       for x := 0; x < modules; x++ {
           if bitmap[y][x] {
               xPos := margin + moduleSize*float64(x)
               yPos := margin + moduleSize*float64(y)
               sb.WriteString(fmt.Sprintf(
                   "<rect x=\"%.2f\" y=\"%.2f\" width=\"%.2f\" height=\"%.2f\" fill=\"%s\"/>",
                   xPos, yPos, moduleSize, moduleSize, fgHex))
           }
       }
   }
   // embed logo if SVG
   if len(svgLogo) > 0 && strings.ToLower(path.Ext(logoName)) == ".svg" {
       // base64 encode raw SVG
       enc := base64.StdEncoding.EncodeToString(svgLogo)
       // place logo centered at 24% area
       rawSide := float64(size) * 0.24
       logoOffset := (float64(size) - rawSide) / 2
       sb.WriteString(fmt.Sprintf(
           "<image x=\"%.2f\" y=\"%.2f\" width=\"%.2f\" height=\"%.2f\" href=\"data:image/svg+xml;base64,%s\" preserveAspectRatio=\"xMidYMid meet\"/>",
           logoOffset, logoOffset, rawSide, rawSide, enc))
   }
   sb.WriteString("</svg>")
   return []byte(sb.String()), nil
}

func parseHexColor(s string) (c color.Color, err error) {
	// dozwolone formaty: #RRGGBB, RRGGBB, #RGB, RGB
	s = strings.TrimPrefix(s, "#")
	if len(s) == 6 {
		r, _ := strconv.ParseInt(s[0:2], 16, 0)
		g, _ := strconv.ParseInt(s[2:4], 16, 0)
		b, _ := strconv.ParseInt(s[4:6], 16, 0)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
	} else if len(s) == 3 {
		r, _ := strconv.ParseInt(strings.Repeat(string(s[0]), 2), 16, 0)
		g, _ := strconv.ParseInt(strings.Repeat(string(s[1]), 2), 16, 0)
		b, _ := strconv.ParseInt(strings.Repeat(string(s[2]), 2), 16, 0)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
	}
	return nil, fmt.Errorf("invalid color: %s", s)
}

// Funkcja generująca kod QR jako image.Image (z kolorem QR i tła)
func generateQR(code, qrCol, bgCol string, size int) (image.Image, error) {
	qrColor, err := parseHexColor(qrCol)
	if err != nil {
		qrColor = color.Black
	}
	bgColor, err := parseHexColor(bgCol)
	if err != nil {
		bgColor = color.White
	}
	// Use highest error correction to accommodate logo overlay
	qr, err := qrcode.New(code, qrcode.Highest)
	if err != nil {
		return nil, err
	}
	qr.BackgroundColor = bgColor
	qr.ForegroundColor = qrColor
	img := qr.Image(size)
	return img, nil
}

// Rysuje SVG logo na środku wygenerowanego QR code, z wyczyszczeniem obszaru pod logiem
func overlaySVG(base image.Image, svgData []byte, scale float64, bgCol color.Color) image.Image {
   // Ustal docelowy rozmiar loga (jako procent szerokości QR), zachowując proporcje SVG
   w, h := base.Bounds().Dx(), base.Bounds().Dy()
   // padding rate (fraction of width) for logo inside square
   const paddingRate = 0.02
   // raw square side: defines area to clear under logo
   rawSquare := int(float64(w) * scale)
   // inner max dimension for logo (square side minus padding on each side)
   innerMax := rawSquare - int(float64(w)*paddingRate*2)
   icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
   if err != nil {
       return base
   }
   // Oryginalne proporcje SVG
   origW := icon.ViewBox.W
   origH := icon.ViewBox.H
   var logoW, logoH int
   if origW > 0 && origH > 0 {
       ratio := origH / origW
       // Preserve aspect ratio within innerMax
       if ratio <= 1 {
           logoW = innerMax
           logoH = int(float64(innerMax) * ratio)
       } else {
           logoH = innerMax
           logoW = int(float64(innerMax) / ratio)
       }
   } else {
       // Fallback to square
       logoW, logoH = innerMax, innerMax
   }
   icon.SetTarget(0, 0, float64(logoW), float64(logoH))
	rgba := image.NewRGBA(image.Rect(0, 0, logoW, logoH))
	scanner := rasterx.NewScannerGV(logoW, logoH, rgba, rgba.Bounds())
	raster := rasterx.NewDasher(logoW, logoH, scanner)
   icon.Draw(raster, 1.0)

	// Pozycja środka
	offsetX := (w - logoW) / 2
	offsetY := (h - logoH) / 2
   dst := image.NewRGBA(base.Bounds())
   // Draw base QR code
   draw.Draw(dst, base.Bounds(), base, image.Point{}, draw.Over)
   // Clear square area under logo to background color (rawSquare side)
   sqOffsetX := (w - rawSquare) / 2
   sqOffsetY := (h - rawSquare) / 2
   draw.Draw(dst, image.Rect(sqOffsetX, sqOffsetY, sqOffsetX+rawSquare, sqOffsetY+rawSquare), &image.Uniform{bgCol}, image.Point{}, draw.Src)
   // Draw logo on top
   draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), rgba, image.Point{}, draw.Over)
   return dst
}

// Rysuje raster (PNG/JPEG) logo na środku QR code, z wyczyszczeniem obszaru pod logiem
func overlayRaster(base image.Image, imgData []byte, scale float64, bgCol color.Color) image.Image {
   w, h := base.Bounds().Dx(), base.Bounds().Dy()
   const paddingRate = 0.02
   rawSquare := int(float64(w) * scale)
   innerMax := rawSquare - int(float64(w)*paddingRate*2)
   img, _, err := image.Decode(bytes.NewReader(imgData))
   if err != nil {
       return base
   }
   // Preserve aspect ratio within innerMax
   origBounds := img.Bounds()
   origW := origBounds.Dx()
   origH := origBounds.Dy()
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
       targetW = uint(innerMax)
       targetH = uint(innerMax)
   }
   scaled := resize.Resize(targetW, targetH, img, resize.Lanczos3)

   // Center position
   logoW := scaled.Bounds().Dx()
   logoH := scaled.Bounds().Dy()
   offsetX := (w - logoW) / 2
   offsetY := (h - logoH) / 2
   dst := image.NewRGBA(base.Bounds())
   // Draw base QR code
   draw.Draw(dst, base.Bounds(), base, image.Point{}, draw.Over)
   // Clear square area under logo to background color
   sqOffsetX := (w - rawSquare) / 2
   sqOffsetY := (h - rawSquare) / 2
   draw.Draw(dst, image.Rect(sqOffsetX, sqOffsetY, sqOffsetX+rawSquare, sqOffsetY+rawSquare), &image.Uniform{bgCol}, image.Point{}, draw.Src)
   // Draw logo on top
   draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), scaled, image.Point{}, draw.Over)
   return dst
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

		var qrReq QRRequest
		var qrImg image.Image
		// Logo upload data: svg bytes and filename
		var svgData []byte
		var logoName string
	// Use larger QR code size for better resilience under logo overlay
	const qrSize = 1024
	// Obsługa: tylko dane JSON (bez logo) lub multipart (z logo)
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, "Can't read multipart: "+err.Error(), 400)
			return
		}
            for {
                part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
               if part.FormName() == "svg_logo" && part.FileName() != "" {
                   svgData, _ = io.ReadAll(part)
                   logoName = part.FileName()
			} else if part.FormName() == "payload" {
				payloadBytes, _ := io.ReadAll(part)
				json.Unmarshal(payloadBytes, &qrReq)
			}
		}
		qrImg, err = generateQR(qrReq.Text, qrReq.QRColor, qrReq.BgColor, qrSize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
           if len(svgData) > 0 {
               // Parse background color for clearing under logo
               bgColColor, err2 := parseHexColor(qrReq.BgColor)
               if err2 != nil {
                   bgColColor = color.White
               }
               // Determine file type by extension
               ext := strings.ToLower(path.Ext(logoName))
               if ext == ".svg" {
                   qrImg = overlaySVG(qrImg, svgData, 0.24, bgColColor)
               } else {
                   qrImg = overlayRaster(qrImg, svgData, 0.24, bgColColor)
               }
           }
	} else {
		// Zwykły JSON, bez SVG
		err := json.NewDecoder(r.Body).Decode(&qrReq)
		if err != nil {
			http.Error(w, "Bad JSON: "+err.Error(), 400)
			return
		}
		qrImg, err = generateQR(qrReq.Text, qrReq.QRColor, qrReq.BgColor, qrSize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
           // Handle download format
           format := strings.ToLower(qrReq.Format)
           size := qrReq.Size
           if size <= 0 {
               size = qrSize
           }
           switch format {
           case "svg":
               svgBytes, err := generateSVG(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size, svgData, logoName)
               if err != nil {
                   http.Error(w, err.Error(), 500)
                   return
               }
               w.Header().Set("Content-Type", "image/svg+xml")
               w.Write(svgBytes)
               return
           case "jpg", "jpeg":
               w.Header().Set("Content-Type", "image/jpeg")
               var imgBuf bytes.Buffer
               if err := jpeg.Encode(&imgBuf, qrImg, &jpeg.Options{Quality: 80}); err != nil {
                   http.Error(w, "encoding error", 500)
                   return
               }
               w.WriteHeader(200)
               w.Write(imgBuf.Bytes())
               return
           default:
               w.Header().Set("Content-Type", "image/png")
               var buf bytes.Buffer
               if err := png.Encode(&buf, qrImg); err != nil {
                   http.Error(w, "encoding error", 500)
                   return
               }
               w.WriteHeader(200)
               w.Write(buf.Bytes())
               return
           }
}

func main() {
	// Register API endpoint
	http.HandleFunc("/api/generate-qr", qrHandler)
	// Serve static files (Vue frontend)
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
