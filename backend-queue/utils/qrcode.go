package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQueueQRCodeBase64 รับ queueID แล้วสร้าง QR Code เป็น Base64 string
func GenerateQueueQRCodeBase64(queueID string) (string, error) {
	// ข้อมูลที่จะเก็บใน QR Code
	data := fmt.Sprintf("http://localhost:3000/api/queue/%s/verify", queueID)

	// สร้าง QR code ขนาด 256x256
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// แปลงเป็น Base64
	qrBase64 := base64.StdEncoding.EncodeToString(png)
	return qrBase64, nil
}
