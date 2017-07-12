// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skydb

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAuthInfo(t *testing.T) {
	info := NewAuthInfo("secret")

	if info.ID == "" {
		t.Fatalf("got empty info.ID, want non-empty string")
	}

	if bytes.Equal(info.HashedPassword, nil) {
		t.Fatalf("got info.HashPassword = %v, want non-empty value", info.HashedPassword)
	}
}

func TestNewAnonymousAuthInfo(t *testing.T) {
	info := NewAnonymousAuthInfo()

	if info.ID == "" {
		t.Fatalf("got info.ID = %v, want \"\"", info.ID)
	}

	if len(info.HashedPassword) != 0 {
		t.Fatalf("got info.HashPassword = %v, want zero-length bytes", info.HashedPassword)
	}
}

func TestNewProviderInfoAuthInfo(t *testing.T) {
	k := "com.example:johndoe"
	v := map[string]interface{}{
		"hello": "world",
	}

	Convey("Test Provied ProviderInfo", t, func() {
		info := NewProviderInfoAuthInfo(k, v)
		So(info.ProviderInfo[k], ShouldResemble, v)
		So(len(info.HashedPassword), ShouldEqual, 0)
	})
}

func TestAuthData(t *testing.T) {
	Convey("Test AuthData", t, func() {
		Convey("valid AuthData", func() {
			So(AuthData{
				"username": "johndoe",
			}.IsValid(), ShouldBeTrue)

			So(AuthData{
				"email": "johndoe@example.com",
			}.IsValid(), ShouldBeTrue)

			So(AuthData{
				"username": "johndoe",
				"email":    "johndoe@example.com",
			}.IsValid(), ShouldBeTrue)
		})

		Convey("invalid AuthData", func() {
			So(AuthData{}.IsValid(), ShouldBeFalse)
			So(AuthData{
				"iamyourfather": "johndoe",
			}.IsValid(), ShouldBeFalse)
			So(AuthData{
				"username": nil,
			}.IsValid(), ShouldBeFalse)
		})

		Convey("empty AuthData", func() {
			So(AuthData{}.IsEmpty(), ShouldBeTrue)
			So(AuthData{
				"username": nil,
			}.IsEmpty(), ShouldBeTrue)
			So(AuthData{
				"iamyourfather": "johndoe",
			}.IsEmpty(), ShouldBeFalse)
		})
	})
}

func TestSetPassword(t *testing.T) {
	info := AuthInfo{}
	info.SetPassword("secret")
	err := bcrypt.CompareHashAndPassword(info.HashedPassword, []byte("secret"))
	if err != nil {
		t.Fatalf("got err = %v, want nil", err)
	}
	if info.TokenValidSince == nil {
		t.Fatalf("got info.TokenValidSince = nil, want non-nil")
	}
	if info.TokenValidSince.IsZero() {
		t.Fatalf("got info.TokenValidSince.IsZero = true, want false")
	}
}

func TestIsSamePassword(t *testing.T) {
	info := AuthInfo{}
	info.SetPassword("secret")
	if !info.IsSamePassword("secret") {
		t.Fatalf("got AuthInfo.HashedPassword = %v, want a hashed \"secret\"", info.HashedPassword)
	}
}

func TestGetSetProviderInfoData(t *testing.T) {
	Convey("Test Get/Set ProviderInfo Data", t, func() {
		k := "com.example:johndoe"
		v := map[string]interface{}{
			"hello": "world",
		}

		Convey("Test Set ProviderInfo", func() {
			info := AuthInfo{}
			info.SetProviderInfoData(k, v)

			So(info.ProviderInfo[k], ShouldResemble, v)
		})

		Convey("Test nonexistent Get ProviderInfo", func() {
			info := AuthInfo{
				ProviderInfo: ProviderInfo{},
			}

			So(info.GetProviderInfoData(k), ShouldBeNil)
		})

		Convey("Test Get ProviderInfo", func() {
			info := AuthInfo{
				ProviderInfo: ProviderInfo(map[string]map[string]interface{}{
					k: v,
				}),
			}

			So(info.GetProviderInfoData(k), ShouldResemble, v)
		})

		Convey("Test Remove ProviderInfo", func() {
			info := AuthInfo{
				ProviderInfo: ProviderInfo(map[string]map[string]interface{}{
					k: v,
				}),
			}

			info.RemoveProviderInfoData(k)
			v, _ = info.ProviderInfo[k]
			So(v, ShouldBeNil)
		})
	})
}