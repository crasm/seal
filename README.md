shield
======

**Warning: This project is under heavy development. Breaking changes will most
likely be introduced in the spec and the tool itself until 1.0.**

`shield` is a container format that lets you check for file corruption without
dealing with separate checksum files.

Distilled, shield is just a file with a prepended sha512 hash, optionally
truncated.

The header has the following format:

    SHD%v0{cf83e135}

Above is a fully valid shielded file (ending in a newline), where the
contents are completely empty.

Check the examples folder for more.

( I was disturbingly able to trash many bits within squirrel.jpg before the
corruption was visible when I opened the file. Hence why this project must exist. )

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

Nothing yet.

However, the reason why I started this was because I became paranoid of disk rot
when one of my music files became corrupted during a copy, and the best
solutions I found were:

1. Checksum files (messy, easy to become outdated) 
2. parchive (messy, easy to lose)
3. ZFS (not portable)

My criteria was that it needed to be portable and simple to manage, so I
designed my own using linux [shell
commands](https://github.com/crasm/shield.sh). This repository is the evolution
of that proof-of-concept.

I plan on using this to guard my lossless music collection, but I need to
add shield support to mpv (`shield-c`) and beets (`shield-python`) before that can
happen.

(It's possible to just pipe the output of a music or video file from shield to
mpv -- or the player of your choice -- but seeking doesn't work and it's
generally a pain.)

Installation
------------

    ; go get -v github.com/crasm/shield
    ; go install -v github.com/crasm/shield

Usage
-----

Once installed, run shield with no arguments.

    ; shield

Examples
--------

    # Extracts to squirrel.jpg.
    ; shield -X squirrel.jpg.shd 

    # Prints to stdout. (Be careful with binary.)
    ; shield -C < LICENSE

    # Shields the text and then extracts it. (Does a lot of nothing.)
    ; echo 'shield pipe!' | shield -C | shield -X

Mission
-------

1. Make checking file integrity easy and easier to automate, so that it can be
   done across filesystems, networking protocols, and operating systems.
2. Keep the format simple enough to generate, extract, and verify file contents
   "by hand" with basic \*nix tools.

Stretch goals
-------------

- [ ] Signify support as an alternative to sha512.
- [ ] Backup client that uses shield to verify integrity while copying files.
- [ ] Browser plugins and apps for automatic verification and extraction of downloads.
- [ ] HTTP middleware for go. (Shield protected HTML? Why not.)

Manual shield header generation
-------------------------------

To generate a shield file by hand:

    ; sha512sum LICENSE
    53331cbf3149b47ba0be481c1cfd61d60282ce13652909a17a25626...  LICENSE
    ; echo 'SHD%v0{<paste the hash>}' > LICENSE.shd
    ; cat LICENSE >> LICENSE.shd

To check a shielded file by hand:

    ; head -n 1 LICENSE.shd
    ; tail -n +2 LICENSE.shd > LICENSE
    ; sha512sum LICENSE
      <compare the hashes starting from the left>


vim: tw=80 et sw=4 sts=4
