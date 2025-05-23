package main

import (
   "fmt"
   "image/color"
   "strconv"
   "strings"
)

// parseHexColor parses a hex color string to color.Color.
// Supported formats: "#RRGGBB", "RRGGBB", "#RGB", "RGB".
func parseHexColor(s string) (color.Color, error) {
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