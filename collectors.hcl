debug = true

system {
  cpu {
    cron = "* * * * * *"
  }
}

postgres {
  pg_stat_user_indexes {
    cron = "*/2 * * * * *"
  }
}