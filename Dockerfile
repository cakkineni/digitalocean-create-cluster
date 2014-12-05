FROM progrium/busybox
ADD create_cluster create_cluster
RUN chmod +x create_cluster
CMD create_cluster
#docker run -e API_TOKEN=a37a4ba5a6ab6a9140bc2d1950776e901db71139fa59797ddd4deba57f5feabf -e REGION=nyc3 -e "KEY_NAME=macbook air" -e NODE_COUNT=1 cakkineni\digital-ocean