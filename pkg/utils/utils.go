package utils

import (
    "log"
)

func LogError(message string, err error) {
    if err != nil {
        log.Printf("%s: %v\n", message, err)
    }
}
