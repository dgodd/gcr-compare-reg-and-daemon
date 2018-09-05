# Simple go-containerregistry benchmark

I wish to investigate the reason behind the performance difference between
daemon and registry use cases for go-containerregistry. The below was generated
by running main.go and shows the times for simply reading an image and then
writing the same image back to the either the docker daemon or registry.

Note that using go-containerregistry against the docker daemon is the slowest
by a significant margin.

I have only just started my own investigation, and have created and posted this
as part of that investigation.

### GO-CONTAINERREGISTRY DAEMON

- NewTag took 7.922µs
- Image Read took 2.128819343s
- Image Write took 13.140210495s
- __TOTAL took 15.269110922s__

### GO-CONTAINERREGISTRY REGISTRY

- ParseReference took 8.732µs
- Auth creation took 135.086µs
- Image Read took 681.436574ms
- ...
- index.docker.io/dgodd/some-build-image:latest: digest: sha256:3baba78fa5ba0ea33f7844de8689217b77c8ec6d9153d74ac8959d0c6b469acb size: 3030
- Image Write took 635.492487ms
- __TOTAL took 1.317129602s__

### DAEMON API

- http.Client took 1.293µs
- get took 1.813008217s
- read body took 368.19261ms
- size of buffer: c. 131M
- write body took 188.48427ms
- __TOTAL took 2.369823875s__
