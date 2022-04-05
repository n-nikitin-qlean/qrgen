package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const outputDir = "./qrs"

type QR struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Data struct {
		Title string `json:"title"`
	} `json:"data"`
}

func main() {
	_ = os.RemoveAll(outputDir)
	if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
		log.Fatalf("can not create dir `%s`: %s\n", outputDir, err)
	}

	for _, pack := range getPackages() {
		packageDirPath := fmt.Sprintf("%s/%s", outputDir, pack.title)
		if err := os.Mkdir(packageDirPath, os.ModePerm); err != nil {
			log.Fatalf("can not create dir `%s`: %s\n", packageDirPath, err)
		}

		// qr комплекта
		if err := createSkuV1QRFile(packageDirPath, pack.skuId, pack.title); err != nil {
			log.Fatalf("can not create SKU_V1_QR file (`%s`, %s): %s\n", packageDirPath, pack.title, err)
		}

		// подкаталог для айтемов комплекта
		packageDirPath = fmt.Sprintf("%s/%s", packageDirPath, pack.title)
		if err := os.Mkdir(packageDirPath, os.ModePerm); err != nil {
			log.Fatalf("can not create dir `%s`: %s\n", packageDirPath, err)
		}

		for _, item := range pack.items {
			err := createSkuV1QRFile(packageDirPath, item.skuId, item.title)
			if err != nil {
				log.Fatalf("can not create SKU_V1_QR file (`%s`, %s): %s\n", packageDirPath, item.title, err)
			}
		}
	}
}

func createSkuV1QRFile(fileDirPath string, skuId int, title string) error {
	filePath := fmt.Sprintf("%s/%s (SKU ID: %d).png", strings.TrimRight(fileDirPath, "/"), title, skuId)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	qrCode, _ := json.Marshal(QR{
		Type: "SKU_V1",
		Id:   strconv.Itoa(skuId),
		Data: struct {
			Title string `json:"title"`
		}{
			Title: title,
		},
	})
	err = generateQrCode(f, (string)(qrCode))
	if err != nil {
		return err
	}

	return nil
}

func generateQrCode(file io.Writer, data string) error {
	qrCode, _ := qr.Encode(data, qr.M, qr.Unicode)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	return png.Encode(file, qrCode)
}
