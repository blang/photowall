Photowall
======

Photowall is a small, functional photowall written in golang.

A photowall is a picture slideshow, running on your TV at a party for example, and guests can upload photos which are added to the slideshow.

The package exposes the main webserver serving the photowall and upload functionality. It supports basic upload validation (size, image) and image processing / resizing.

Current state
-----
The backend is fine, but some frontend parts need polish like missing proper css. If your familiar with js customizing the wall for your need should be quite easy.

Usage
-----
```bash
$ go get github.com/blang/photowall
$ photowall
```

Your images are stored at `$GOPATH/src/github.com/blang/photowall/imgs` by default which is subject to change.

Frontend
-----
[supersized](https://github.com/buildinternet/supersized) is used as the frontend slideshow with modifications to poll the backend server. Also see [LICENSE](LICENSE).

Backend
-----
By default the photowall listens on `127.0.0.1:8000` with following routes:

- `/`: Upload new photos
- `/wall`: View the photowall

License (MIT)
-----

See [LICENSE](LICENSE) file.