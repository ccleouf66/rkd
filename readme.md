# rkd
Rancher Kubernetes Downloader, tool to list and download Rancher release for Air Gapped Environment.

## To do
- bug when downloading multiple images in the same archive
- download rke binary to complet package


## List Rancher version
To list the latests version of Rancher:

```bash
# List rancher version 
$ ./rkd list
Num. Name - TagName
0. Release v2.4.8 - v2.4.8
1.  - v2.4.7
2. v2.4.6 - v2.4.6
3. Release v2.3.9 - v2.3.9
4. Release v2.4.5 - v2.4.5
```

## Download latest package
To downlaod the package in order install a Rancher on Air Gapped environment (helm chart + images):

```bash
# Downlaod latest chart & images
./rkd downlaod
```

```bash
# Downlaod only chart v2.4.5
./rkd downlaod --helm v2.4.5
```

```bash
# Downlaod only images v2.4.5
./rkd downlaod --images v2.4.5
```

```bash
# Downlaod chart & images v2.4.7
./rkd downlaod --helm v2.4.7 --images v2.4.7
```

## Building dependencies

It should work in more environments (e.g. for native macOS builds)
It does not require root privileges (after dependencies are installed)
Install the necessary dependencies:

```bash
# Fedora:
$ sudo dnf install gpgme-devel libassuan-devel btrfs-progs-devel device-mapper-devel
```
```bash
# Ubuntu (`libbtrfs-dev` requires Ubuntu 18.10 and above):
$ sudo apt install libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev
```
```bash
# macOS:
$ brew install gpgme
```