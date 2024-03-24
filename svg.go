package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
)

type SVG struct {
	XMLName xml.Name  `xml:"svg"`
	XMLns   string    `xml:"xmlns,attr"`
	Width   string    `xml:"width,attr"`
	Height  string    `xml:"height,attr"`
	ViewBox string    `xml:"viewBox,attr"`
	Text    []SVGText `xml:"text"`
}

type SVGText struct {
	XMLName xml.Name `xml:"text"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Fill    string   `xml:"fill,attr"`
	Style   string   `xml:"style,attr"`
	TSpan   []TSpan  `xml:"tspan"`
}

type TSpan struct {
	XMLName xml.Name `xml:"tspan"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"dy,attr"`
	Space   string   `xml:"xml:space,attr"`
	Content string   `xml:",chardata"`
}

func textToSVG(text string) (string, error) {
	// Dimensions and padding
	padding := 20                        // Padding around the text
	height := 600                        // Fixed height before adding padding
	width := int(float64(height) * 1.91) // Width based on aspect ratio before adding padding

	// Adjust width and height to include padding
	paddedWidth := width + 2*padding
	paddedHeight := height + 2*padding

	// Split text into lines
	lines := strings.Split(text, "\n")

	// Prepare SVG text with tspans for each line
	svgText := SVGText{
		X:     fmt.Sprintf("%d", padding),
		Y:     "0", // Initial Y position set to 0; dy will handle vertical positioning
		Fill:  "black",
		Style: "dominant-baseline:middle;font-size:24px;font-family:'Courier New',Courier,monospace;font-weight:bold;",
	}

	dyOffset := 30 // Fixed vertical offset for each line
	for _, line := range lines {
		// For blank lines, insert a space to ensure the line's presence
		if line == "" {
			line = " "
		}
		tspan := TSpan{
			X:       fmt.Sprintf("%d", padding),
			Y:       fmt.Sprintf("%d", dyOffset), // Use dyOffset for vertical spacing
			Space:   "preserve",                  // Ensure whitespace is preserved
			Content: line,
		}
		svgText.TSpan = append(svgText.TSpan, tspan)
	}

	svg := SVG{
		XMLns:   "http://www.w3.org/2000/svg",
		Width:   fmt.Sprintf("%d", paddedWidth),
		Height:  fmt.Sprintf("%d", paddedHeight),
		ViewBox: fmt.Sprintf("0 0 %d %d", paddedWidth, paddedHeight),
		Text:    []SVGText{svgText},
	}

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(svg); err != nil {
		return "", err
	}

	// Encode the SVG to a data URI
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
