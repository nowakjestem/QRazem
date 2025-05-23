package main

// QRRequest represents the JSON payload for QR generation.
type QRRequest struct {
   Text    string `json:"text"`
   QRColor string `json:"qr_color"`
   BgColor string `json:"bg_color"`
   // Download format: svg, png, jpg
   Format  string `json:"format"`
   // Image size (px)
   Size    int    `json:"size"`
}