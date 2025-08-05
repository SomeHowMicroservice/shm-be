package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/imagekit"
	"github.com/SomeHowMicroservice/shm-be/services/product/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/product/mq"
)

func startUploadImageConsumer(mqc *initialization.MQConnection, imagekit imagekit.ImageKitService) {
	if err := mq.ConsumeMessage(mqc.Chann, "image.upload", func(body []byte) error {
		var imageMsg common.Base64UploadRequest
		if err := json.Unmarshal(body, &imageMsg); err != nil {
			return fmt.Errorf("unmarshal json thất bại: %w", err)
		}
		
		ctx := context.Background()
		res, err := imagekit.UploadFromBase64(ctx, &imageMsg)
		if err != nil {
			return fmt.Errorf("tải ảnh thất bại: %w", err)
		}

		log.Printf("Đã tải lên hình ảnh thành công: %s", res.URL)
		return nil
	}); err != nil {
		log.Printf("Lỗi khởi tạo upload image consumer: %v", err)
	}
}
