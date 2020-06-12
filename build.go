package main

// #cgo linux CFLAGS: -I/usr/local/libressl/include -Wno-deprecated-declarations
// #cgo linux LDFLAGS: -L/usr/local/libressl/lib -lssl -lcrypto
import "C"
