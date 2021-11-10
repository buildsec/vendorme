//
// Copyright 2021 Tim Miller.
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

package cli

import (
	"context"
	"fmt"
)

type PullCommand struct {
		Foo		string
}

func(c *PullCommand) Exec(ctx context.Context, directory []string) (err error){
	fmt.Println("TOTALLY VENDORING RIGHT NOW")

	return nil
}
