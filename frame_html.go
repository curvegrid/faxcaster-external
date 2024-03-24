package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const frameTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta property="fc:frame" content="vNext" />
    <meta property="fc:frame:post_url" content="{{.PostURL}}" />
    <meta property="fc:frame:image" content="{{.ImageURL}}" />
    <meta property="og:image" content="{{.ImageURL}}" />
    {{range $idx, $element := .Buttons}}
        {{$buttonIdx := inc $idx}}
        <meta property="fc:frame:button:{{$buttonIdx}}" content="{{$element.Label}}" />
        <meta property="fc:frame:button:{{$buttonIdx}}:action" content="{{$element.Action}}" />
        {{if ne $element.Target ""}}
            <meta property="fc:frame:button:{{$buttonIdx}}:target" content="{{$element.Target}}" />
        {{end}}
    {{end}}
    {{if ne .Input ""}}
        <meta property="fc:frame:input:text" content="{{.Input}}" />
    {{end}}
    <title>Faxcaster</title>
</head>
<body>
    <img src="{{.ImageURL}}" alt="{{.ImageAlt}}" />
</body>
</html>
`

type FrameData struct {
	PostURL  string
	ImageURL string
	ImageAlt string
	Frame
}

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

// stripNewlines replaces newlines with spaces in a string.
func stripNewlines(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

func generateFrameHTML(postURL string, frame Frame, values ...interface{}) (string, error) {
	// parse the template with the frameTemplate string
	tmpl, err := template.New("frame").Funcs(funcMap).Parse(frameTemplate)
	if err != nil {
		return "", err
	}

	// print the values to the frame text
	text := fmt.Sprintf(frame.FormatText, values...)

	// generate SVG
	svg, err := textToSVG(text)
	if err != nil {
		return "", fmt.Errorf("failed to generate SVG: %w", err)
	}

	frameData := FrameData{
		PostURL:  postURL,
		ImageURL: svg,
		ImageAlt: stripNewlines(text),
		Frame:    frame,
	}

	// render the template to the ResponseWriter
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, frameData); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
