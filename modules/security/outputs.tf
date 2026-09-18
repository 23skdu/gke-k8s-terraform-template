output "policy_id" {
  description = "Binary Authorization policy ID"
  value       = var.enable_binary_authorization ? google_binary_authorization_policy.policy[0].id : ""
}

output "attestor_id" {
  description = "Binary Authorization attestor ID"
  value       = var.enable_binary_authorization ? google_binary_authorization_attestor.attestor[0].id : ""
}
