FROM progrium/busybox
RUN mkdir -p /etc/ssl && mkdir -p /etc/ssl/certs
ADD certs /etc/ssl/certs/
ADD create_cluster create_cluster
RUN chmod +x create_cluster
CMD "./create_cluster"
#docker run -e '/etc/ssl/certs:/etc/ssl/certs'  -e ETCD_API=172.17.42.1:4001 -e DIGITALOCEAN_TOKEN=a37a4ba5a6ab6a9140bc2d1950776e901db71139fa59797ddd4deba57f5feabf  -e REGION=nyc3 -e "SSH_KEY_NAME=macbook air" -e NODE_COUNT=1 -e VM_SIZE=512mb cakkineni/digital-ocean
