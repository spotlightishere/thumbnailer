package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v58/objectstorage"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"io/ioutil"
	// Acceptable image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// decodeImage decodes a passed image, up to 10 MiB.
func decodeImage(in io.Reader) image.Image {
	img, _, err := image.Decode(io.LimitReader(in, 10*1024*1024))
	DieIfErr(err)

	return img
}

// encodeImage encodes the passed image as a PNG and returns its bytes.
func encodeImage(img image.Image) *bytes.Buffer {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	DieIfErr(err)
	return buf
}

// generateVariants generates an original sized image, half sized and thumbnail image as a PNG.
func generateVariants(ctx context.Context, client objectstorage.ObjectStorageClient, img image.Image) {
	// Output this current image in its original form PNG.
	original := encodeImage(img)

	uploadImage(ctx, client, "", original)

	// Resize to half its width or height, whatever comes first.
	resizedHalfImg := resize(img, img.Bounds().Size().X/2, img.Bounds().Size().Y/2)
	resizedHalf := encodeImage(resizedHalfImg)
	uploadImage(ctx, client, "-half", resizedHalf)

	// Finally, our thumbnail.
	resizedThumbImg := resize(img, 120, 120)
	resizedThumb := encodeImage(resizedThumbImg)
	uploadImage(ctx, client, "-thumb", resizedThumb)
}

// resize outputs an image with an optimal size.
func resize(origImage image.Image, maxWidth int, maxHeight int) image.Image {
	width := origImage.Bounds().Size().X
	height := origImage.Bounds().Size().Y

	if width > maxWidth {
		height = height * maxWidth / width
		width = maxWidth
	}

	if height > maxHeight {
		width = width * maxHeight / height
		height = maxHeight
	}

	if width != maxWidth && height != maxHeight {
		// No resize needs to occur.
		return origImage
	}

	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), origImage, origImage.Bounds(), draw.Over, nil)
	return newImage
}

// uploadImage uploads the given file to the bucket provider for the given media ID.
func uploadImage(ctx context.Context, client objectstorage.ObjectStorageClient, attributes string, buf *bytes.Buffer) {
	id := mediaId(ctx)
	filename := fmt.Sprintf("%s/%s%s.png", id, id, attributes)

	_, err := client.PutObject(ctx, objectstorage.PutObjectRequest{
		NamespaceName: namespaceName(ctx),
		BucketName:    bucketName(ctx),
		ObjectName:    &filename,
		PutObjectBody: ioutil.NopCloser(buf),
	})
	DieIfErr(err)
}
