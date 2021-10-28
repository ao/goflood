# GoFlood

A simple CLI tool that lets you make GET requests to a domain, and ignores the output

## Where to download GoFlood

Packaged binaries for all platforms are available on the [GoFlood Github Release page](https://github.com/ao/goflood/releases/)


## Usage once Downloaded/Installed

### Option 1: Follow the wizard

You can simply run the CLI application and follow the prompts

```
./goflood
```

### Option 2: Specify the arguments

```
./goflood example.com   10       1
          ^             ^        ^
          domain        count    batch
```

#### Commandline Arguments

| argument | description | example value |
|----------|-------------|-------|
| domain   | a domain name | example.com
| count    | amount of concurrent reqs | 50 |
| batch    | amount of times to repeat | 3 |
