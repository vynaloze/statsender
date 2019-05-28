datasource {
  host = "localhost"
  port = 5432
  username = "postgres"
  password = "hasL0"
  dbname = "foo"
  sslmode = "disable"
  tags = {
    env = "prod",
    iteration = "v1"
  }
}

//datasource {
//  host = "localhost"
//  port = 6432
//  username = "user"
//  password = "pass"
//  dbname = "testdb"
//  sslmode = "disable"
//  tags = { env = "dev", iteration = "v1" }
//}