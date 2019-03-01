# Dependency Extraction Service - DES

Dependency Extraction Service (DES) is a simple gRPC service that encapsulates the logic for extracting library dependency information.
It does this by parsing well known dependency management files (pom.xml, build.gradle, go.mod, package.json to name a few).
After parsing out the information, it returns a standard representation making it easy to store and query.

## Support

## Getting Started

```
docker build . -t mjpitz/des
docker run --rm mjpitz
```
