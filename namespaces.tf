resource "kubernetes_namespace" "prod" {
  metadata {
    annotations = {
      name = "prod"
    } 
    name = "prod"
  }
}
