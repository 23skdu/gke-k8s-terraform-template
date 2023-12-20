# gke-k8s-terraform-template
## skeleton for Google k8s cluster and GCS bucket for TF state
#### you need to populate the vars for project, region, etc
easy create a GKE k8s cluster on their free tier
```
need to comment out bucket, create bucket, then import the bucket into itself,
OR
manuaally create the tf-state bucket, then just use it and have the tf update in place (easier)
```
```
gke-k8s-terraform-template $ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding hashicorp/kubernetes versions matching "2.24.0"...
- Finding gavinbunney/kubectl versions matching ">= 1.7.0"...
- Finding hashicorp/google versions matching "5.8.0"...
- Installing hashicorp/kubernetes v2.24.0...
- Installed hashicorp/kubernetes v2.24.0 (signed by HashiCorp)
- Installing gavinbunney/kubectl v1.14.0...
- Installed gavinbunney/kubectl v1.14.0 (self-signed, key ID AD64217B5ADD572F)
- Installing hashicorp/google v5.8.0...
- Installed hashicorp/google v5.8.0 (signed by HashiCorp)

Partner and community providers are signed by their developers.
If you'd like to know more about provider signing, you can read about it here:
https://www.terraform.io/docs/cli/plugins/signing.html

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!
```
