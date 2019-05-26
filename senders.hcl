console {}

http {
  url = "http://localhost:8000"
  endpoint = "/metric"
  retryDelay = 7
  maxRetries = 2
}

http {
  url = "http://localhost:8000"
  endpoint = "/another"
  retryDelay = 7
  maxRetries = 2
}
