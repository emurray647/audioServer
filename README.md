


- src
    -internal
        - model
        - db
    main.go


## endpoints

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
- add mp3