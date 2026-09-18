resource "google_binary_authorization_policy" "policy" {
  count    = var.enable_binary_authorization ? 1 : 0
  project  = var.project

  admission_whitelist_patterns {
    name_pattern = "gcr.io/${var.project}/*"
  }

  default_admission_rule {
    evaluation_mode  = "REQUIRE_ATTESTATION"
    enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"

    require_attestations_by = [
      google_binary_authorization_attestor.attestor[0].name,
    ]
  }

  global_policy_evaluation_mode = "ENABLE"
}

resource "google_binary_authorization_attestor" "attestor" {
  count    = var.enable_binary_authorization ? 1 : 0
  name     = "${var.environment}-attestor"
  project  = var.project

  attestation_authority_note {
    note_reference = google_binary_authorization_note.note[0].name
  }
}

resource "google_binary_authorization_note" "note" {
  count    = var.enable_binary_authorization ? 1 : 0
  name     = "${var.environment}-note"
  project  = var.project

  attestation_authority {
    hint {
      human_readable_name = "Attestor for ${var.environment} environment"
    }
  }
}
