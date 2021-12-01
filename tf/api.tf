resource "docker_image" "api" {
  name = "127.0.0.1:5000/test-api:${var.api_tag}"
}


resource "docker_service" "api" {
  name = "api"

  endpoint_spec {
    mode = "vip"

    ports {
      protocol       = "tcp"
      publish_mode   = "ingress"
      published_port = 80
      target_port    = 8000
    }
  }


  mode {
    replicated {
      replicas = var.api_replicas
    }
  }

  task_spec {
    force_update = 0
    networks = [
      docker_network.bus_network.id
    ]
    runtime = "container"

    container_spec {
      image     = docker_image.api.repo_digest
      isolation = "default"
      read_only = false

      stop_signal       = "SIGINT"
      stop_grace_period = "${var.STOP_TIMEOUT}s"

      env = {
        URL          = docker_service.nats.name
        API_SHUTDOWN = var.API_SHUTDOWN
        API_VERSION  = var.API_VERSION
        STOP_TIMEOUT = var.STOP_TIMEOUT
      }

      healthcheck {
        test     = ["CMD", "wget", "--spider", "http://127.0.0.1:8000/health"]
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
      condition    = "on-failure"
      delay        = "${var.API_STARTS_DELAY}s"
      max_attempts = var.API_STARTS_COUNT
      window       = "2s"
    }
  }

  update_config {
    parallelism       = 1
    delay             = "${var.API_STARTS_DELAY}s"
    failure_action    = "rollback"
    monitor           = "${var.API_STARTS_DELAY * (var.API_STARTS_COUNT + 1)}s"
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
