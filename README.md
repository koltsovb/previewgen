# previewgen
Previewgen is a http-servise intended for creating image 100x100 preview.

## API
```
### 1. Post image in json
#### Request:
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

Responce:
HTTP/1.1 200 OK
{
  "f1":6,
  "f2":7
}

2. Post image-url in json

3. Post image in form-data
```

### Dependances
```
github.com/nfnt/resize
```

###
