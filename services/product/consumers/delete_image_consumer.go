package consumers

import (
	"context"
	"log"

	"github.com/SomeHowMicroservice/shm-be/services/product/imagekit"
	"github.com/SomeHowMicroservice/shm-be/services/product/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/product/mq"
)

func StartDeleteImageConsumer(mqc *initialization.MQConnection, imagekit imagekit.ImageKitService) {
	if err := mq.ConsumeMessage(mqc.Chann, "image.delete", func(body []byte) error {
		fileID := string(body)
		ctx := context.Background()
		if err := imagekit.DeleteFile(ctx, fileID); err != nil {
			return err
		}

		log.Printf("Xóa hình ảnh thành công: %s", fileID)
		return nil
	}); err != nil {
		log.Printf("Lỗi khởi tạo delete image consumer: %v", err)
	}
}