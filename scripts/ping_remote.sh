if [ $# -eq 0 ];
then
  target=10
else
  target=$1
fi

echo $target

sum=0
while [[ sum -lt $target ]];
do
grpcurl -import-path protos -proto server.proto -H 'requesting_user_email: '$TEST_USER_EMAIL -cert certs/remote/client.crt -key certs/remote/client.key -cacert certs/remote/CA.crt api.airetreat.co:9000 protos.CandidateTrackerGo/CheckConnection;

sum=$(($sum+1));
echo $sum
sleep 0.1;

done
