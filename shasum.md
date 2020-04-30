# Shasum

We can run shasum in OSX, what get sha1 sum value:
```shell script
$ echo '\0hello sha1' | shasum
dcbc883f6bac8566fbd1a8bc0cdc3eb750d5fde0  -
```

If we compute the sha1 sum value in Golang, what text we should passed? It that 'hello sha1'?

Well, it is not. If we run the follow codes:

```shell script
func main() {
    h := sha1.New()
    h.Write([]byte(`\0hello sha1`))
    s := h.Sum(nil)
    fmt.Printf("%x\n", s)
}
```

The result we get is `b56c209e2ad8aabd6fd25ef63fa436ef29cf11e1`.

So, how about the follow codes?
```shell script
func main() {
    h := sha1.New()
    b1 := []byte("")
    b2 := []byte(`hello sha1`)
    b1 = append(b1, b2...)
    h.Write(b1)
    s := h.Sum(nil)
    fmt.Printf("%x\n", s)
}
```
Well, the result is `64faca92dec81be17500f67d521fbd32bb3a6968`, which is still not same as what we get using shasum command.

The problem here is `\0`. It is not an empty string, it is not an empty byte slice. It is the value of `0x00` which is []byte{0x00}

So would replacing empty byte slice with `[]byte{0x00}` give us the correct answer? If we run the following codes:
```shell script
func main() {
    h := sha1.New()
    b1 := []byte{0x00}
    b2 := []byte(`hello sha1`)
    b1 = append(b1, b2...)
    h.Write(b1)
    s := h.Sum(nil)
    fmt.Printf("%x\n", s)
}
``` 

We get the result `522468a98b3c489953d9e7ef59f43fd65d184d11`. It is still not the same as the `shasum` output. Why? Because `echo` 
has added a 'line breaker' at the end of the string. That means, we need append a 'line breaker', which is `\n` in OSX:
```shell script
func main() {
    h := sha1.New()
    b1 := []byte{0x00}
    b2 := []byte("hello sha1\n")
    b1 = append(b1, b2...)
    h.Write(b1)
    s := h.Sum(nil)
    fmt.Printf("%x\n", s)
}
```

Bingo! Finally, we get the same result as `shasum` one.

### Summary
- `\0` which is `0x00` which is not an empty byte slice;
- Do not forget to add a 'line breaker'.
