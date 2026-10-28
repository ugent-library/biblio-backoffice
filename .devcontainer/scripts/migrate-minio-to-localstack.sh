#!/bin/bash
# Copies the files in the old MinIO volume of the devcontainer into LocalStack.
#
# Run it once on the host (not inside the devcontainer), with the devcontainer
# stack running. It starts a temporary MinIO on the old volume, copies the
# biblio bucket with rclone and removes the temporary MinIO again. The old
# volume is left as is; remove it yourself once you're happy:
#
#   docker volume rm biblio-backoffice_devcontainer_miniodata
#
# MinIO images are no longer published, so this only works when the image is
# still cached locally. Without it, the script skips the migration.
set -euo pipefail

PROJECT="${PROJECT:-biblio-backoffice_devcontainer}"
MINIO_VOLUME="${MINIO_VOLUME:-${PROJECT}_miniodata}"
NETWORK="${PROJECT}_default"
BUCKET="biblio"

MINIO_RELEASE="RELEASE.2023-12-09T18-17-51Z"
MINIO_IMAGE="minio/minio:${MINIO_RELEASE}"
MINIO_CONTAINER="${PROJECT}-minio-migrate"

RCLONE_IMAGE="rclone/rclone:1.75.2"

if ! docker volume inspect "$MINIO_VOLUME" >/dev/null 2>&1; then
	echo "No MinIO volume $MINIO_VOLUME, nothing to migrate."
	exit 0
fi

# MinIO keeps every object as a directory in the bucket directory.
first_object=$(docker run --rm -v "$MINIO_VOLUME:/data:ro" alpine:3.24 sh -c "ls -A /data/$BUCKET 2>/dev/null | head -1")
if [[ -z "$first_object" ]]; then
	echo "The $BUCKET bucket in $MINIO_VOLUME is empty, nothing to migrate."
	exit 0
fi

if [[ -z "$(docker ps -q --filter "label=com.docker.compose.project=$PROJECT" --filter "label=com.docker.compose.service=localstack")" ]]; then
	echo "LocalStack isn't running in $PROJECT. Start the devcontainer first." >&2
	exit 1
fi

if ! docker image inspect "$MINIO_IMAGE" >/dev/null 2>&1; then
	echo "$MINIO_IMAGE isn't cached locally, skipping the migration."
	exit 0
fi

cleanup() {
	docker rm -f "$MINIO_CONTAINER" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "Starting a temporary MinIO on $MINIO_VOLUME..."
docker run -d --name "$MINIO_CONTAINER" --network "$NETWORK" \
	-v "$MINIO_VOLUME:/data" \
	-e MINIO_ROOT_USER=minioadmin \
	-e MINIO_ROOT_PASSWORD=minioadmin \
	--entrypoint /usr/bin/minio \
	"$MINIO_IMAGE" server /data --address :9000 >/dev/null

rclone() {
	docker run --rm --network "$NETWORK" \
		-e RCLONE_CONFIG_MINIO_TYPE=s3 \
		-e RCLONE_CONFIG_MINIO_PROVIDER=Minio \
		-e RCLONE_CONFIG_MINIO_ENDPOINT="http://$MINIO_CONTAINER:9000" \
		-e RCLONE_CONFIG_MINIO_ACCESS_KEY_ID=minioadmin \
		-e RCLONE_CONFIG_MINIO_SECRET_ACCESS_KEY=minioadmin \
		-e RCLONE_CONFIG_LOCALSTACK_TYPE=s3 \
		-e RCLONE_CONFIG_LOCALSTACK_PROVIDER=Other \
		-e RCLONE_CONFIG_LOCALSTACK_ENDPOINT=http://localstack:3071 \
		-e RCLONE_CONFIG_LOCALSTACK_ACCESS_KEY_ID=test \
		-e RCLONE_CONFIG_LOCALSTACK_SECRET_ACCESS_KEY=test \
		-e RCLONE_CONFIG_LOCALSTACK_REGION=us-east-1 \
		"$RCLONE_IMAGE" "$@"
}

for i in $(seq 1 30); do
	rclone lsf --max-depth 1 "minio:$BUCKET" >/dev/null 2>&1 && break
	if [[ $i -eq 30 ]]; then
		echo "MinIO didn't come up. Its logs:" >&2
		docker logs "$MINIO_CONTAINER" >&2
		exit 1
	fi
	sleep 2
done

echo "Copying minio:$BUCKET to localstack:$BUCKET..."
rclone copy --stats-one-line --stats 10s -v "minio:$BUCKET" "localstack:$BUCKET"
rclone check --one-way "minio:$BUCKET" "localstack:$BUCKET"
echo "Done. The files are in LocalStack; $MINIO_VOLUME is untouched."
