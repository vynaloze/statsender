/*
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
Using the second approch, you won't be able to manage
this collector using CLI anymore, though.

*/

system {
  cpu {
    cron = "*/5 * * * * *"
    enabled = true
  }
  virtual_memory {
    cron = "*/5 * * * * *"
    enabled = true
  }
  swap_memory {
    cron = "*/5 * * * * *"
    enabled = true
  }
  disk_io {
    cron = "*/5 * * * * *"
    enabled = true
  }
  disk_usage {
    cron = "*/5 * * * * *"
    enabled = true
  }
  network {
    cron = "*/5 * * * * *"
    enabled = true
  }
  // unsupported on Windows
  load {
    cron = "*/5 * * * * *"
    enabled = true
  }
}

postgres {
  pg_stat_user_indexes {
    cron = "@hourly"
    enabled = true
  }
}
