#!bin/bash
rm .terraform.lock.hcl || true
cd provider

go build -o ~/.terraform.d/plugins/github.com/$1/homeauto/$2/linux_amd64/terraform-provider-homeauto_v$2

echo "Built file"
cd ..
terraform init
terraform plan
terraform apply -auto-approve