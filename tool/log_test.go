package tool

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
)

func TestLog(t *testing.T) {
	// GetLogs().Info("msg info")
	Vlog.Info("msg info")
}

func TestImg(t *testing.T) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("/Users/jackielee/Desktop/modern-kiln-397411-90473590fcdb.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	filePath := "/Users/jackielee/Desktop/wechat2.png"
	reader, _ := os.Open(filePath)
	image, err := vision.NewImageFromReader(reader)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to recognize text: %v", err)
	}

	result := ""
	for _, annotation := range annotations {
		result += (annotation.Description)
	}
	fmt.Println(result)
}
