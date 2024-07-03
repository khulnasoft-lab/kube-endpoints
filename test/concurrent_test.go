/*
Copyright 2022 iconmamundentist@gmail.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"testing"

	"github.com/khulnasoft-lab/kube-endpoints/apis/network/v1beta1"
	"github.com/khulnasoft-lab/kube-endpoints/test/testhelper"
)

func TestAddCrs(t *testing.T) {
	cep := &v1beta1.ClusterEndpoint{}
	cep.Namespace = "default"
	testhelper.CreateTestCRs(10, "sealos", cep, t)
}

func TestDeleteCrs(t *testing.T) {
	testhelper.DeleteTestCRs("sealos", "default", t)
}
