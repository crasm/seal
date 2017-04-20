seal
====

**Warning: This project is still under development. Breaking changes may
potentially be introduced until 1.0.**

`seal` is a container format that lets you check for file corruption without
dealing with separate checksum files.

Distilled, seal is just a file with a prepended sha512 hash, optionally
truncated.

The header has the following format:

    SL%v0{cf83e135}

Above is a fully valid sealed file (ending in a newline), where the
contents are completely empty.

Check the examples folder for more.

Potential use cases
-------------------

- Catch corruption early and automatically when copying documents to/from flash
  drives or across unreliable networks.
- Detect corruption in large archives and video from disk rot without needing
  filesystem support.
- Checksum entire directories without worrying about where to store the checksum
  files or squishing it all into a tarball.

What I'm using it for
---------------------

Right now, I'm sealing tarballs.

For instance, each set of the [nexus factory images][] come in a ~1GB tarball.
Rather than storing the hash separately, I seal the files. In the event I need
to extract them later (I bricked my phone...), it's easy:

    ; tar vczf - * | seal -Wo bullhead.tgz.sl  # Archive, wrap with a seal
    ; seal -U < bullhead.tgz.sl | tar vxzf -   # Verify, seal and unwrap

[nexus factory images]: https://developers.google.com/android/nexus/images

_tar actually has CRC checksums, so it'll most likely catch any accidental
corruption by itself. However, the performance overhead of seal is less than 5%
(informally benchmarked) for large files, so it doesn't hurt much._

Why
---

The reason why I started seal was because I became paranoid of disk rot
when one of my music files became corrupted during a copy, and the best
solutions I found were:

1. Checksum files (messy, easy to become outdated)
2. parchive (messy, easy to lose)
3. ZFS (not portable)

My criteria was that it needed to be portable and simple to manage, so I
designed my own using linux [shell commands][shield.sh]. This repository is the
evolution of that proof-of-concept.

[shield.sh]: https://github.com/crasm/shield.sh

Installation
------------

    ; go get -v github.com/crasm/seal
    ; go install -v github.com/crasm/seal

Usage
-----

Once installed, run seal with no arguments.

    ; seal

Examples
--------

    # Extracts to LICENSE
    ; seal -U LICENSE.sl

    # Prints to stdout. (Be careful with binary.)
    ; seal -W < LICENSE

    # Seals the text and then extracts it. (Does a lot of... nothing.)
    ; echo 'seal pipe!' | seal -W | seal -U

Mission
-------

1. Make checking file integrity easy and easier to automate, so that it can be
   done across filesystems, networking protocols, and operating systems.
2. Keep the format simple enough to generate, extract, and verify file contents
   "by hand" with basic \*nix tools.

Stretch goals
-------------

- [ ] Signify support as an alternative to sha512.
- [ ] Backup client that uses seal to verify integrity while copying files.
- [ ] Browser plugins and apps for automatic verification and extraction of downloads.
- [ ] HTTP middleware for go. (Sealed HTML? Why not.)

Manual seal generation
----------------------

To generate a seal file by hand:

    ; sha512sum LICENSE
    53331cbf3149b47ba0be481c1cfd61d60282ce13652909a17a25626...  LICENSE
    ; echo 'SL%v0{<paste the hash>}' > LICENSE.sl
    ; cat LICENSE >> LICENSE.sl

To check a sealed file by hand:

    ; head -n 1 LICENSE.sl
    ; tail -n +2 LICENSE.sl > LICENSE
    ; sha512sum LICENSE
      <compare the hashes starting from the left>


vim: tw=80 et sw=4 sts=4
