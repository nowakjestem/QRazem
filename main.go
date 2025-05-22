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
	// Ustal docelowy rozmiar loga (jako procent szerokości QR)
	w, h := base.Bounds().Dx(), base.Bounds().Dy()
	logoW := int(float64(w) * scale)
	logoH := int(float64(h) * scale)

	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
	if err != nil {
		return base
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
   // Clear area under logo to background color
   draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), &image.Uniform{bgCol}, image.Point{}, draw.Src)
   // Draw logo on top
   draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), rgba, image.Point{}, draw.Over)
   return dst
}

// Rysuje raster (PNG/JPEG) logo na środku QR code, z wyczyszczeniem obszaru pod logiem
func overlayRaster(base image.Image, imgData []byte, scale float64, bgCol color.Color) image.Image {
   w, h := base.Bounds().Dx(), base.Bounds().Dy()
   // Resize logo to percentage of QR width
   maxDim := int(float64(w) * scale)
   img, _, err := image.Decode(bytes.NewReader(imgData))
   if err != nil {
       return base
   }
   // Preserve aspect ratio
   scaled := resize.Resize(uint(maxDim), 0, img, resize.Lanczos3)

   // Center position
   logoW := scaled.Bounds().Dx()
   logoH := scaled.Bounds().Dy()
   offsetX := (w - logoW) / 2
   offsetY := (h - logoH) / 2
   dst := image.NewRGBA(base.Bounds())
   // Draw base QR code
   draw.Draw(dst, base.Bounds(), base, image.Point{}, draw.Over)
   // Clear area under logo to background color
   draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+logoW, offsetY+logoH), &image.Uniform{bgCol}, image.Point{}, draw.Src)
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
	// Use larger QR code size for better resilience under logo overlay
	const qrSize = 1024
	// Obsługa: tylko dane JSON (bez logo) lub multipart (z logo)
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, "Can't read multipart: "+err.Error(), 400)
			return
		}
            var svgData []byte
            var logoName string
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
	w.Header().Set("Content-Type", "image/png")
	var buf bytes.Buffer
	if err := png.Encode(&buf, qrImg); err != nil {
		http.Error(w, "encoding error", 500)
		return
	}
	w.WriteHeader(200)
	w.Write(buf.Bytes())
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
