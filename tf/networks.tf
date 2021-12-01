resource "docker_network" "bus_network" {
  name   = "bus_network"
  driver = "overlay"
}
