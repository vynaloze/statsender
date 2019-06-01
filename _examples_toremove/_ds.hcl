//
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