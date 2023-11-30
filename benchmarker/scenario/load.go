package scenario

import (
	"context"

	"github.com/isucon/isucandar/worker"
	"github.com/rosylilly/performania/benchmarker/scenario/model"
	"github.com/rosylilly/performania/benchmarker/scenario/testdata"
)

func LoadIconImages(ctx context.Context) ([]*model.Image, error) {
	dirents, err := testdata.IconFiles.ReadDir("icons")
	if err != nil {
		return nil, err
	}

	imgs := make([]*model.Image, len(dirents))

	wrk, err := worker.NewWorker(func(ctx context.Context, i int) {
		file, err := testdata.IconFiles.Open("icons/" + dirents[i].Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, err := model.NewImage(file)
		if err != nil {
			panic(err)
		}

		imgs[i] = img
	},
		worker.WithLoopCount(int32(len(dirents))),
		worker.WithMaxParallelism(20),
	)
	if err != nil {
		return imgs, err
	}

	wrk.Process(ctx)

	return imgs, nil
}

func LoadCoverImages(ctx context.Context) ([]*model.Image, error) {
	dirents, err := testdata.CoverFiles.ReadDir("covers")
	if err != nil {
		return nil, err
	}

	imgs := make([]*model.Image, len(dirents))

	wrk, err := worker.NewWorker(func(ctx context.Context, i int) {
		file, err := testdata.CoverFiles.Open("covers/" + dirents[i].Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, err := model.NewImage(file)
		if err != nil {
			panic(err)
		}

		imgs[i] = img
	},
		worker.WithLoopCount(int32(len(dirents))),
		worker.WithMaxParallelism(20),
	)
	if err != nil {
		return imgs, err
	}

	wrk.Process(ctx)

	return imgs, nil
}
