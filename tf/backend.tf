resource "docker_image" "backend" {
  name = "127.0.0.1:5000/test-backend:${var.backend_tag}"
}



resource "docker_service" "backend" {
  name = "backend"

  mode {
    replicated {
      replicas = var.backend_replicas
    }
  }

  task_spec {
    force_update = 0
    networks = [
      docker_network.bus_network.id
    ]
    runtime = "container"

    container_spec {
      image     = docker_image.backend.repo_digest
      isolation = "default"
      read_only = false

      stop_signal       = "SIGINT"
      stop_grace_period = "${var.STOP_TIMEOUT}s"

      env = {
        URL              = docker_service.nats.name
        BACKEND_SHUTDOWN = var.BACKEND_SHUTDOWN
        STOP_TIMEOUT     = var.STOP_TIMEOUT
      }

      healthcheck {
        test     = ["CMD", "/prober"]
        interval = "10s"
        timeout  = "2s"
        retries  = 4
      }
    }

    placement {
      max_replicas = 10
      platforms {
        architecture = "amd64"
        os           = "linux"
      }
    }

    resources {
      limits {
        memory_bytes = 214748364
        nano_cpus    = 500000000
      }

      reservation {
        memory_bytes = 214748364
        nano_cpus    = 500000000
      }
    }

    restart_policy {
      condition    = "any"
      delay        = "${var.BACKEND_STARTS_DELAY}s"
      max_attempts = var.BACKEND_STARTS_COUNT
      window       = "2s"
    }
  }

  update_config {
    parallelism       = 1
    delay             = "${var.BACKEND_STARTS_DELAY}s"
    failure_action    = "rollback"
    monitor           = "${var.BACKEND_STARTS_DELAY * (1 + var.BACKEND_STARTS_COUNT)}s"
    max_failure_ratio = "0.0"
    order             = "start-first"
  }

  rollback_config {
    parallelism       = 1
    delay             = "5ms"
    failure_action    = "pause"
    monitor           = "10h"
    max_failure_ratio = "0.9"
    order             = "stop-first"
  }
}
