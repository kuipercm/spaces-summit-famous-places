# Spaces Summit Famous Places

## How-to

1. Use a terminal to authenticate gcloud: `gcloud auth login`
2. Manually create a project in CloudPlayground. See https://docs.bol.io/cloud-setup/way-of-working/using_playground.html#_using_gcloud
3. Set default project: `gcloud config set project <projectid>` 
4. Manually create a service account and assign cloud build roles
```
gcloud iam service-accounts create terraform-sa \
    --description="Terraform Automation" \
    --display-name="Terraform Automation"
```
```
gcloud projects add-iam-policy-binding <projectid> \
    --member="serviceAccount:terraform-sa@<projectid>.iam.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.editor" 
```
```
gcloud projects add-iam-policy-binding <projectid> \
    --member="serviceAccount:terraform-sa@<projectid>.iam.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.builder"
```
```
gcloud projects add-iam-policy-binding <projectid> \
    --member="serviceAccount:terraform-sa@<projectid>.iam.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.viewer"
```
5. Create a service account key and download the json
```
gcloud iam service-accounts keys create ./spacessummit-terraform-sa-private-key.json \
    --iam-account=terraform-sa@<projectid>.iam.gserviceaccount.com
```

IMPORTANT: the next steps (6-8) can also be done in the Terraform console, which is a lot easier.
Also, if you don't run Terraform a lot, you can run Terraform by hand.

6. Setup new workspace in Terraform IO - to automatically configure the cloud build. Perform
```
curl \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @payload.json \
  https://app.terraform.io/api/v2/organizations/3ldn/workspaces
```
where `$TOKEN` is the personal identification token for Terraform, `3ldn` is the organization name
and the `payload.json` is
```json
{
  "data": {
    "attributes": {
      "name": "spaces-summit-famous-places-<suffix>",
      "resource-count": 0,
      "working-directory": "",
      "vcs-repo": {
        "identifier": "kuipercm/spaces-summit-famous-places",
        "oauth-token-id": "<auth-token>",
        "branch": "master"
      },
      "auto-apply": true,
      "updated-at": "<current-date>"
    },
    "type": "workspaces"
  }
}
```
where the `auth-token` is retrieved by first retrieving the oauth-client-id 
```
curl --header "Authorization: Bearer $TOKEN" \
   https://app.terraform.io/api/v2/organizations/3ldn/oauth-clients
```
and subsequently retrieving the oauth-token by
```
curl --header "Authorization: Bearer $TOKEN" \
   https://app.terraform.io/api/v2/oauth-clients/<oauth-client-id>/oauth-tokens
```
7. Set appropriate Terraform workspace variables
```
curl \
  --header "Authorization: Bearer $TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request POST \
  --data @payload.json \
  https://app.terraform.io/api/v2/vars
```
where `payload.json` contains
```
{
  "data": {
    "type":"vars",
    "attributes": {
      "key":"google_credentials",
      "value":"<json content of service account token>",
      "description":"",
      "category":"terraform",
      "hcl":false,
      "sensitive":true
    },
    "relationships": {
      "workspace": {
        "data": {
          "id":"<workspace-id>",
          "type":"workspaces"
        }
      }
    }
  }
}
```
where the json token content is from the previously created `./spacessummit-terraform-sa-private-key.json` file.

At this point, running the terraform files should create the correct build trigger in GCP, which should in turn
create the appropriate deployment to Cloud Run to deploy the app.

The app itself creates a bucket and a topic (currently) to store uploaded files.

To test the app, go into the GCP console -> cloud run and find out the app url after deployment is successful.

To test the app, perform the following command

```
curl --location --request PUT 'https://famous-places-3guutasc6a-ez.a.run.app/api/upload' \
--header 'Content-Type: multipart/form-data' \
--header "Authorization: Bearer $(gcloud auth print-identity-token)" \
 --form 'photo=@<some image>'
```

The file should appear in the bucket with its name replaced by a UUID.