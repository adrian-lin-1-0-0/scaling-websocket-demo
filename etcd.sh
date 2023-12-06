docker run \
-p 2379:2379 \
quay.io/coreos/etcd:v3.4.0 \
/usr/local/bin/etcd \
--data-dir /etcd-data \
--listen-client-urls http://0.0.0.0:2379 \
--advertise-client-urls http://0.0.0.0:2379 \
--log-level info \
--logger zap \
--log-outputs stderr
