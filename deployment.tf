resource "kubernetes_deployment" "nginx" {
  metadata {
    name = "nginx"
    labels = {
      production = "nginx"
    }
    namespace = "prod"
  }
  spec {
    replicas = 2 
    selector {
      match_labels = {
      production = "nginx" 
     }
    }
    template {
      metadata {
        labels = {
        production = "nginx"
        }
      }
      spec {
        container {
          image = "nginx:latest"
          name  = "nginx"
          resources {
            limits = {
              cpu    = "0.25"
              memory = "24Mi"
            }
            requests = {
              cpu    = "0.25"
              memory = "24Mi"
            }
          }
          liveness_probe {
            http_get {
              path = "/"
              port = 80
            }
            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}
