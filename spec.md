shield v0
=========

**Note: Not all spec features may be in the reference implementation
yet.**

"shield" is a simple container format that embeds the information needed to
check for corruption (hashes) and/or authenticity information (signatures)
into the file itself, allowing verification across filesystems, networking
protocols, and operating systems.

The goal of shield is to reduce the barrier to checking for file corruption and
file authenticity, and therefore make storing and transferring digital files
provably reliable.

shield also aims to match or exceed the integrity and security guarantees offered by
PGP-signed checksum files, while offering acceptable performance and being
simple to use and understand.

Header Identification
---------------------

The shield magic number is encoded in UTF-8 in the first 4 bytes of the file:


        0                      31
        +-----+-----+-----+-----+
        | 'S' | 'H' | 'D' | '%' |
        +-----+-----+-----+-----+
        | 'v' | '0' | ..........


`SHD%` is the magic number for shield files and can be used to identify them.

The byte immediately following the literal `v` identifies the hex-encoded major
version number of the file. 

Header Structure
----------------

    SHD%v0{variant:<claim>}

Variants and Claims
-----------------------------

There are two shield variants with different properties. The variant determines
how the claim is generated and interpreted.

### sha512

For the `sha512` variant, the claim is the hex-encoded sha512 hash of the given
file. It may optionally be truncated to the first `x` bytes, where `1 ≤ x < 64`.

The following short form may be used:

    SHD%v0{<claim>}

### signify

The `signify` variant targets compatibility with OpenBSD's signify tool for
signing and verification. The claim is the base64-encoded signature generated
by signify on a given file.

vim: tw=80 et sw=4 sts=4
