package scheduler

import (
	"github.com/jasonlvhit/gocron"
)

type DayConfig struct {
	Start string `yaml:"start_time"`
	End   string `yaml:"end_time"`
}

type Config struct {
	Sunday    *DayConfig `yaml:"sunday,omitempty"`
	Monday    *DayConfig `yaml:"monday,omitempty"`
	Tuesday   *DayConfig `yaml:"tuesday,omitempty"`
	Wednesday *DayConfig `yaml:"wednesday,omitempty"`
	Thursday  *DayConfig `yaml:"thursday,omitempty"`
	Friday    *DayConfig `yaml:"friday,omitempty"`
	Saturday  *DayConfig `yaml:"saturday,omitempty"`
}

func Schedule(config Config, onStart func(), onEnd func()) {

	if config.Monday != nil {

		if err := gocron.Every(1).Monday().At(config.Monday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Monday().At(config.Monday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Tuesday != nil {
		if err := gocron.Every(1).Tuesday().At(config.Tuesday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Tuesday().At(config.Tuesday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Wednesday != nil {
		if err := gocron.Every(1).Wednesday().At(config.Wednesday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Wednesday().At(config.Wednesday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Thursday != nil {
		if err := gocron.Every(1).Thursday().At(config.Thursday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Thursday().At(config.Thursday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Friday != nil {
		if err := gocron.Every(1).Friday().At(config.Friday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Friday().At(config.Friday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Saturday != nil {
		if err := gocron.Every(1).Saturday().At(config.Saturday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Saturday().At(config.Saturday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}

	if config.Sunday != nil {
		if err := gocron.Every(1).Sunday().At(config.Saturday.Start).Do(onStart); err != nil {
			panic(err)
		}
		if err := gocron.Every(1).Sunday().At(config.Sunday.End).Do(onEnd); err != nil {
			panic(err)
		}
	}
}

func Start() chan bool {
	return gocron.Start()
}
