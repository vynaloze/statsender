//
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
//   url = "http://10.0.1.3:8080"
//   endpoint = "/stats"
//   retryDelay = 7
//   maxRetries = 2
// }
//
