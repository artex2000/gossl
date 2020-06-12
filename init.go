package main

/*
#include <string.h>

#include <openssl/conf.h>

#include <openssl/bio.h>
#include <openssl/crypto.h>
#include <openssl/engine.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/ssl.h>

int openssl_init(void) {
	int rc = 0;

	OPENSSL_config(NULL);
	ENGINE_load_builtin_engines();
	SSL_load_error_strings();
	SSL_library_init();
	OpenSSL_add_all_algorithms();

	return 0;
}
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

func init() {
	if rc := C.openssl_init(); rc != 0 {
		panic(fmt.Errorf("openssl_init failed with %d", rc))
	}
}

// errorFromErrorQueue needs to run in the same OS thread as the operation
// that caused the possible error
func errorFromErrorQueue() error {
	var errs []string
	for {
		err := C.ERR_get_error()
		if err == 0 {
			break
		}
		errs = append(errs, fmt.Sprintf("%s:%s:%s",
			C.GoString(C.ERR_lib_error_string(err)),
			C.GoString(C.ERR_func_error_string(err)),
			C.GoString(C.ERR_reason_error_string(err))))
	}
	return fmt.Errorf("SSL errors: %s", strings.Join(errs, "\n"))
}

type GoBIO C.BIO

func readBio(bio *GoBIO) ([]byte, error) {
	r := make([]byte, 0, 0)
	tmp := make([]byte, 4096)

	for {
		n := int(C.BIO_read((*C.BIO)(bio), unsafe.Pointer(&tmp[0]), C.int(len(tmp))))
		if n < 0 {
			return nil, fmt.Errorf("Failed to read BIO")
		} else if n == 0 {
			return r, nil
		} else if n < 4096 {
			tmp = tmp[:n]
			r = append(r, tmp...)
			return r, nil
		} else {
			r = append(r, tmp...)
		}
	}
}
