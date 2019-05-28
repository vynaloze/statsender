debug = false

system {
  cpu {
    cron = "*/2 * * * * *"
  }
  virtual_memory {
    cron = "*/5 * * * * *"
  }
  //  swap_memory {
  //    cron = "* * * * * *"
  //  }
  //  disk_io {
  //    cron = "* * * * * *"
  //  }
  //  disk_usage {
  //    cron = "* * * * * *"
  //  }
  //  network {
  //    cron = "* * * * * *"
  //  }
  //  load {
  //    cron = "* * * * * *"
  //  }
}

postgres {
  pg_stat_user_indexes {
    cron = "*/2 * * * * *"
  }
}