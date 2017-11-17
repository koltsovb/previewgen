# previewgen
Previewgen is a http-service intended for creating image 100x100 preview.

### Dependances
```
github.com/nfnt/resize
```

### Build
```
go build
```
### Test
```
go test -v
```
### Run
```
previewgen -p=8080 -d=/tmp/pg 
# -p - service port
# -d - service working directory (for log file and image' thumb)
# to display log - tail -f /tmp/pg/previewgen.log
```


## API

### 1. Post image in json
### Request:
```
POST /files HTTP/1.1
Host: 127.0.0.1
Content-Type: application/json
 
{
    "type" : "files",
    "files" : [
      { 
          "name" : "f1",
          "mimeType" : "image/png",
          "content" : "... image in base64 ..."
      },
      { 
          "name" : "f2",
          "mimeType" : "image/jpeg",
          "content" : "... image in base64 ..."
      }
    ]
}
```
### Responce:
```
HTTP/1.1 200 OK
{
  "f1":6,
  "f2":7
}
```

### 2. Post image-url in json
### Request:
```
POST /files HTTP/1.1
Host: 127.0.0.1
Content-Type: application/json

{ 
	"Type" : "urls",
	"Files" : [ 
	       {"URL" : "http://redis.io/images/redis-white.png"},
	       { "URL" : "https://memcached.org/images/memcached_banner75.jpg" }
	]
}
```
### Responce:
```
HTTP/1.1 200 OK
{
 "http://redis.io/images/redis-white.png":4,
 "https://memcached.org/images/memcached_banner75.jpg":5
}
```
### 3. Post image in form-data
### Request:
```
POST /files HTTP/1.1
Host: 127.0.0.1
Content-Type: multipart/form-data; boundary=------------------------dfa1a27f9d9cca0c
Content-Disposition: form-data; name="image_field"; filename="1.jpg"
Content-Type: image/jpeg
# file data
boundary=------------------------dfa1a27f9d9cca0c
```
### Responce:
```
HTTP/1.1 200 OK
{
 "1.jpg":9
}
```
### curl example
```
curl -v -X POST -F "image_field=@./test/1.jpg" '127.0.0.1/files'
```



