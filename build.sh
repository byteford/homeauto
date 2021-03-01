#!bin/bash
rm .terraform.lock.hcl || true
cd provider

go build -o ~/.terraform.d/plugins/github.com/byteford/homeauto/$1/darwin_amd64/terraform-provider-homeauto_v$1

echo "Built file"
cd ..
terraform init
terraform plan
terraform apply -auto-approve