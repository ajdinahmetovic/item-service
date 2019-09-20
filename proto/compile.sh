svc=user

protoc $svc.proto --go_out=plugins=grpc:.
ls $svc.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

mkdir -p $1

mv $svc.pb.go $1/$svc.pb.go
cp $svc.proto $1/$svc.proto

echo $svc.proto compiled