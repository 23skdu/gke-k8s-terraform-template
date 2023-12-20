#data "kubectl_filename_list" "manifests" {
#    pattern = "./manifests/*.yaml"
#}
#resource "kubectl_manifest" "CRDs" {
#    count     = length(data.kubectl_filename_list.manifests.matches)
#    yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
#}
