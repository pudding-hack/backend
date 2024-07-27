package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// ConvertBase64ToImage converts a Base64 encoded string to a byte slice representing the image.
func ConvertBase64ToImage(base64Data string) ([]byte, error) {
	var imageFormat string
	if strings.Contains(base64Data, "data:image/png;base64,") {
		imageFormat = "png"
	} else if strings.Contains(base64Data, "data:image/jpeg;base64,") {
		imageFormat = "jpeg"
	} else {
		return nil, errors.New("invalid image format")
	}

	base64Data = strings.Replace(base64Data, fmt.Sprintf("data:image/%s;base64,", imageFormat), "", 1)

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
