output "namespace_names" {
  description = "Created namespace names"
  value       = [for ns in kubernetes_namespace_v1.namespaces : ns.metadata[0].name]
}
