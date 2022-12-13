


- src
    -internal
        - model
        - db
    main.go


## endpoints

### `POST /files`

Upload a file to the server

**Parameters**

| Name | Description |
| --- | --- |
| name  | optional name to give the file. if none is provided, a hash of the file will be used

**Response**

A StatusMessage response..

| StatusCode | Reasons |
| --- | --- |
| 200 | Success |
| 400 | Attempted to upload unknown file type, or file was not of the specified type |
| 409 | File with specified name already exists on the server |
| 500 | Server had an error processing the request |


### `DELETE /files`

Delete a file from the server

**Parameters**

| Name | Description |
| --- | --- |
| name  | the file to delete |

A StatusMessage response.

| StatusCode | Reasons |
| --- | --- |
| 200 | Success |
| 400 | No file name provided or file does not exist on the server |
| 500 | Server had an error processing the request |

### `GET /list`

List all files that match provided conditions

**Parameters**

| Name | Description |
| --- | --- |
| name | |
| minname | |
| maxname | |
| file_size | |
| minfile_size | |
| maxfile_size | |
| duration | |
| minduration | |
| maxduration | |
| minnum_channels | |
| maxnum_channels | |
| sample_rate | |
| minsample_rate | |
| maxsample_rate | |
| audio_format | |
| minaudio_format | |
| maxaudio_format | |
| avg_bytes_per_second | |
| minavg_bytes_per_second | |
| maxavg_bytes_per_second | |

**Response**

On success, will return a JSON object such as
```
[
  {
    "name": "gettysburg.wav",
    "file_size": 441180,
    "duration": 10.0039,
    "num_channels": 1,
    "sample_rate": 22050,
    "audio_format": 1,
    "avg_bytes_per_second": 44100
  },
  {
    "name": "cantina.wav",
    "file_size": 2646044,
    "duration": 60.0008,
    "num_channels": 1,
    "sample_rate": 22050,
    "audio_format": 1,
    "avg_bytes_per_second": 44100
  }
]
```
### GET /download

Downloads a file with the specified name

**Parameters**

| Name | Description |
| --- | --- |
| name | The name of the file to download |

**Response**

On success, will downlad the specified file.  On failure, will return a `StatusMessage` object

| StatusCode | Reasons |
| --- | --- |
| 400 | No file name provided |
| 404 | File is not on server |
| 500 | Server had an error processing the request |


### GET /info

Returns information about a specific file

**Parameters**

| Name | Description |
| --- | --- |
| name | The name of the file on which to get information |

**Response**

On success will return a message of the form

```
{
  "name": "cantina.wav",
  "file_size": 2646044,
  "duration": 60.0008,
  "num_channels": 1,
  "sample_rate": 22050,
  "audio_format": 1,
  "avg_bytes_per_second": 44100
}
```

On failure, will return a `StatusMessage` with the following possible errors

| StatusCode | Reasons |
| --- | --- |
| 400 | No file name provided |
| 404 | File does not exist on server |
| 500 | Server had an error processing the request |

### files

invalid
duplicate
unknown error / server error

| asdf | asdf |
| --- | --- |
| 200 | File successfully uploaded |
| 400 | asdf |

### download

### list

### info


for `file` endpoint, store the file to the file system, and put its location
into the db

we are going to want other fields in the db as well, such as length
(anything we will want for `list`)

```
id, name, length_minutes, ..., filesystem_uri
```

`curl -X POST --data-binary @/home/emurray/Downloads/CantinaBand60.wav http://localhost:8080/files/`


maxduration=300 -> max,duration,300 -> duration <= 300

TODO
- index db
- have successful delete return something
- query for format
