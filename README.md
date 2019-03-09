# Dependency Extraction Service - DES

Dependency Extraction Service (DES) is a simple gRPC service that encapsulates the logic for extracting library dependency information.
It does this by parsing well known dependency management files (pom.xml, build.gradle, go.mod, package.json to name a few).
After parsing out the information, it returns a standard representation making it easy to store and query.

## Support

![GitHub](https://img.shields.io/github/license/mjpitz/rds.svg)
[![Build Status](https://travis-ci.com/mjpitz/des.svg?branch=master)](https://travis-ci.com/mjpitz/des)
[![](https://images.microbadger.com/badges/image/mjpitz/des.svg)](https://microbadger.com/images/mjpitz/des)
[![](https://images.microbadger.com/badges/version/mjpitz/des.svg)](https://microbadger.com/images/mjpitz/des)

## Getting Started

A docker image is regularly built and uploaded to docker.io.

```bash
docker run --rm mjpitz/des
```

If you make changes and want to test them out, you can run via a locally built docker image or using `npm start`.

```bash
docker build . -t mjpitz/des
```
