# Url Shortener

The service implementing the storage
of urls and their short aliases

It is actually just simple CRUD

## The workflow
1. Requester make request via HTTP or gRPC handlers.
2. Handlers read request data and construct domain model from it.
3. If model is broken, handler return error immediately.
4. After that handler calls corresponding use case method.
5. Use case layer can process the data, call corresponding repository method
6. Repository (in memory or SQLite or anything elsa) return requested data to use case layer.
7. Use case layer method send data back to handler.
8. Handler return data to Requester.