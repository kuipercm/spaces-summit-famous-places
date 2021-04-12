# Spaces Summit Famous Places

## How-to

1. Use a terminal to authenticate gcloud: `gcloud auth login`
2. Manually create a project in CloudPlayground. See https://docs.bol.io/cloud-setup/way-of-working/using_playground.html#_using_gcloud
3. Set default project: `gcloud config set project <projectid>` 
4. Create the build trigger
```
gcloud beta builds triggers create github \
--repo-name=spaces-summit-famous-places \
--repo-owner=kuipercm \
--branch-pattern=^master$ \
--build-config=cloudbuild.yaml \
```
5. At this point, there should be a cloud build trigger created which should run and deploy
the app to Cloud Run.
6. The app itself creates a bucket and a topic (currently) to store uploaded files.
7. To test the app, go into the GCP console -> cloud run and find out the app url after deployment is successful. 
To test the app, perform the following command
```
curl --location --request PUT 'https://famous-places-3guutasc6a-ez.a.run.app/api/upload' \
--header 'Content-Type: multipart/form-data' \
--header "Authorization: Bearer $(gcloud auth print-identity-token)" \
 --form 'photo=@<some image>'
```
The file should appear in the bucket with its name replaced by a UUID.