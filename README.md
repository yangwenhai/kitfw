# kitfw
A micro service framework built on go-kit that uses grpc and capnp as a network protocol and is compatible with etcd and zipkin

# Third party service dependency

## zipkin:

$ docker pull openzipkin/zipkin

$ docker run -d -p 9411:9411 openzipkin/zipkin     

## etcd:

$ docker pull quay.io/coreos/etcd:v2.2.0    

$ HostIP="192.168.88.68"

$ docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
 --name etcd quay.io/coreos/etcd:v2.2.0 \
 -name etcd0 \
 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${HostIP}:2380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${HostIP}:2380 \
 -initial-cluster-state new

# build  

$ git clone https://github.com/yangwenhai/kitfw.git

$ cd kitfw && source devenv.sh

$ cd src/kitfw/vendor && govendor sync

$ cd ../../../

$ go build kitfw/sg/server

$ go build kitfw/sg/client 


# run

./server -zipkinAddr=http://192.168.88.68:9411/api/v1/spans -etcdAddr=http://192.168.99.102:2379

./client 10001 hello kitfw

# result

 ![image](https://github.com/yangwenhai/kitfw/blob/master/image/server.png)

 ![image](https://github.com/yangwenhai/kitfw/blob/master/image/client.png)

 ![image](https://github.com/yangwenhai/kitfw/blob/master/image/zipkin.png)

 ![image](https://github.com/yangwenhai/kitfw/blob/master/image/zipkin_detail.png.png)