// +build windows

package config

var defaultCollectorConfig = `/*
------------------------------------
statsender - collector configuration
------------------------------------

Example collector block structure:

name_of_collector {   // Name of a collector
  cron = "@hourly"    // Cron expression, which defines data update period
  enabled = true      // Switch to enable/disable a collector
}

There are two ways to disable the collector:
1. set 'enabled' to 'false'
2. comment or delete the whole block
Using the second approach, you won't be able to manage
this collector using CLI anymore, though.

*/

system {
  cpu {
    cron = "0 */5 * * * *"
    enabled = true
  }
  virtual_memory {
    cron = "0 */5 * * * *"
    enabled = true
  }
  swap_memory {
    cron = "0 */5 * * * *"
    enabled = true
  }
  disk {
    cron = "0 */5 * * * *"
    enabled = true
  }
  network {
    cron = "0 */5 * * * *"
    enabled = true
  }
  // unsupported on Windows
  // load {
  //   cron = "0 */5 * * * *"
  //   enabled = false
  // }
}

postgres {
  pg_stat_statements {
    cron = "@hourly"
    enabled = true
  }
  pg_stat_activity {
    cron = "@hourly"
    enabled = true
  }
  pg_stat_user_tables {
    cron = "@hourly"
    enabled = true
  }
  pg_stat_user_indexes {
    cron = "@hourly"
    enabled = true
  }
  pg_locks {
    cron = "@hourly"
    enabled = true
  }
  pg_stat_archiver {
    cron = "@hourly"
    enabled = true
  }
}
`

var defaultDatasourceConfig = `//
// ---------------------------------
// statsender - datasource configuration
// ---------------------------------
//
// Example datasource block structures:
//
// datasource {
//   host = "10.0.3.13"
//   port = 5432
//   username = "user"
//   password = "pass"
//   dbname = "proddb"
//   sslmode = "verify-ca"
//   sslcert = "~/.postgreqsl/client.crt"
//   sslkey = "~/.postgreqsl/client.key"
//   sslrootkey = "~/.postgreqsl/server.key"
//   tags = {
//     env = "prod",
//     iteration = "v1"
//   }
// }
//
// datasource {
//   host = "localhost"
//   port = 6432
//   username = "user"
//   password = "pass"
//   dbname = "testdb"
//   sslmode = "disable"
//   tags = {
//     env = "dev",
//     iteration = "v1"
//   }
// }
//
// The required attributes are:
// - host
// - port
// - username
// - password
// - dbname
//
// Default value of sslmode is 'require'
//
`

var defaultSenderConfig = `//
// ---------------------------------
// statsender - sender configuration
// ---------------------------------
//
// Example sender block structures:
//
// Console sender - useful for debugging purposes
// console {}
//
// Http sender - used to send statistics to remote http locations
// http {
//   target = "https://10.0.1.3:8443/stats"
//   retryDelay = 7
//   maxRetries = 2
//   rootCAs = ["C:\\Path\\to\\your\\rootCA.pem", "and\\optionally\\another\\ones"]
// }
//
`
