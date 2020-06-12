package main

/*

#include <openssl/bn.h>
#include <openssl/bio.h>
#include <openssl/rsa.h>
#include <openssl/pem.h>

int BIO_reset_wrapper(BIO* bio) { return BIO_reset(bio); }

*/
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	pubFile     = flag.String("pb", "public_key.pem", "public key storage filename")
	privateFile = flag.String("pr", "key.pem", "private key storage filename")
)

var workDir string = ""

func init() {
	flag.Parse()
}

func main() {
	bn := C.BN_new()
	defer C.BN_free(bn)

	r := C.BN_set_word(bn, C.ulong(0x10001))
	if r != 1 {
		fmt.Println("BN_set_word failed")
		return
	}

	rsa := C.RSA_new()
	defer C.RSA_free(rsa)

	r = C.RSA_generate_key_ex(rsa, C.int(2048), bn, nil)
	if r != 1 {
		fmt.Println("RSA_generate_key_ex failed")
		return
	}

	bio := C.BIO_new(C.BIO_s_mem())
	if bio == nil {
		fmt.Println("BIO_new failed")
		return
	}
	defer C.BIO_free(bio)

	//process public key
	//export public key to memory BiO
	r = C.PEM_write_bio_RSAPublicKey(bio, rsa)
	if r != 1 {
		fmt.Println("PEM_write_bio_RSAPublicKey failed")
		return
	}

	//grab data from memory BiO
	pub_key, err := readBio((*GoBIO)(bio))
	if err != nil {
		log.Fatal(err)
	}

	//write public key to file
	pub_f, err := os.Create(getAbsPath(*pubFile))
	if err != nil {
		log.Fatal(err)
	}
	_, err = pub_f.Write(pub_key)
	if err != nil {
		log.Fatal(err)
	}

	//process private key
	C.BIO_reset_wrapper(bio)
	r = C.PEM_write_bio_RSAPrivateKey(bio, rsa, nil, nil, 0, nil, nil)
	if r != 1 {
		fmt.Println("PEM_write_bio_RSAPrivateKey failed")
		return
	}

	//grab data from memory BiO
	pr_key, err := readBio((*GoBIO)(bio))
	if err != nil {
		log.Fatal(err)
	}

	//write key to file
	pr_f, err := os.Create(getAbsPath(*privateFile))
	if err != nil {
		log.Fatal(err)
	}
	_, err = pr_f.Write(pr_key)
	if err != nil {
		log.Fatal(err)
	}
}

func getAbsPath(fpath string) string {
	var err error

	if workDir == "" {
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	if ok := filepath.IsAbs(fpath); !ok {
		return path.Join(workDir, fpath)
	} else {
		return fpath
	}
}
