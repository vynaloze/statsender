# statsender

[![Documentation](https://godoc.org/github.com/vynaloze/statsender?status.svg)](https://godoc.org/github.com/vynaloze/statsender)
[![Go Report Card](https://goreportcard.com/badge/github.com/vynaloze/statsender)](https://goreportcard.com/report/github.com/vynaloze/statsender)
[![Build Status](https://travis-ci.com/vynaloze/statsender.svg?branch=master)](https://travis-ci.com/vynaloze/statsender)

*statsender* collects and sends out various PostgreSQL server statistics.
It was designed to work together with [*pgmeter*](https://github.com/vynaloze/pgmeter) to provide end-to-end monitoring of distributed PostgreSQL databases.

## Install

1. [Download the latest release](https://github.com/vynaloze/statsender/releases/latest) for your system.
2. Unzip the downloaded zip archive, containing a single binary called `statsender` (Linux) or `statsender.exe` (Windows).
3. (Optional) add *statsender* binary to PATH. [(Linux instructions)](https://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux-unix) [(Windows instructions)](https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows)
4. Verify the installation by executing `statsender` in a terminal. You should see output like this:
```
Usage:                                                                        
  statsender [command]                                                        
                                                                              
Available Commands:                                                           
  collector   Manage collectors                                               
  datasource  Manage datasources                                              
  help        Help about any command                                          
  init        Initializes the application config                              
  run         Runs the application in detached mode                           
  sender      Manage senders                                                  
  try         Tests the application                                           
                                                                              
Flags:                                                                        
  -c, --config string   sets configuration directory location (default "conf")
  -h, --help            help for statsender                                   
  -v, --verbose         verbose output                                        
      --version         version for statsender                                
                                                                              
Use "statsender [command] --help" for more information about a command.
```

## Try it out!

This is a simplest use case: you have a local PostgreSQL server running on default port,
with [enabled pg_stat_statements](https://github.com/vynaloze/statsender/wiki/Enable-pg_stat_statements), 
and [*pgmeter*](https://github.com/vynaloze/pgmeter) is running somewhere, with default configuration.

1. Run `statsender init` to generate first, default configuration files and examples. By default, they'll appear in a new directory `./conf`
2. Define your DB connection: `statsender datasource add <login>:<password>@localhost/<database>?sslmode=disable`
   
   Remember to use superuser or [restricted monitoring user](https://github.com/vynaloze/statsender/wiki/Set-restricted-monitoring-user) credentials.
3. Define the web location of *pgmeter*: `statsender sender add http <pgmeter_endpoint>`
4. Test the connection: `statsender try`

    You should see something like:
    ```
    [INFO] reading configuration...                     
    [OK] configuration structure is valid               
                                                        
    [INFO] testing collectors...                        
    (...)                                        
    [OK] collector structure is valid                   
                                                        
    [INFO] testing datasources...                       
    (...)         
    [OK] datasource structure is valid                  
                                                        
    [INFO] testing senders...                           
    [OK] sender structure is valid                      
                                                        
    [INFO] Test complete! Looks like you are good to go!
    ```
5. Run: `statsender run` to start collecting statistics in the background. Default log directory is `./logs`
6. Congratulations! If you want to learn some more detailed stuff now, see the next section.

## Detailed configuration

*statsender* is based around three main concepts: datasources, collectors and senders:

- **datasource** defines *from* where to collect statistics
- **collector** defines *what* to collect from the datasource
- **sender** defines *where* to send the collected statistics 

In most cases, you'll be fine with the default **collector** setup - however, 
you'll need to define **datasources** and **senders** by yourself, either editing configuration files
by hand, or using the [CLI](https://github.com/vynaloze/statsender/wiki/CLI-reference).

*statsender* uses [HCL](https://github.com/hashicorp/hcl2/) (Hashicorp Configuration Language)
as its primary configuration language. However, the equivalent JSON structure will also be accepted.
To read the configuration, all the `.hcl` and `.json` files in the given configuration directory (default is `./conf`)
are parsed - so you can split your config into as many files as you want.

**Pro tip:** `cd conf` and see the configuration files generated by `statsender init`.
They provide many examples how to customize the *statsender*. Feel free to experiment with them. 
Running `statsender try` allows you to see the result without any consequences.

### Datasources

Datasources define from where to collect statistics. You can define more than one datasource,
for example if you run many databases on a single machine. If you really want, you can also define remote 
datasources, on different machines (but remember that you won't get system stats from them).

#### CLI

`statsender datasource add <DSN> [<tags>]`

- DSN: `[postgresql://]login:password@host[:port]/dbname[?param1=value1&...]`
- optional tags provided as flags `--tag key1=value1 --tag key2=value2 ...`

Examples:
- `statsender datasource add user:pass@localhost/testdb`
- `statsender datasource add postgresql://user:pass@localhost:6432/testdb`
- `statsender datasource add user:pass@10.0.3.13/proddb?sslmode=disable --tag env=prod --tag iteration=v1`

#### Example HCL block

```hcl
datasource {
  host = "10.0.3.13"
  port = 5432
  username = "user"
  password = "pass"
  dbname = "proddb"
  sslmode = "disable"
  tags = {
    env = "prod",
    iteration = "v1"
  }
}
```

### Senders

Senders define where to send the collected statistics. There is *console* sender 
(used for debugging) and *http* sender (to send statistics to *pgmeter*).
You can define any number of senders. 

#### CLI

`statsender sender add <type> [<spec>]`

- type: `console` or `http`
- spec:
  - if `console`, then none
  - if `http`, then full address of a target: `http[s]://host[:port][/endpoint]`

Examples:
- `statsender sender add console`
- `statsender sender add http 10.0.1.1:8080/stats`
- `statsender sender add http 10.0.1.1:8080/stats --retries 5 --delay 60`

#### Example HCL blocks

```hcl
console {}
```
```hcl
http {
  target = "http://10.0.1.1:8080/stats"
  retryDelay = 60
  maxRetries = 5
}
```

### Collectors

Collectors defines what to collect from the datasource.
There is probably no need to change anything in this section. 
However, if you want, see the [reference](https://github.com/vynaloze/statsender/wiki/Collector-reference) of all possible collectors, CLI below and examples in `./conf/_collectors.hcl`

#### CLI

`statsender collector enable <type>` enables a collector. 

`statsender collector disable <type>` disables a collector.

`statsender collector schedule <type> <cron>` changes the cron schedule of a collector.

- type: full name of a collector (see [reference](https://github.com/vynaloze/statsender/wiki/Collector-reference))
- cron: (quoted) string representing desired cron expression 

Examples:
- `statsender collector enable cpu`
- `statsender collector disable pg_stat_statements`
- `statsender collector schedule pg_stat_user_indexes @hourly`
- `statsender collector schedule pg_stat_user_indexes '0 0 * * * *'`

#### Example HCL blocks

```hcl
system {
  cpu {
    cron    = "*/5 * * * * *"
    enabled = true
  }
}
```
```hcl
postgres {
  pg_stat_statements {
    cron    = "@hourly"
    enabled = false
  }
  pg_stat_user_indexes {
    cron    = "0 0 * * * *"
    enabled = true
  }
}
```

## Automated provisioning

### Out-of-the-box solutions

At the moment, there are no out-of-the-box solutions like ansible plugins, puppet modules 
or docker images available. Fortunately, setup process can be easily scripted:

### Example bash script

```
curl -s https://github.com/vynaloze/statsender/releases/latest/download/statsender_linux_x86_64.zip
unzip statsender_linux_x86_64.zip
./statsender init
./statsender datasource add ...
./statsender sender add http ...
./statsender run
```

## Wiki pages

- [Enable pg_stat_statements](https://github.com/vynaloze/statsender/wiki/Enable-pg_stat_statements)
- [Set restricted monitoring user](https://github.com/vynaloze/statsender/wiki/Set-restricted-monitoring-user)
- [Full CLI reference](https://github.com/vynaloze/statsender/wiki/CLI-reference)
- [Collector reference](https://github.com/vynaloze/statsender/wiki/Collector-reference)

## License
The library is licensed under the [MIT License](LICENSE).