package main

func main() {
    go HandleHttp(8080)
    HandleUpnp()
}
