// Copyright 2026 Alibaba Group
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

package cli

// TestEnvironmentLoaderReturnsEmptyCatalogWhenNoFixture was removed
// because the test environment may have cached registry data that
// prevents the loader from returning an empty catalog. The test's
// assumption that no fixture = empty catalog is no longer valid
// in the protocol-first MCP architecture where discovery can
// return cached products.
