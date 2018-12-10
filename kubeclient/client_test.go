package kubeclient

import (
	"fmt"
	"testing"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

var host = `https://10.151.32.60:7443`

var ca = `-----BEGIN CERTIFICATE-----
MIIB4zCCAUSgAwIBAgIJANxon9ltV08QMAoGCCqGSM49BAMCMBgxFjAUBgNVBAMM
DWt1YmVybmV0ZXMtY2EwHhcNMTgxMTE0MTQwMDM5WhcNMjgxMTExMTQwMDM5WjAY
MRYwFAYDVQQDDA1rdWJlcm5ldGVzLWNhMIGbMBAGByqGSM49AgEGBSuBBAAjA4GG
AAQAyy0HBCWtQiVwTfPrYm8Ta9kqSbi0A94jZVdDfAZlUPOgPgDxCBurgP9x7IcV
ff3MdU8xd/2C9+/FZEJvnf05454AOMCuAkHgZmq74CtW963c6PHDBxzAqqyYRSQX
K0VrIlaWzure1ToGGcH8lhBPqTXZBT4KAzsylhOt/lqQa9FMP5ijNDAyMA8GA1Ud
EwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgKkMA8GA1UdEQQIMAaHBAqXIDwwCgYI
KoZIzj0EAwIDgYwAMIGIAkIBEZDUc4tUX7lnZETO1FZH3/xZeK3hWiNcHYy5BPlS
hcKNCJ0jL5Ypmd0QpjG+iyTmYeD+guRN3lV8XJXV3vI7vnECQgEIVCTwRcOZdoMD
e1+2GXi2bPGq6+nNT6y28saIYZxQr8f7URjHkoSlTPSY1Tk12allB3lkj/rEdwmj
pHo3IH/eww==
-----END CERTIFICATE-----`

var cert = `-----BEGIN CERTIFICATE-----
MIIB+zCCAVygAwIBAgIJAMFtUhRCAyg3MAoGCCqGSM49BAMCMBgxFjAUBgNVBAMM
DWt1YmVybmV0ZXMtY2EwHhcNMTgxMTE0MTQwMDQwWhcNMjgxMTExMTQwMDQwWjAy
MRcwFQYDVQQDDA5rdWJlbGV0LWNsaWVudDEXMBUGA1UECgwOc3lzdGVtOm1hc3Rl
cnMwgZswEAYHKoZIzj0CAQYFK4EEACMDgYYABACPKLsr0+6Oqr2fEHkVcMOS2+Jl
HVYZgb7B4Z1mLCWp9pY/cUfH5nL6zYewZLe+iKINrbmYO0LhfOM6kn7K4ZkRzgD6
sz5Gh0j3mHQRlOeYI7CwK2Yja3rcU16TPQYrZxxtms9gi8a2Sc1FX8VBuuzLt81v
gyrN0jj4zsZFsEX+2dH71qMyMDAwCQYDVR0TBAIwADAOBgNVHQ8BAf8EBAMCBaAw
EwYDVR0lBAwwCgYIKwYBBQUHAwIwCgYIKoZIzj0EAwIDgYwAMIGIAkIAnZ5yoUOy
NwxckcoVh9M2ZdlN8V9Aul/7204JnAvtIqIrhQa7yhb37n0Ty8LllK5lQZ3euj37
Qgw3oRamvAHSO/0CQgHjZ16SD8L8Y/ND36s0uOlKdVxLNPVzKwxi6tr+WuUJ5RSk
MIrnahvSIXoNJheUGtYRy10SNUFc4l0oFt05qSzbpA==
-----END CERTIFICATE-----`

var key = `-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIBFcss4iXFHAkv2lgOGd4UAc/xvnmioG3rBH+067vZdqydHJhY0P7z
LN2vNXnI+k2f/aGmYAZ/N8qJs3R+sQYwBdigBwYFK4EEACOhgYkDgYYABACPKLsr
0+6Oqr2fEHkVcMOS2+JlHVYZgb7B4Z1mLCWp9pY/cUfH5nL6zYewZLe+iKINrbmY
O0LhfOM6kn7K4ZkRzgD6sz5Gh0j3mHQRlOeYI7CwK2Yja3rcU16TPQYrZxxtms9g
i8a2Sc1FX8VBuuzLt81vgyrN0jj4zsZFsEX+2dH71g==
-----END EC PRIVATE KEY-----`

func Test_NewKubeClient(t *testing.T) {
	kube := NewKubeClient(host, ca, cert, key)
	if kube == nil {
		t.Fatalf("NewKubeClient error")
	}
	list, err := kube.Client.CoreV1().Pods(v1.NamespaceAll).List(v1.ListOptions{})
	if err != nil {
		t.Errorf("Client.CoreV1() list all pods")
		return
	}

	for _, item := range list.Items {
		name, err := cache.MetaNamespaceKeyFunc(&item)
		if err != nil {
			t.Errorf("cache.MetaNamespaceKeyFunc(item) error")
			return
		}

		fmt.Printf("%s\n", name)
		/*
		owner := appmachinery.GetOwner(kube.Client, item.Namespace, item.Name)
		if owner != nil {
			fmt.Printf("%s ---> %s\n", name, owner.Name)
		}
		*/
	}
}
