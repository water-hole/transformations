package kubernetes

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

// Create Transformations for generic annotations
func AddAnnotation(annotations map[string]string) jsonpatch.Patch {
	patchJSON := `[`
	i := 0
	for key, value := range annotations {
		if i == 0 {
			patchJSON = fmt.Sprintf(`%v
{"op": "add", "path": "/metadata/annotations/%v", "value": "%v"}`, patchJSON, key, value)
		} else {
			patchJSON = fmt.Sprintf(`%v,
{"op": "add", "path": "/metadata/annotations/%v", "value": "%v"}`, patchJSON, key, value)
		}

		i++
	}
	patchJSON = fmt.Sprintf("%v]", patchJSON)
	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}

func UpdatePodSpecable(containerImagePath, updatedImage string) jsonpatch.Patch {
	patchJSON := fmt.Sprintf(`[
{ "op": "replace", "path": "%v", "value": "%v"}
]`, containerImagePath, updatedImage)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}

func RemovePodSelectedNode() jsonpatch.Patch {
	patchJSON := fmt.Sprintf(`[
{ "op": "remove", "path": "/spec/nodeName"}
]`)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}

func UpdateNamespace(newNamespace string) jsonpatch.Patch {
	patchJSON := fmt.Sprintf(`[
{"op": "replace", "path": "/namespace", "value": "%v"}
]`, newNamespace)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}

func UpdateRoleBindingSVCACCTNamespace(newNamespace string, numberOfSubjects int) jsonpatch.Patch {
	patchJSON := "["
	for i := 0; i < numberOfSubjects; i++ {
		if i != 0 {
			patchJSON = fmt.Sprintf("%v,", patchJSON)
		}
		patchJSON = fmt.Sprintf(`%v
{"op": "replace", "path": "/subjects/%v/namespace", "value": "%v"}`, patchJSON, i, newNamespace)
	}

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}

func RemoveServiceClusterIPs() jsonpatch.Patch {
	patchJSON := fmt.Sprintf(`[
{"op": "remove", "path": "/spec/clusterIP"}
]`)
	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		fmt.Printf("%v", err)
	}
	return patch
}
