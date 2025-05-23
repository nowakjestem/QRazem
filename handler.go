package main

import (
   "bytes"
   "encoding/json"
   "image/jpeg"
   "image/png"
   "io"
   "net/http"
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

   // Generate SVG (with optional embedded logo)
   svgBytes, err := generateSVG(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size, svgData, logoName)
   if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }

   // Respond in requested format
   format := strings.ToLower(qrReq.Format)
   if format == "svg" {
       w.Header().Set("Content-Type", "image/svg+xml")
       w.Write(svgBytes)
       return
   }

   // Raster formats: convert SVG to image.Image
   img, err := rasterizeSVG(svgBytes, size)
   if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }
   switch format {
   case "jpg", "jpeg":
       w.Header().Set("Content-Type", "image/jpeg")
       var buf bytes.Buffer
       if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
       w.WriteHeader(http.StatusOK)
       w.Write(buf.Bytes())
   default:
       w.Header().Set("Content-Type", "image/png")
       var buf bytes.Buffer
       if err := png.Encode(&buf, img); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
       w.WriteHeader(http.StatusOK)
       w.Write(buf.Bytes())
   }
}