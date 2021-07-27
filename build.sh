#!bin/bash
rm .terraform.lock.hcl || true
cd provider

# $1 is equal to NAME
# $2 is equal to version number
go build -o ~/.terraform.d/plugins/github.com/$1/homeauto/$2/linux_amd64/terraform-provider-homeauto_v$2

echo "Built file"
cd ..
terraform init
terraform plan
terraform apply -auto-approve