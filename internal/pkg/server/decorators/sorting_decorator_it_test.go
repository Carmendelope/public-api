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
func CreateApplicationInstance(name string) *grpc_application_go.AppInstance {

	instance := grpc_application_go.AppInstance{
		OrganizationId:  uuid.New().String(),
		AppDescriptorId: uuid.New().String(),
		AppInstanceId:   uuid.New().String(),
		Name:            name,
	}

	return &instance
}

var _ = ginkgo.Describe("Helper", func() {

	ginkgo.Context("Order decorator", func() {
		ginkgo.It("should be able to order a list of appDescriptor by Name", func() {

			num := 10
			list := make([]*grpc_application_go.AppDescriptor, 0)
			for i := 0; i < num; i++ {
				list = append(list, CreateApplicationDescriptor(uuid.New().String()))
			}

			decorator := NewOrderDecorator(OrderOptions{Field: "name", Asc: true})

			res := ApplyDecorator(list, decorator)
			gomega.Expect(res.Error).Should(gomega.BeNil())
			gomega.Expect(res.AppDescriptorList).ShouldNot(gomega.BeNil())
			gomega.Expect(len(res.AppDescriptorList)).Should(gomega.Equal(num))

			for i := 0; i <= len(res.AppDescriptorList)-2; i++ {
				aux := res.AppDescriptorList[i]
				aux2 := res.AppDescriptorList[i+1]
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

			res := ApplyDecorator(list, decorator)
			gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
			fmt.Println(res.Error.Error())

		})
		ginkgo.It("should not be able to order a list of appInstance by any field", func() {

			num := 1
			list := make([]*grpc_application_go.AppInstance, 0)
			for i := 0; i < num; i++ {
				list = append(list, CreateApplicationInstance(uuid.New().String()))
			}

			decorator := NewOrderDecorator(OrderOptions{Field: "OrganizationId", Asc: true})

			res := ApplyDecorator(list, decorator)
			gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
			fmt.Println(res.Error.Error())

		})
	})
})
