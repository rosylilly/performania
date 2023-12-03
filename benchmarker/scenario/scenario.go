package scenario

import (
	"context"
	"log"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/parallel"
	"github.com/rosylilly/performania/benchmarker"
	"github.com/rosylilly/performania/benchmarker/scenario/model"
	"github.com/rosylilly/performania/benchmarker/scenario/random"
)

type Scenario struct {
	icons  []*model.Image
	covers []*model.Image
}

func NewScenario() *Scenario {
	return &Scenario{}
}

func (s *Scenario) Prepare(ctx context.Context, step *isucandar.BenchmarkStep) error {
	icons, err := LoadIconImages(ctx)
	if err != nil {
		return err
	}
	log.Printf("loaded %d icons", len(icons))

	covers, err := LoadCoverImages(ctx)
	if err != nil {
		return err
	}
	log.Printf("loaded %d covers", len(covers))

	s.icons = icons
	s.covers = covers

	ag, err := agent.NewAgent(
		agent.WithBaseURL("http://localhost:9292"),
		agent.WithDefaultTransport(),
	)
	if err != nil {
		return err
	}
	req, err := ag.POST("/api/initialize", nil)
	if err != nil {
		return err
	}
	res, err := ag.Do(ctx, req)
	if err != nil {
		return err
	}
	if err := res.Body.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Scenario) Load(ctx context.Context, step *isucandar.BenchmarkStep) error {
	ag, err := agent.NewAgent(
		agent.WithBaseURL("http://localhost:9292"),
		agent.WithDefaultTransport(),
	)
	if err != nil {
		return err
	}

	userCreation := parallel.NewParallel(ctx, 30)
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		userCreation.Do(func(ctx context.Context) {
			user := &model.User{
				Login: random.RandomLogin(),
				Icon:  random.RandomItem(s.icons),
				Cover: random.RandomItem(s.covers),
			}
			contentType, body, err := user.SignupRequestBody()
			req, err := ag.POST("/api/users", body)
			if err != nil {
				step.AddError(err)
				return
			}
			req.Header.Set("Content-Type", contentType)
			res, err := ag.Do(ctx, req)
			if err != nil {
				step.AddError(err)
				return
			}

			benchmarker.DevLogger.Printf("user created: %s", user.Login)

			if err := res.Body.Close(); err != nil {
				step.AddError(err)
				return
			}

			step.AddScore(ScoreUserCreation)
		})
	}
}

func (s *Scenario) Validate(ctx context.Context, step *isucandar.BenchmarkStep) error {
	return nil
}
