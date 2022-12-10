


- src
    -internal
        - model
        - db
    main.go



for `file` endpoint, store the file to the file system, and put its location
into the db

we are going to want other fields in the db as well, such as length
(anything we will want for `list`)

```
id, name, length_minutes, ..., filesystem_uri
```

`curl -X POST --data-binary @/home/emurray/Downloads/CantinaBand60.wav http://localhost:8080/files/`