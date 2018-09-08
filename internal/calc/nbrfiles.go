// Copyright © 2018 Julien SENON <julien.senon@gmail.com>
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

// Export Metric of number of CNI Files

package nbrfiles

import (
	"context"

	"github.com/rs/zerolog/log"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var (
	// Size of file
	MSize = stats.Int64("file/size", "The Size of file", "KBs")

	// Counts the number of file
	MFile = stats.Int64("file/nbre", "The number of file in folder", "1")
)

var (
	KeyMethod, _ = tag.NewKey("cni")
)

var (
	FileView = &view.View{
		Name:        "cni/size",
		Measure:     MSize,
		Description: "The size of files",
		Aggregation: view.Count(),
	}

	LineCountView = &view.View{
		Name:        "cni/number",
		Measure:     MFile,
		Description: "The number of files",
		Aggregation: view.Count(),
	}
)

func StatsFiles(ctx context.Context, cnifile string) error {
	_, span := trace.StartSpan(ctx, "(*api).StatsFiles")
	defer span.End()
	ctx2, err := tag.New(context.Background(), tag.Insert(KeyMethod, "file"))

	if err != nil {
		return err
	}
	if err := view.Register(FileView, LineCountView); err != nil {
		log.Fatal().Msgf("Failed to register views: %v", err)
	}
	// Fake Values
	stats.Record(ctx2, MSize.M(1), MFile.M(2))

	// TODO count number of files
	// TODO records size of files

	return nil
}
