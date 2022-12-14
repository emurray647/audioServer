# Audio API

## Table of Contents

- [Description](#desription)
- [Requirements](#requirements)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Design](#design)
- [Future Improvements](#future-improvements)

## Desription

This is an API server to store audio projects. Files can be uploaded with the `/files` endpoint,
and can be downloaded with the `/download` endpoint.  Furthermore, one can search for audio
files on the server with the `/list` endpoint.

Currently, this server only supports storing wav or mp3 files,

## Requirements

To run this program, you need the following installed on your machine

- `docker`
- `docker-compose`
- `make` (optional)

This project uses some go third party modules to read in the wav and mp3 files, including

- `github.com/go-audio/wav`
- `github.com/tcolgate/mp3`

For the full list of third party packages, see [go.mod](https://github.com/emurray647/audioServer/blob/main/src/go.mod).

## Usage

To build this application, simply run `make` or `docker-compose build`.  Then to run,
use the command `make start` or `docker-compose up`.  This will start up a server at
`localhost:8080` that you can send requests to.  For the format of the requests, see 
[Endpoints](#endpoints).

For running tests, the command `make test` will start up a test API docker container as well
as a test client docker container that will send some requests to the server and assert
on the responses.

## Endpoints

QuickLinks:
- [`POST /files`](#post-files)
- [`DELETE /files`](#delete-files)
- [`GET /download`](#get-download)
- [`GET /info`](#get-info)
- [`GET /list`](#get-list)

Unless otherwise specified, all endpoints will return a StatusMessage JSON object of the form:

```
{
    status_code, // http response code
    success, // boolean to indicate if the request was successful
    message, // optional message to provide an error or other information
}
```

### `POST /files`

Uploads a file to the server.  The file should be passed as part of the request payload (ie,
`curl -X POST --data-binary @<filename> http://localhost:8080/files).  Because of this, an
optional `name` value can be passed as a query param as a name for the file.  Otherwise
onen will be generated.

**Parameters**

| Name | Description |
| --- | --- |
| name  | optional name to give the file

**Response**

A StatusMessage response.

| StatusCode | Reasons |
| --- | --- |
| 200 | Success |
| 400 | Attempted to upload unknown file type, or file was not of the specified type |
| 409 | File with specified name already exists on the server |
| 500 | Server had an error processing the request |


### `DELETE /files`

Deletes a file from the server

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
| name | the name of the file |
| format | format of the file ("wav", "mp3",...)|
| file_size | the size of the file in bytes |
| minfile_size | |
| maxfile_size | |
| duration | the length of the audio in seconds |
| minduration | |
| maxduration | |
| num_channels | the number of channels in the file |
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

Downloads a file with the provided name

**Parameters**

| Name | Description |
| --- | --- |
| name | The name of the file to download |

**Response**

On success, will download the specified file.  On failure, will return a `StatusMessage` object

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

# Design

The layout of the project is as follows:

```
├── docker-compose.test.yml
├── docker-compose.yml
├── init.sql
├── Makefile
├── README.md
├── src
│   ├── Dockerfile
│   ├── ...
└── test
    ├── Dockerfile
    └── samples
        ├── ...
```

A brief description of the files:

- `docker-compose.test.yml`: docker configuration for setting up tests
- `docker-compose.yml`: docker configuration to run the audioapi server and its backing DB
- `init.sql`: the database schema
- `Makefile`: defines `make` rules for starting and stopping the server, as well as running tests
- `README.md`: description of the project
- `src`: directory containing all the go source code for the project
- `test`: directory containing the test setup Dockerfile, as well as sample files with which to test

A MySQL database is 

Since databases aren't optimized to contain large file, the decision was made to write the file
to disk, and then store the metadata for the file in MySQL, as well as the URI of where the file
was written.  Then the `/list` and `/info` endpoints just have to query the DB to get the appropriate
metdata, while `/download` hits the DB to get the URI, and then returns a copy of the file.

Name generation


## Future Improvements

### Customizable port 

This server is hardcoded to listen for traffic on 8080. It would be nice if this were customizable

### Unit tests 

As of now the only testing is more of a system test that spins up a test-backend server, 
hits it with some requests, and verifies the output.  The downsides are that this process
is slow and makes it difficult to iterate on while developing.

The better alternative (addition) would be for unit tests.  In this case we would want to 
mock calls to the `DBConnection` and then verify that our functions return what we would 
expect.  Unfortunately, I ran out of time to add these tests.

### Form Data

As of now we can only upload files via `POST /files` with the data in the request payload.
This means that we are not passing in the filename, and have to either manually pass it 
as a query parameter, or generate a name on the server.  If we were to have the capability
to pass in the file as form data, we could use the filename from there.
For ex: `curl -F name=<filename> -F upload=@<local_file> http://localhost:8080/files`.
This would have the same need to add the name twice, but it would mesh better 
for any frontend developers trying to use this API.


### Store Files somewhere other than to disk

Currently this API stores all uploaded files into the `/data` directory.  This has a huge downside
of preventing scalability.  If this API were to be hosted in two nodes rather than one, files 
would be written to disk on the node in which the original request was made.  Then if a 
download request were to be made to the other node, it would not be able to find the file.
Furthermore, if one of the nodes were to go down, the files on that node would be lost
forever.

The best alternative here is to use some sort of cloud storage, such as AWS S3.  If we were
to instead write files to S3 and store the S3 bucket/key as the URI, we could then retrieve
the file from any of our nodes.

### Additional Metadata

The metadata that is pulled from each file is currently minimal.  There is additional data,
such as artist, which would be neat to be able to search for.