# containers-from-scratch
Writing a container in a few lines of Go code, as seen at [DockerCon 2017](https://www.youtube.com/watch?v=MHv6cWjvQjM&t=1316s) and on [O'Reilly Safari](https://www.safaribooksonline.com/library/view/how-to-containerize/9781491982310/)

You need root permissions for this to work. 

Also note that the Go code uses some syscall definitions that are only available when building with GOOS=linux.


# Download ubuntu root filesystem

You can download an ubuntu root filesystem from the [Ubuntu Cloud Images](https://cloud-images.ubuntu.com) page.
For example, for bionic:
```
mkdir -p /root/ubuntufs
cd /root/ubuntufs
curl -O -L https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64-root.tar.xz
xz -d bionic-server-cloudimg-amd64-root.tar.xz
tar xf bionic-server-cloudimg-amd64-root.tar
rm -f bionic-server-cloudimg-amd64-root.tar
```
