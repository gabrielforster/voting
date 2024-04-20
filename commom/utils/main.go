package utils

import "crypto/sha1"

func StringToHash(s string) string {
    hasher := sha1.New()
    hasher.Write([]byte(s))
    return string(hasher.Sum(nil))
}

func Ternary[T any](condition bool, truthy T, notTruthy T) T {
    if condition == true {
        return truthy
    }

    return notTruthy
}
