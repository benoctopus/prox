# prox

Simple reverse proxy written in go for use with docker. Work in progress.

## How to use
Have these installed on your system:
- Docker
- docker-compose
- a bash shell

In dev mode everything runs in docker WITH hot reloading on unix based systems out of the box. Just run `start_dev.sh` and you're golden. On Windows, the hot reloading may require additional software due to the way windows mounts th... blah blah, try using [this](https://github.com/merofeev/docker-windows-volume-watcher) script. Other than that it works fine with wsl and docker for windows.

There is currently no production mode, but i'll figure that out soon since I would like to use it with my own projects. The proxy is configured through a json file `config.json` which in dev mode lives in the src folder. If it is not present, the server will automatically use the config.dev.json instead, which just proxies to some other docker containers containing benign applications for testing purposes (all started in the docker-compose.yml). For an example of the config structure check out the dev config.

## Things I would like to do

- make this readme more useful
- global application session management with redis!
- csrf protection baked in (it's kindof there but doesn't work that well)
- static file serving (also kindof there)
- overlapping routes (ie template served at / but also static file folder at /, go makes this difficult :'( )
- gzip?
- scalable?

## Thanks bro

[gorilla](https://github.com/gorilla)
