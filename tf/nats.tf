resource "docker_image" "nats" {
  name         = "nats:latest"
  keep_locally = true
}


resource "docker_service" "nats" {
  name = "nats"

  mode {
    replicated {
      replicas = 1
    }
  }

  task_spec {
    force_update = 0
    networks = [
      docker_network.bus_network.id
    ]
    runtime = "container"

    container_spec {
      image     = docker_image.nats.repo_digest
      isolation = "default"
      read_only = false
    }

    placement {
      max_replicas = 1
      platforms {
        architecture = "amd64"
        os           = "linux"
      }
    }

    resources {
      limits {
        memory_bytes = 536870912
        nano_cpus    = 2000000000
      }

      reservation {
        memory_bytes = 214748364
        nano_cpus    = 500000000
      }
    }

    restart_policy {
      condition    = "any"
      delay        = "3s"
      max_attempts = 20
      window       = "10s"
    }
  }

  update_config {
    delay             = "3s"
    failure_action    = "pause"
    max_failure_ratio = "0.0"
    monitor           = "0s"
    order             = "start-first"
    parallelism       = 1
  }
}
