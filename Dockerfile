# This file is used as a base for .goreleaser for easy releasing.
# The image is intended to be used as a standalone "product"
FROM ***REMOVED***.dev.adskengineer.net/container-hardening/alpine-hardened-min:latest
LABEL maintainer="shore@autodesk.com"

CMD ["shore"]
COPY shore /usr/local/bin/shore
