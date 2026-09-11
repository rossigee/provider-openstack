//go:build generate

//go:generate rm -rf ../package/crds
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen crd:allowDangerousTypes=true,crdVersions=v1 output:artifacts:config=../package/crds paths=./...
//go:generate go run -tags generate github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets --header-file=../hack/boilerplate.go.txt ./...

package apis
