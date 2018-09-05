package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func rundaemon() error {
	start := time.Now()

	tag, err := name.NewTag("dgodd/some-build-image", name.WeakValidation) // "Size": 135M
	if err != nil {
		return err
	}

	log.Printf("NewTag took %s", time.Since(start))
	start = time.Now()

	image, err := daemon.Image(tag)
	if err != nil {
		return err
	}

	log.Printf("Image Read took %s", time.Since(start))
	start = time.Now()

	if _, err := daemon.Write(tag, image, daemon.WriteOptions{}); err != nil {
		return err
	}

	log.Printf("Image Write took %s", time.Since(start))
	// start = time.Now()

	return nil
}

func runregistry() error {
	start := time.Now()

	ref, err := name.ParseReference("dgodd/some-build-image", name.WeakValidation)
	if err != nil {
		return err
	}

	log.Printf("ParseReference took %s", time.Since(start))
	start = time.Now()

	auth, err := authn.DefaultKeychain.Resolve(ref.Context().Registry)
	if err != nil {
		return err
	}

	log.Printf("Auth creation took %s", time.Since(start))
	start = time.Now()

	image, err := remote.Image(ref, remote.WithAuth(auth))
	if err != nil {
		return err
	}

	log.Printf("Image Read took %s", time.Since(start))
	start = time.Now()

	if err := remote.Write(ref, image, auth, http.DefaultTransport, remote.WriteOptions{}); err != nil {
		return err
	}

	log.Printf("Image Write took %s", time.Since(start))
	// start = time.Now()

	return nil
}

func main() {
	fmt.Println("*** DAEMON ***")
	if err := rundaemon(); err != nil {
		panic(err)
	}

	fmt.Println("\n*** REGISTRY ***")
	if err := runregistry(); err != nil {
		panic(err)
	}
}
