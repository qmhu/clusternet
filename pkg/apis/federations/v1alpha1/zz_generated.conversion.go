// +build !ignore_autogenerated

/*
Copyright 2021 The Clusternet Authors.

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
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	url "net/url"

	federations "github.com/clusternet/clusternet/pkg/apis/federations"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Declaration)(nil), (*federations.Declaration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Declaration_To_federations_Declaration(a.(*Declaration), b.(*federations.Declaration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*federations.Declaration)(nil), (*Declaration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_federations_Declaration_To_v1alpha1_Declaration(a.(*federations.Declaration), b.(*Declaration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*Declaration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1alpha1_Declaration(a.(*url.Values), b.(*Declaration), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Declaration_To_federations_Declaration(in *Declaration, out *federations.Declaration, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_Declaration_To_federations_Declaration is an autogenerated conversion function.
func Convert_v1alpha1_Declaration_To_federations_Declaration(in *Declaration, out *federations.Declaration, s conversion.Scope) error {
	return autoConvert_v1alpha1_Declaration_To_federations_Declaration(in, out, s)
}

func autoConvert_federations_Declaration_To_v1alpha1_Declaration(in *federations.Declaration, out *Declaration, s conversion.Scope) error {
	return nil
}

// Convert_federations_Declaration_To_v1alpha1_Declaration is an autogenerated conversion function.
func Convert_federations_Declaration_To_v1alpha1_Declaration(in *federations.Declaration, out *Declaration, s conversion.Scope) error {
	return autoConvert_federations_Declaration_To_v1alpha1_Declaration(in, out, s)
}

func autoConvert_url_Values_To_v1alpha1_Declaration(in *url.Values, out *Declaration, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	return nil
}

// Convert_url_Values_To_v1alpha1_Declaration is an autogenerated conversion function.
func Convert_url_Values_To_v1alpha1_Declaration(in *url.Values, out *Declaration, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1alpha1_Declaration(in, out, s)
}
