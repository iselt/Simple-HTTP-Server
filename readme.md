# Simple HTTP Server

Designed for environments with limited operations to transfer files, supports POST/PUT to upload files and GET for static files.

## Usage

### Start Server

```shell
simple_http_server <HOST> <PORT> [<ROOT_DIR>]
```

### Download File

```shell
curl http://<HOST>:<PORT>/<FILE_PATH>
```

### Upload File

```shell
curl -X POST http://<HOST>:<PORT>/<FILE_PATH> --data-binary @<LOCAL_FILE_PATH>
```
