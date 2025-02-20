package helper

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"web-service/model"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func GetUserID(userId string) (int, error) {
	userID, err := strconv.Atoi(userId)
	if err != nil {
		return 0, errors.New("invalid userId, must be an integer")
	}
	return userID, nil
}

func ParseFormAndCreateProduct(r *http.Request) (model.Product, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return model.Product{}, fmt.Errorf("failed to parse form data: %w", err)
	}

	product := model.Product{
		ProductID:          r.FormValue("productId"),
		ProductTitle:       r.FormValue("productTitle"),
		ProductDescription: r.FormValue("productDescription"),
		ProductPostDate:    r.FormValue("productPostDate"),
		ProductLocation:    r.FormValue("productLocation"),
		ProductImage:       r.FormValue("productImage"),
	}

	if product.ProductID == "" {
		product.ProductID = uuid.NewString()
	}

	if err := parseNumericalFormValues(r, &product); err != nil {
		log.Printf("Error parsing numerical form values: %v", err)
		return model.Product{}, err
	}

	if err := product.Validate(); err != nil {
		log.Printf("Product validation failed: %v", err)
		return model.Product{}, err
	}

	return product, nil
}

func parseNumericalFormValues(r *http.Request, product *model.Product) error {
	if condition := r.FormValue("productCondition"); condition != "" {
		if _, err := fmt.Sscanf(condition, "%d", &product.ProductCondition); err != nil {
			return fmt.Errorf("invalid product condition: %v", err)
		}
	}

	if price := r.FormValue("productPrice"); price != "" {
		if _, err := fmt.Sscanf(price, "%f", &product.ProductPrice); err != nil {
			return fmt.Errorf("invalid product price: %v", err)
		}
	}

	return nil
}

func ParseProductImage(r *http.Request) (bytes.Buffer, string, error) {
	file, _, err := r.FormFile("productImage")
	if err != nil {
		return bytes.Buffer{}, "", fmt.Errorf("error retrieving file: %v", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return bytes.Buffer{}, "", fmt.Errorf("error decoding image: %v", err)
	}

	compressedImage, err := compressAndResizeImage(img)
	if err != nil {
		return bytes.Buffer{}, "", fmt.Errorf("error compressing and resizing image: %v", err)
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, compressedImage, &jpeg.Options{Quality: 85}) // Adjust quality here
		if err != nil {
			return bytes.Buffer{}, "", fmt.Errorf("error encoding compressed image: %v", err)
		}
	case "png":
		err = png.Encode(&buf, compressedImage)
		if err != nil {
			return bytes.Buffer{}, "", fmt.Errorf("error encoding compressed image: %v", err)
		}
	default:
		return bytes.Buffer{}, "", fmt.Errorf("unsupported image format: %s", format)
	}

	return buf, format, nil
}

func compressAndResizeImage(img image.Image) (image.Image, error) {
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)
	return resizedImg, nil
}
