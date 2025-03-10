package helper

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"web-service/model"

	"github.com/google/uuid"
)

func GetUserID(userId string) (int, error) {
	userID, err := strconv.Atoi(userId)
	if err != nil {
		return 0, errors.New("invalid userId, must be an integer")
	}
	return userID, nil
}

func ParseLimit(limitStr string) int {
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			return parsedLimit
		}
	}
	return 10
}

func ParseFormAndCreateProduct(r *http.Request, userId int) (model.Product, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return model.Product{}, fmt.Errorf("failed to parse form data: %w", err)
	}

	product := model.Product{
		UserID:             userId,
		ProductID:          r.FormValue("productId"),
		ProductTitle:       r.FormValue("productTitle"),
		ProductDescription: r.FormValue("productDescription"),
		ProductLocation:    r.FormValue("productLocation"),
		ProductImage:       r.FormValue("productImage"),
	}

	if productPostDate := r.FormValue("productPostDate"); productPostDate != "" {
		parsedDate, err := time.Parse("01-02-2006", productPostDate)
		if err != nil {
			return model.Product{}, fmt.Errorf("invalid product post date format: %v", err)
		}
		product.ProductPostDate = parsedDate
	} else {
		return model.Product{}, fmt.Errorf("product post date is required")
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
