/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 */

package decorators

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/grpc-application-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func CreateApplicationDescriptor(name string) *grpc_application_go.AppDescriptor {

	descriptor := grpc_application_go.AppDescriptor{
		OrganizationId:  uuid.New().String(),
		AppDescriptorId: uuid.New().String(),
		Name:            name,
		Rules: []*grpc_application_go.SecurityRule{
			{Name: uuid.New().String()},
		},
	}

	return &descriptor
}

var _ = ginkgo.Describe("Helper", func() {

	ginkgo.Context("Order decorator", func() {
		ginkgo.It("should be able to order a list of appDescriptor by Name", func() {

			num := 10
			list := make([]*grpc_application_go.AppDescriptor, 0)
			for i := 0; i < num; i++ {
				list = append(list, CreateApplicationDescriptor(uuid.New().String()))
			}

			decorator := NewOrderDecorator(OrderOptions{Field: "Name", Asc: true})
			AppDescList := &grpc_application_go.AppDescriptorList{
				Descriptors: list,
			}

			res := ApplyDecorator(AppDescList, decorator)
			gomega.Expect(res.Error).Should(gomega.BeNil())
			gomega.Expect(res.AppDescriptorList).ShouldNot(gomega.BeNil())
			gomega.Expect(len(res.AppDescriptorList.Descriptors)).Should(gomega.Equal(num))

			for i := 0; i <= len(res.AppDescriptorList.Descriptors)-2; i++ {
				aux := res.AppDescriptorList.Descriptors[i]
				aux2 := res.AppDescriptorList.Descriptors[i+1]
				fmt.Println(aux.Name, " < ", aux2.Name)
				minor := aux.Name < aux2.Name
				gomega.Expect(minor).Should(gomega.BeTrue())
			}
		})
		ginkgo.It("should not be able to order a list of appDescriptor by OrganizationId", func() {

			num := 10
			list := make([]*grpc_application_go.AppDescriptor, 0)
			for i := 0; i < num; i++ {
				list = append(list, CreateApplicationDescriptor(uuid.New().String()))
			}

			decorator := NewOrderDecorator(OrderOptions{Field: "OrganizationId", Asc: true})
			AppDescList := &grpc_application_go.AppDescriptorList{
				Descriptors: list,
			}

			res := ApplyDecorator(AppDescList, decorator)
			gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
			fmt.Println(res.Error.Error())

		})
	})
})
