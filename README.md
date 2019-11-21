## To make a deploy with Docker 

- heroku container:push web -a dev01-go-base
- heroku container:release web -a dev01-go-base


## Local testing

The env file allows to use environments variables running heroku locally.
When you launch the command heroku local, the Procfile is readed and runned the command contained inside. 
For this example app, the procfile contains the command `web: bin/go-base`, this means that before you launch the heroku local command,
you need to build your application. For this specific case you only need to run the following command from the root of this project:
```
$ mkdir bin
$ go build -o bin/go-base
```

As a result, you have your fantastic app, builded into the bin folder.
Now you can launch the heroku local.

Remember that, if your `.env` file, contains the following line:
```
MY_LOCAL_VARIABLE=my-local-value
```
To access this variable, in your code you need to call:
```
os.GetEnv("MY_LOCAL_VARIABLE")
```

## TODO

- Configure docker compose in order to use also DB an redis