package main

import (
	"context"
	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v57/common/auth"
	"github.com/oracle/oci-go-sdk/v57/example/helpers"
	"github.com/oracle/oci-go-sdk/v57/objectstorage"
	"io"
	"os"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(generateThumbnails))
}

// bucketName returns the name of the bucket as provided within config.
func bucketName() string {
	name := os.Getenv("BUCKET_NAME")
	if name == "" {
		panic("no bucket name provided")
	}

	return name
}

// mediaId returns the current media ID from the context.
func mediaId(ctx context.Context) string {
	return ctx.Value("MediaId").(string)
}

// generateThumbnails is our main request handler.
func generateThumbnails(ctx context.Context, in io.Reader, out io.Writer) {
	// Generate a media ID.
	mediaCtx := context.WithValue(ctx, "MediaId", RandStringBytesMaskImprSrc(16))

	// Now, attempt to read our supposed image.
	// We only want to read up to 10 megabytes - our imposed maximum.
	// As a bonus, we read directly to the image - format issues will be instantly exposed.
	img := decodeImage(in)

	// Authenticate to obtain a usable client with our storage bucket.
	provider, err := auth.ResourcePrincipalConfigurationProvider()
	helpers.FatalIfError(err)
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

	generateVariants(mediaCtx, client, img)
	out.Write([]byte("We did it!"))
}
