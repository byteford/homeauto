#!bin/bash
rm .terraform.lock.hcl || true
cd provider

go build -o ~/.terraform.d/plugins/dgp.com/byteford/homeauto/0.3.1/darwin_amd64/terraform-provider-homeauto_v0.3.1
echo "Built file"
cd ..
terraform init
terraform plan