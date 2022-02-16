package main

import (
	"context"
	"encoding/json"
	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/v58/common/auth"
	"github.com/oracle/oci-go-sdk/v58/objectstorage"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(generateThumbnails))
}

// bucketName returns the name of the bucket as provided within config.
func bucketName(ctx context.Context) *string {
	name := fdk.GetContext(ctx).Config()["BUCKET_NAME"]
	if name == "" {
		panic("no bucket name provided")
	}

	return &name
}

// namespaceName retrieves the current namespace from the context.
func namespaceName(ctx context.Context) *string {
	name := ctx.Value("NamespaceName").(string)
	return &name
}

// mediaId returns the current media ID from the context.
func mediaId(ctx context.Context) string {
	return ctx.Value("MediaId").(string)
}

// ResponseFormat dictates what we'll write back on success.
type ResponseFormat struct {
	MediaID string `json:"id"`
}

// generateThumbnails is our main request handler.
func generateThumbnails(ctx context.Context, in io.Reader, out io.Writer) {
	// Authenticate to obtain a usable client with our storage bucket.
	provider, err := auth.ResourcePrincipalConfigurationProvider()
	DieIfErr(err)
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	DieIfErr(err)

	// Generate a media ID.
	mediaCtx := context.WithValue(ctx, "MediaId", RandStringBytesMaskImprSrc(16))

	// Get the current storage namespace.
	namespace, err := client.GetNamespace(ctx, objectstorage.GetNamespaceRequest{})
	DieIfErr(err)
	mediaCtx = context.WithValue(mediaCtx, "NamespaceName", *namespace.Value)

	// Now, attempt to read our supposed image.
	// We only want to read up to 10 megabytes - our imposed maximum.
	// As a bonus, we read directly to the image - format issues will be instantly exposed.
	img := decodeImage(in)

	// Generate!
	generateVariants(mediaCtx, client, img)

	// We're all good!
	response := ResponseFormat{
		MediaID: mediaId(mediaCtx),
	}
	json.NewEncoder(out).Encode(response)
}
