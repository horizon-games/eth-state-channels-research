#!/bin/bash
set -e # exit on any error
set -u # error out on unbound variables (empty $WORKDIR etc.)

# Check if the image is available in Docker registry first.
sudo docker pull arcadeum/server:$VERSION && { echo "Image already built."; exit 0; } || :

# Create temporary directory for build, clean it up on exit.
DIR=$(mktemp -d) && trap "sudo rm -rf $DIR" EXIT && cd ${DIR}
git clone --depth 1 --single-branch -b ${GITTAG:-$GITBRANCH} git@github.com:horizon-games/arcadeum.git ./
wait $(jobs -p)

cd $DIR/server

# Build Docker image. Tag it with version and re-tag latest.
sudo docker build -t arcadeum/server:$VERSION .

# Push to Docker registry.
sudo docker push arcadeum/server:$VERSION

echo
echo "arcadeum/server:$VERSION built successfully"
echo
echo "Now, you can deploy it with:"
echo "sup -e VERSION=$VERSION $SUP_NETWORK deploy"
echo
