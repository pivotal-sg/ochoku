package images

import (
	"bytes"
	"image"

	proto "github.com/pivotal-sg/ochoku/images/proto"
	"golang.org/x/net/context"

	_ "image/jpeg"
	_ "image/png"
)

type ImageService struct{}

func (*ImageService) StoreImage(ctx context.Context, imgData *proto.ImageData, resp *proto.StatusResponse) error {
	buf := &bytes.Buffer{}
	buf.Write(imgData.Image)

	_, _, err := image.Decode(buf)

	if err == image.ErrFormat {
		*resp = proto.StatusResponse{
			Message: "Bad image format, jpeg and png only please.",
			Success: false,
		}
		return nil
	}
	*resp = proto.StatusResponse{
		Message: "Image Saved!",
		Success: true,
	}
	return nil
}
