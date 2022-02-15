# thumbnailer
An [Oracle Cloud Function](https://www.oracle.com/cloud-native/functions/) designed to render different image sizes for a given image. This is excellent in the case of thumbnailing for a media gallery, or avatar images.

## Setup
Within the application holding this function, configure the key `BUCKET_NAME` to be the bucket you wish to store thumbnail output.