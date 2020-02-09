# design-upload-roboter

Projekt-ID: design-upload-roboter
Projjekt-Number: 553986649280

gcloud config set project design-upload-roboter
gcloud config set run/region europe-west1

gsutil mb gs://design-upload-roboter-import
gsutil mb gs://design-upload-roboter-output

gcloud pubsub topics create design-upload-roboter-import-topic

gcloud builds submit --tag gcr.io/design-upload-roboter/convertpsd2png


gcloud run deploy design-upload-roboter --image gcr.io/design-upload-roboter/convertpsd2png --set-env-vars=CONVERTED_BUCKET_NAME=design-upload-roboter-output --platform managed

https://design-upload-roboter-fn5hvfqesa-uc.a.run.app

gcloud projects add-iam-policy-binding design-upload-roboter --member=serviceAccount:service-553986649280@gcp-sa-pubsub.iam.gserviceaccount.com --role=roles/iam.serviceAccountTokenCreator

gcloud iam service-accounts create cloud-run-pubsub-invoker --display-name "Cloud Run Pub/Sub Invoker"

gcloud beta run services add-iam-policy-binding design-upload-roboter --member=serviceAccount:cloud-run-pubsub-invoker@design-upload-roboter.iam.gserviceaccount.com --role=roles/run.invoker

gcloud beta pubsub subscriptions create convertSubscription --topic design-upload-roboter-import-topic --push-endpoint=https://design-upload-roboter-fn5hvfqesa-ew.a.run.app --push-auth-service-account=cloud-run-pubsub-invoker@design-upload-roboter.iam.gserviceaccount.com

gsutil notification create -t design-upload-roboter-import-topic -e OBJECT_FINALIZE -f json gs://design-upload-roboter-import

gsutil cp sampledocs/wonder_to_flight.psd gs://design-upload-roboter-import
