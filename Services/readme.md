# Services

## Docker guide

### Categories service

Must be built from the path `<libery dungeon src>/Services/`. Use the command:

```bash
    docker buildx build -t defalt115/libery-dungeon-categories:<new version tag> -f ./Categories/Dockerfile --load .
```


### Medias service

Must be built from its own service directory e.g: `<libery dungeon src>/Services/Medias/`. Use the command:

```bash
    docker buildx build --progress=plain -t defalt115/libery-dungeon-medias:<new version tag> -f ./Dockerfile --load .. 
```

**Important**: When building the medias service for a arm64 architecture, you must provide the precomplied version of the `ffmpeg` binary.


### Metadata service

Must be built from the path `<libery dungeon src>/Services/`. Use the command:

```bash
    docker buildx build -t defalt115/libery-dungeon-metadata:<new version tag> -f ./Metadata/Dockerfile --load .
```

### Collect service

Must be built from the path `<libery dungeon src>/Services/`. Use the command:

```bash
    docker buildx build -t defalt115/libery-dungeon-collect:<new version tag> -f ./Collect/Dockerfile --load .
```

### Downloads service

Must be built from the path `<libery dungeon src>/Services/`. Use the command:

```bash
    docker buildx build -t defalt115/libery-dungeon-downloads:<new version tag> -f ./Downloads/Dockerfile --load .
```