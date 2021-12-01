
// docker images tags

variable "api_tag" {
  type = string
}

variable "backend_tag" {
  type = string
}

variable "api_replicas" {
  type = number
}

variable "backend_replicas" {
  type = number
}


// envs

variable "API_SHUTDOWN" {
  type = number
}

variable "API_STARTS_COUNT" {
  type = number
}

variable "API_STARTS_DELAY" {
  type = number
}



variable "BACKEND_SHUTDOWN" {
  type = number
}

variable "BACKEND_STARTS_COUNT" {
  type = number
}


variable "BACKEND_STARTS_DELAY" {
  type = number
}

//

variable "API_VERSION" {
  type = string
}

variable "STOP_TIMEOUT" {
  type = number
}


variable "LISTEN_PORT" {
  type = number
}

variable "API_PORT" {
  type = number
}
