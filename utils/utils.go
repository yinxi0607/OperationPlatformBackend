package utils

import "k8s.io/apimachinery/pkg/util/intstr"

func Int32Ptr(i int32) *int32 { return &i }

func IntOrStringPtr(i int) *intstr.IntOrString {
	intOrStr := intstr.FromInt(i)
	return &intOrStr
}
