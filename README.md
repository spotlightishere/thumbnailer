# thumbnailer
An [Oracle Cloud Function](https://www.oracle.com/cloud-native/functions/) designed to render different image sizes for a given image.
This is excellent in the case of thumbnailing for a media gallery, or avatar images.

## Setup
Setting up this service is rather nuanced.

You will need to perform a few steps as a prerequisite:
 - Create a bucket for this function. You can grant it whatever name you wish.
 - Create an application for this function as described in ["Get Started using the CLI"](https://docs.oracle.com/en-us/iaas/developer-tutorials/tutorials/functions/func-setup-cli/01-summary.htm).
 - Using the same article above, deploy the function to your application. The function name `thumbnailer` is recommended for consistency with other services.
 - Within the application holding this function, configure the key `BUCKET_NAME` to be the bucket you wish to store thumbnail output.
The namespace will be discovered automatically via the API during invocation.

You will then need to grant the thumbnailing function access to the bucket:
 - Within your application, select the function. Take note of its OCID.
 - Go to "Dynamic Groups" within the Identity section.
 - Create a new group.
   - You may use any name and description.
   - For rules, input the following, substituting the default value for your own:
```
resource.id = 'ocid1.fnfunc1.oc1.iad.xxxxxxxxx'
```
 - Within the Identity section, go to "Policies" and create a new policy.
   - You may use any name and description.
   - Within the Policy Builder, toggle "Show manual editor" and input the following, substituting values for your own configuration:
```
allow dynamic-group <DYNAMIC_GROUP> to read buckets in tenancy where target.bucket.name='<BUCKET_NAME>'
allow dynamic-group <DYNAMIC_GROUP> to manage objects in tenancy where any {request.permission='OBJECT_CREATE', request.permission='OBJECT_INSPECT',target.bucket.name='<BUCKET_NAME>'}
```