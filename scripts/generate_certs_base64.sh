#!/bin/zsh

# ensure certstrap is setup.
certstap=/usr/local/bin/certstrap
if ! command -v $certstap &> /dev/null
then
    echo "certstap could not be found"
    exit
fi

printFileAndBase64 () {
    echo "\n----------------------"
    echo "Name: $1.$2"
    contents=$(cat ./out/$1.$2)
    echo $contents
    contentsBase64=$(echo $contents | base64)
    echo $contentsBase64
    echo "----------------------\n"
}

saveFile () {
  contents=$(cat ./out/$1.$2)
  echo $contents > $3
}

printBase64 () {
  contents=$(cat ./out/$1.$2)
  contentsBase64=$(echo $contents | base64)
  echo "$3=$contentsBase64"
}

# Delete old certs and keys if any
rm -rf ./out

caName=CandidateTrackerCA
serverName=api.candidatetracker.co
# for local uncomment the below line
# serverName=candidatetracker
clientName=candidatetracker

# Create the CA
certstrap init --passphrase "" --common-name $caName

# Create the certificates for the client and servers
certstrap request-cert --passphrase "" --domain $serverName
certstrap request-cert --passphrase "" --domain $clientName

# Sign the certificates for the client and servers
certstrap sign $serverName --CA $caName
certstrap sign $clientName --CA $caName

if [[ $1 = "debug" ]]
then
  # convert required keys and certs to Base64 and print to stdout
  printFileAndBase64 $caName "crt"
  printFileAndBase64 $serverName "crt"
  printFileAndBase64 $serverName "key"
  printFileAndBase64 $clientName "crt"
  printFileAndBase64 $clientName "key"
else
  # Delete old certs
  rm -rf ./tmp
  mkdir ./tmp

  saveFile $caName "crt" "./tmp/CA.crt"
  saveFile $clientName "crt" "./tmp/client.crt"
  saveFile $clientName "key" "./tmp/client.key"
  echo "---server-env---"
  printBase64 $caName "crt" "CA_CERT_BASE64"
  printBase64 $serverName "crt" "SERVER_CERT_BASE64"
  printBase64 $serverName "key" "SERVER_KEY_BASE64"
  echo "---client-env---"
  printBase64 $caName "crt" "CA_CERT_BASE64"
  printBase64 $clientName "crt" "CLIENT_CERT_BASE64"
  printBase64 $clientName "key" "CLIENT_KEY_BASE64"
fi

# Delete all generated certs and keys for safety
rm -rf ./out
