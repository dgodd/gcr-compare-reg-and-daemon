package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func rundaemon() error {
	totalstart := time.Now()
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
	log.Printf("TOTAL took %s", time.Since(totalstart))

	return nil
}

func runregistry() error {
	totalstart := time.Now()
	start := time.Now()

	ref, err := name.ParseReference("dgodd/some-build-image", name.WeakValidation) // "Size": 135M
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
	log.Printf("TOTAL took %s", time.Since(totalstart))

	return nil
}

func rundaemonapi() error {
	totalstart := time.Now()
	start := time.Now()

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	log.Printf("http.Client took %s", time.Since(start))
	start = time.Now()

	res, err := httpc.Get("http://unix/images/dgodd/some-build-image/get")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("NON-200 code on image get: %d", res.StatusCode)
	}

	log.Printf("get took %s", time.Since(start))
	start = time.Now()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("read body took %s", time.Since(start))
	log.Printf("size of buffer: c. %dM", len(buf)>>20)
	start = time.Now()

	res, err = httpc.Post("http://unix/images/load", "", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("NON-200 code on image get: %d", res.StatusCode)
	}

	log.Printf("write body took %s", time.Since(start))
	log.Printf("TOTAL took %s", time.Since(totalstart))

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

	fmt.Println("\n*** DAEMON API ***")
	if err := rundaemonapi(); err != nil {
		panic(err)
	}
}
