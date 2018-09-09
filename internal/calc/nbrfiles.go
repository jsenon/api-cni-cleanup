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
	"io/ioutil"
	"time"

	"github.com/rs/zerolog/log"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

var nbrfile int64
var folderoldsize int64
var foldersize int64

// StatsFiles will provide metrics for cni folder: folder size, number of elements
func StatsFiles(ctx context.Context, cnifiles string) error {
	_, span := trace.StartSpan(ctx, "(*api).StatsFiles")
	defer span.End()

	// New Metrics Number of file in CNI Folder
	nbr := stats.Int64("cni/file/nbre", "Number of File", "")
	viewCount := &view.View{
		Name:        "number_count",
		Description: "number of files",
		TagKeys:     nil,
		Measure:     nbr,
		Aggregation: view.LastValue(),
	}
	// New Metrics Size of CNI Folder
	size := stats.Int64("cni/file/size", "Size of folder", "Bytes")
	viewSize := &view.View{
		Name:        "size_bytes",
		Description: "Size of the folder",
		TagKeys:     nil,
		Measure:     size,
		Aggregation: view.LastValue(),
	}
	err := view.Register(viewCount, viewSize)
	if err != nil {
		log.Error().Msgf("Error registering view:", err.Error())
		return err
	}
	view.SetReportingPeriod(10 * time.Second)
	log.Debug().Msgf("Time Reporting: %d", 10*time.Second)

	go func() {
		log.Debug().Msgf("Entering into go func for counting number of file in %s", cnifiles)
		log.Debug().Msgf("Entering into go func for size of folder %s", cnifiles)

		for {
			files, err := ioutil.ReadDir(cnifiles)
			if err != nil {
				log.Fatal().Msgf("Failed to read folder: %v", err)
			}
			nbrfile = 0
			foldersize = 0
			for _, f := range files {
				// log.Debug().Msgf("File Name: %s File size %d", f.Name(), f.Size())
				if !f.IsDir() {
					nbrfile = nbrfile + 1
					folderoldsize = f.Size()
					foldersize = foldersize + folderoldsize
				}
			}
			log.Debug().Msgf("File size %d, File Number %d", foldersize, nbrfile)
			stats.Record(ctx, nbr.M(nbrfile), size.M(foldersize))
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
