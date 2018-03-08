#!/bin/bash
set -e # exit on any error
set -u # error out on unbound variables (empty $WORKDIR etc.)

# TODO: rename davatar org to arcadeum once billing is setup there
# for now this is fine.

# Check if the image is available in Docker registry first.
sudo docker pull davatar/dgame-server:$VERSION && { echo "Image already built."; exit 0; } || :

# Create temporary directory for build, clean it up on exit.
DIR=$(mktemp -d) && trap "sudo rm -rf $DIR" EXIT && cd ${DIR}
git clone --depth 1 --single-branch -b ${GITTAG:-$GITBRANCH} git@github.com:horizon-games/dgame-server.git ./
wait $(jobs -p)

# Build Docker image. Tag it with version and re-tag latest.
sudo docker build -t davatar/dgame-server:$VERSION .

# Push to Docker registry.
sudo docker push davatar/dgame-server:$VERSION

echo
echo "arcadeum/dgame-server:$VERSION built successfully"
echo
echo "Now, you can deploy it with:"
echo "sup -e VERSION=$VERSION $SUP_NETWORK deploy"
echo
