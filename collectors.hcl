debug = true

system {
  cpu {
    cron = "* * * * * *"
  }
  virtual_memory {
    cron = "* * * * * *"
  }
  swap_memory {
    cron = "* * * * * *"
  }
  disk_io {
    cron = "* * * * * *"
  }
  disk_usage {
    cron = "* * * * * *"
  }
  network_io {
    cron = "* * * * * *"
  }
  //  load {
  //    cron = "* * * * * *"
  //  }
}

postgres {
  pg_stat_user_indexes {
    cron = "*/2 * * * * *"
  }
}