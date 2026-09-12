resource "kubernetes_namespace_v1" "prod" {
  metadata {
    name = "prod"
  }
}
