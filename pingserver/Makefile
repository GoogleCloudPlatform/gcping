project = computelabs
image = pingserver
zones = asia-east1-a \
		asia-east2-a \
		asia-northeast1-a \
		asia-northeast2-a \
		asia-south1-a \
		asia-southeast1-a \
		australia-southeast1-a \
		europe-north1-a \
		europe-west1-b \
		europe-west2-a \
		europe-west3-a \
		europe-west4-a \
		europe-west6-a \
		northamerica-northeast1-a \
		southamerica-east1-a \
		us-central1-a \
		us-east1-b \
		us-east4-a \
		us-west1-a \
		us-west2-a

build:
	docker build -t gcr.io/$(project)/$(image) .

push:
	docker push gcr.io/$(project)/$(image)

run:
	docker run -p 80:80 gcr.io/$(project)/$(image)

deploy-create:
	for zone in $(zones) ; do \
		gcloud compute instances \
    		create-with-container pingserver-$$zone \
			--zone $$zone \
			--container-image gcr.io/$(project)/$(image); \
	done

deploy-update:
	for zone in $(zones) ; do \
		gcloud compute instances \
    		update-container pingserver-$$zone \
			--zone $$zone \
			--container-image gcr.io/$(project)/$(image); \
	done


firewall:
	gcloud compute firewall-rules \
      create allow-80 --allow tcp:80
