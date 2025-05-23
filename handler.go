package main

import (
   "bytes"
   "encoding/json"
   "image/color"
   "image/jpeg"
   "image/png"
   "io"
   "net/http"
   "path"
   "strings"
)

// qrHandler handles QR generation requests (JSON or multipart with logo).
func qrHandler(w http.ResponseWriter, r *http.Request) {
   if r.Method != http.MethodPost {
       http.Error(w, "POST only", http.StatusMethodNotAllowed)
       return
   }

   // Parse request payload (JSON or multipart with optional logo)
   var qrReq QRRequest
   var svgData []byte
   var logoName string
   const defaultSize = 1024
   if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
       reader, err := r.MultipartReader()
       if err != nil {
           http.Error(w, "Can't read multipart: "+err.Error(), http.StatusBadRequest)
           return
       }
       for {
           part, err := reader.NextPart()
           if err == io.EOF {
               break
           }
           switch part.FormName() {
           case "svg_logo":
               if part.FileName() != "" {
                   svgData, _ = io.ReadAll(part)
                   logoName = part.FileName()
               }
           case "payload":
               payload, _ := io.ReadAll(part)
               json.Unmarshal(payload, &qrReq)
           }
       }
   } else {
       if err := json.NewDecoder(r.Body).Decode(&qrReq); err != nil {
           http.Error(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
           return
       }
   }

   // Determine output size
   size := qrReq.Size
   if size <= 0 {
       size = defaultSize
   }

   // Respond in requested format
   format := strings.ToLower(qrReq.Format)
   if format == "svg" {
       // SVG output: generate SVG with embedded logo
       svgBytes, err := generateSVG(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size, svgData, logoName)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }
       w.Header().Set("Content-Type", "image/svg+xml")
       w.Write(svgBytes)
       return
   }

   // Raster output: generate raster QR, then overlay logo
   qrImg, err := generateQR(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size)
   if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }
   // Overlay logo if provided
   if len(svgData) > 0 {
       bgCol, colErr := parseHexColor(qrReq.BgColor)
       if colErr != nil {
           bgCol = color.White
       }
       ext := strings.ToLower(path.Ext(logoName))
       if ext == ".svg" {
           qrImg = overlaySVG(qrImg, svgData, 0.24, bgCol)
       } else {
           qrImg = overlayRaster(qrImg, svgData, 0.24, bgCol)
       }
   }
   // Encode to PNG or JPEG
   var buf bytes.Buffer
   switch format {
   case "jpg", "jpeg":
       w.Header().Set("Content-Type", "image/jpeg")
       if err := jpeg.Encode(&buf, qrImg, &jpeg.Options{Quality: 80}); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
   default:
       w.Header().Set("Content-Type", "image/png")
       if err := png.Encode(&buf, qrImg); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
   }
   w.WriteHeader(http.StatusOK)
   w.Write(buf.Bytes())
}