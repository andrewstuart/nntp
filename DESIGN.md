# Design

## Connection pool
- Abstracted away. When the client is asked to do something, it gets a connection (or waits for one),
  performs the requested action, and then returns the connection to the pool and returns
  the response to the calling code.
- Returns a reader
  - How do I return a reader (tied to the connection b/c data might still be coming in) 
    and only return the connection to the pool after the reader has been exhausted?
    I want to do this without client code being aware and having to call Done() or 
    some nonsense.
