# This file is used as a base for .goreleaser for easy releasing.
# The image is intended to be used as a standalone "product"
# Local and more specialized images can be found in the `Docker/` directory of this repo.
FROM autodesk-docker-build-images.***REMOVED***/hardened-build/golang-1.16:latest as BASE
LABEL maintainer="shore@autodesk.com"

CMD ["shore"]
COPY shore /usr/local/bin/shore
