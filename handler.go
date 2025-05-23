package main

import (
   "bytes"
   "encoding/json"
   "image"
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

   var qrReq QRRequest
   var svgData []byte
   var logoName string
   const defaultSize = 1024
   var err error
   var qrImg image.Image

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
           if part.FormName() == "svg_logo" && part.FileName() != "" {
               svgData, _ = io.ReadAll(part)
               logoName = part.FileName()
           } else if part.FormName() == "payload" {
               payload, _ := io.ReadAll(part)
               json.Unmarshal(payload, &qrReq)
           }
       }
       size := qrReq.Size
       if size <= 0 {
           size = defaultSize
       }
       qrImg, err = generateQR(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }
       if len(svgData) > 0 {
           bgCol, err2 := parseHexColor(qrReq.BgColor)
           if err2 != nil {
               bgCol = color.White
           }
           ext := strings.ToLower(path.Ext(logoName))
           if ext == ".svg" {
               qrImg = overlaySVG(qrImg, svgData, 0.24, bgCol)
           } else {
               qrImg = overlayRaster(qrImg, svgData, 0.24, bgCol)
           }
       }
   } else {
       if err := json.NewDecoder(r.Body).Decode(&qrReq); err != nil {
           http.Error(w, "Bad JSON: "+err.Error(), http.StatusBadRequest)
           return
       }
       size := qrReq.Size
       if size <= 0 {
           size = defaultSize
       }
       qrImg, err = generateQR(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }
   }

   format := strings.ToLower(qrReq.Format)
   size := qrReq.Size
   if size <= 0 {
       size = defaultSize
   }

   switch format {
   case "svg":
       svgBytes, err := generateSVG(qrReq.Text, qrReq.QRColor, qrReq.BgColor, size, svgData, logoName)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }
       w.Header().Set("Content-Type", "image/svg+xml")
       w.Write(svgBytes)
   case "jpg", "jpeg":
       w.Header().Set("Content-Type", "image/jpeg")
       buf := new(bytes.Buffer)
       if err := jpeg.Encode(buf, qrImg, &jpeg.Options{Quality: 80}); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
       w.Write(buf.Bytes())
   default:
       w.Header().Set("Content-Type", "image/png")
       buf := new(bytes.Buffer)
       if err := png.Encode(buf, qrImg); err != nil {
           http.Error(w, "encoding error", http.StatusInternalServerError)
           return
       }
       w.Write(buf.Bytes())
   }
}