package helper

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	customerrors "web-service/errors"

	"github.com/nfnt/resize"
)

func ParseProductImage(r *http.Request) (bytes.Buffer, string, error) {
	file, _, err := r.FormFile("productImage")
	if err != nil {
		return bytes.Buffer{}, "", customerrors.NewBadRequestError("error retrieving file", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return bytes.Buffer{}, "", customerrors.NewBadRequestError("error decoding image", err)
	}

	compressedImage, err := compressAndResizeImage(img)
	if err != nil {
		return bytes.Buffer{}, "", customerrors.NewBadRequestError("error compressing and resizing image", err)
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, compressedImage, &jpeg.Options{Quality: 85}) // Adjust quality here
		if err != nil {
			return bytes.Buffer{}, "", customerrors.NewBadRequestError("error encoding compressed image", err)
		}
	case "png":
		err = png.Encode(&buf, compressedImage)
		if err != nil {
			return bytes.Buffer{}, "", customerrors.NewBadRequestError("error encoding compressed image", err)
		}
	default:
		return bytes.Buffer{}, "", customerrors.NewBadRequestError(fmt.Sprintf("unsupported image format: %s", format), nil)
	}

	return buf, format, nil
}

func compressAndResizeImage(img image.Image) (image.Image, error) {
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)
	return resizedImg, nil
}
