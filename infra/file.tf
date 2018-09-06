variable "credentials" {
  type = "string"
}

provider "google" {
  credentials = "${file("${var.credentials}")}"
  project     = "cf-sandbox-sijones"
  region      = "us-central1"
}

variable "subnet_cidr" {
  type    = "string"
  default = "10.0.0.0/16"
}

resource "google_compute_network" "production-network" {
  name                    = "production-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "production-subnet" {
  name          = "production-subnet"
  ip_cidr_range = "10.0.1.0/24"
  network       = "${google_compute_network.production-network.self_link}"
}

resource "google_compute_firewall" "external" {
  name    = "production-external"
  network = "${google_compute_network.production-network.name}"

  source_ranges = ["0.0.0.0/0"]

  allow {
    ports    = ["22", "443", "8080", "8000"]
    protocol = "tcp"
  }

  target_tags = ["production"]
}

resource "google_compute_global_address" "prod-lb" {
  name = "prod-lb"
}

resource "google_compute_instance_group" "httplb" {
  // Count based on number of AZs
  count       = 1
  name        = "httpslb"
  description = "terraform generated instance group that is multi-zone for https loadbalancing"
  zone        = "us-central1-a"
  instances = [
    "${google_compute_instance.app_server_cert_issue.self_link}",
    "${google_compute_instance.app_server_ip_issue.self_link}",
    "${google_compute_instance.app_server_cron_app.self_link}",
  ]

}

resource "google_compute_http_health_check" "prod-lb-public" {
  name                = "prod-lb-public"
  port                = 8080
  request_path        = "/healthcheck"
  check_interval_sec  = 5
  timeout_sec         = 2
  healthy_threshold   = 2
  unhealthy_threshold = 2
}

resource "google_compute_backend_service" "http_lb_backend_service" {
  name        = "httpslb"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 900
  enable_cdn  = false

  backend {
    group = "${google_compute_instance_group.httplb.0.self_link}"
  }

  health_checks = ["${google_compute_http_health_check.prod-lb-public.self_link}"]
}

resource "google_compute_url_map" "https_lb_url_map" {
  name = "prod-lb-http"

  default_service = "${google_compute_backend_service.http_lb_backend_service.self_link}"
}

resource "google_compute_target_http_proxy" "http_lb_proxy" {
  name        = "httpproxy"
  description = "really a load balancer but listed as an https proxy"
  url_map     = "${google_compute_url_map.https_lb_url_map.self_link}"
}

resource "google_compute_target_https_proxy" "https_lb_proxy" {
  name             = "httpsproxy"
  description      = "really a load balancer but listed as an https proxy"
  url_map          = "${google_compute_url_map.https_lb_url_map.self_link}"
  ssl_certificates = ["${google_compute_ssl_certificate.cert.self_link}"]
}

resource "google_compute_ssl_certificate" "cert" {
  name_prefix = "lbcert-prod-cert"
  description = "user provided ssl private key / ssl certificate pair"
  certificate = "${file("../certs/mysite/garimash.com.crt")}"
  private_key = "${file("../certs/mysite/garimash.com.key")}"

  lifecycle = {
    create_before_destroy = true
  }
}

resource "google_compute_firewall" "prod-lb-health_check" {
  name    = "prod-lb-health-check"
  network = "${google_compute_network.production-network.name}"

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }

  source_ranges = ["0.0.0.0"]
  target_tags   = ["production"]
}

resource "google_compute_global_forwarding_rule" "prod-lb-http" {
  name       = "prod-lb-lb-http"
  ip_address = "${google_compute_global_address.prod-lb.address}"
  target     = "${google_compute_target_http_proxy.http_lb_proxy.self_link}"
  port_range = "80"
}

resource "google_compute_global_forwarding_rule" "prod-lb-https" {
  name       = "prod-lb-lb-https"
  ip_address = "${google_compute_global_address.prod-lb.address}"
  target     = "${google_compute_target_https_proxy.https_lb_proxy.self_link}"
  port_range = "443"
}

resource "google_compute_address" "app_server_cert_issue" {
  name         = "cert-ip"
}

resource "google_compute_instance" "app_server_cert_issue" {
  name         = "app3"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"
  
  tags = ["production", "http-lb", "http-server", "https-server", "ssh-access"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    subnetwork   = "production-subnet"
    access_config {
          nat_ip = "${google_compute_address.app_server_cert_issue.address}"
    }
  }

}

resource "google_compute_address" "app_server_ip_issue" {
  name         = "ip-ip"
}

resource "google_compute_instance" "app_server_ip_issue" {
  name         = "app2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["production", "http-lb", "http-server", "https-server", "ssh-access"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
  }

  network_interface {
    subnetwork   = "production-subnet"
    access_config {
              nat_ip = "${google_compute_address.app_server_ip_issue.address}"
        }
  }

}


resource "google_compute_address" "app_server_cron_issue" {
  name         = "cron-ip"
}

resource "google_compute_instance" "app_server_cron_app" {
  name         = "app1"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["production", "http-lb", "http-server", "https-server", "ssh-access"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
  }

  network_interface {
    subnetwork   = "production-subnet"
        access_config {
                  nat_ip = "${google_compute_address.app_server_cron_issue.address}"
            }
  }

}

resource "google_compute_address" "pressurevm" {
  name         = "pressurevm"
}

resource "google_compute_instance" "pressurevm" {
  name         = "lbvm"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["production", "http-lb", "http-server", "https-server", "ssh-access"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
  }

  network_interface {
    subnetwork   = "production-subnet"
        access_config {
                  nat_ip = "${google_compute_address.pressurevm.address}"
            }
  }

}
