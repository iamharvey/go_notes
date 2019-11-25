There are three essentials w.r.t. `defer`:
- `defer` is usually used for, e.g., closing I/O operations, DB connections, etc.

- it follows the rule LIFO (last-in-first-out). E.g.,
   ```
   func main() {
      defer func() {
          fmt.Println("first defer")
      }()
      
      defer func() {
          fmt.Println("second defer")
      }()
  }
  
  // the result would be
  // ====================
  // second defer
  // first defer
  ```
  
- the only way for deferred func to access returned result is to use a named result parameter. E.g.,
  ```
  func main() {
      fmt.Printf(deferFirst())
      fmt.Printf(deferSecond())
  }
  
  func deferFirst() (x int) {
     x = 1
     defer func() { x = 2 }()
     return x
  }

  func deferSecond() (int) {
     x := 1
     defer func() { x = 2 }()
     return x
  }
  
  // the result would be
  // ========================
  // 2
  // 1
  
  ```
