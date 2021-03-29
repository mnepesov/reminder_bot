package service_errors

import "errors"

const UserNotFound = "user not found"
var Timeout = errors.New("timeout")