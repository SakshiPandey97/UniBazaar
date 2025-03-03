package helper

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func CreateMockImage(format string) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, fmt.Errorf("error encoding JPEG: %v", err)
		}
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("error encoding PNG: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return buf.Bytes(), nil
}

func TestParseProductImage_ErrorRetrievingFile(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	_, _, err = ParseProductImage(req)
	if err == nil || !strings.Contains(err.Error(), "error retrieving file") {
		t.Errorf("Expected error retrieving file, but got: %v", err)
	}
}

func TestParseProductImage_ErrorEncodingJPEG(t *testing.T) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := CreateMockImage("jpeg")
	if err != nil {
		t.Fatalf("Error creating mock image: %v", err)
	}

	part, err := writer.CreateFormFile("productImage", "image.jpg")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	_, err = part.Write(file)
	if err != nil {
		t.Fatalf("Error writing file to form: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", &requestBody)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, _, err = ParseProductImage(req)
	if err != nil && strings.Contains(err.Error(), "error encoding compressed image") {
		t.Errorf("Expected error encoding compressed image, but got: %v", err)
	}
}

func TestParseProductImage_ErrorEncodingPNG(t *testing.T) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := CreateMockImage("png")
	if err != nil {
		t.Fatalf("Error creating mock image: %v", err)
	}

	part, err := writer.CreateFormFile("productImage", "image.png")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	_, err = part.Write(file)
	if err != nil {
		t.Fatalf("Error writing file to form: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", &requestBody)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, _, err = ParseProductImage(req)
	if err != nil && strings.Contains(err.Error(), "error encoding compressed image") {
		t.Errorf("Expected error encoding compressed image, but got: %v", err)
	}
}

func TestParseProductImage_UnsupportedFormat(t *testing.T) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("productImage", "image.gif")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	_, err = part.Write([]byte("Invalid GIF Data"))
	if err != nil {
		t.Fatalf("Error writing file to form: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", &requestBody)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, _, err = ParseProductImage(req)
	if err == nil || !strings.Contains(err.Error(), "error decoding image") {
		t.Errorf("Expected unsupported image format error, but got: %v", err)
	}
}
