# rkd
Rancher Kubernetes Downloader, tool to list and download Rancher release for Air Gapped Environment.

## To do
- download rke binary to complet package
- add go minio client to upload datapack
- improve .devcontainer docker image with go extension

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
# Downlaod latest chart & images of rancher
./rkd downlaod
```

```bash
# Downlaod only busybox and alpine images
./rkd downlaod --image busybox --image alpine
```

```bash
# Downlaod rancher v2.5.0 and alpine image
./rkd d --rancher v2.4.8 --image alpine
```

## Auth
rkd uses Docker credential helpers to pull images from a registry.

Get your docker registry user and password encoded in base64
```bash
echo -n USER:PASSWORD | base64
```

Create a <b>config.json</b> file with your Docker registry url and the previous generated base64 string
```json
{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "xxxxxxxxxxxxxxx"
		}
	}
}
```

Run rkd with the config.json inside your home dir.
```bash
~/.docker/config.json
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