### Tracking PI

#### Development
Backend located in `/`. Run
```shell
$ cp .env.example .env
$ go get .
$ go run .
```

Frontend located in `views/`. Run 
```shell
$ pnpm add 
```

```shell
$ pnpm build --watch
```

During production set `DEVELOPMENT=false` or omit it from env.
