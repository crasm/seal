shield
======

"shield" is a container format that lets you store a sha256 hash with a file,
allowing file integrity to be verified across filesystems, networking
protocols, and operating systems.

This is a work-in-progress. See [shield.sh][] for a limited, but functional
version. See the [spec][] for details and planned features.

[shield.sh]: https://github.com/crasm/shield.sh
[spec]: https://github.com/crasm/shield-spec

Header
------

A shield header has the following form.

    SHD%v0{e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855}

The total length of the header is 72 characters plus a newline.
