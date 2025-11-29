package workout

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type ParsedWorkoutData struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Author        string            `json:"author"`
	SportType     string            `json:"sport_type"`
	TotalDuration int               `json:"total_duration"`
	Segments      []WorkoutSegment  `json:"segments"`
}

type zwoFile struct {
	Author      string     `xml:"author,attr"`
	Name        string     `xml:"name,attr"`
	Description string     `xml:"description,attr"`
	SportType   string     `xml:"sportType,attr"`
	Workout     zwoWorkout `xml:"workout"`
}

type zwoWorkout struct {
	Warmups       []zwoWarmup       `xml:"Warmup"`
	SteadyStates  []zwoSteadyState  `xml:"SteadyState"`
	Cooldowns     []zwoCooldown     `xml:"Cooldown"`
	Intervals     []zwoInterval     `xml:"Interval"`
	Ramps         []zwoRamp         `xml:"Ramp"`
	FreeRides     []zwoFreeRide     `xml:"FreeRide"`
}

type zwoWarmup struct {
	Duration  int     `xml:"Duration,attr"`
	PowerLow  float64 `xml:"PowerLow,attr"`
	PowerHigh float64 `xml:"PowerHigh,attr"`
	Cadence   int     `xml:"Cadence,attr"`
}

type zwoSteadyState struct {
	Duration int     `xml:"Duration,attr"`
	Power    float64 `xml:"Power,attr"`
	Cadence  int     `xml:"Cadence,attr"`
}

type zwoCooldown struct {
	Duration  int     `xml:"Duration,attr"`
	PowerLow  float64 `xml:"PowerLow,attr"`
	PowerHigh float64 `xml:"PowerHigh,attr"`
	Cadence   int     `xml:"Cadence,attr"`
}

type zwoInterval struct {
	Duration  int     `xml:"Duration,attr"`
	PowerLow  float64 `xml:"PowerLow,attr"`
	PowerHigh float64 `xml:"PowerHigh,attr"`
	Cadence   int     `xml:"Cadence,attr"`
}

type zwoRamp struct {
	Duration  int     `xml:"Duration,attr"`
	PowerLow  float64 `xml:"PowerLow,attr"`
	PowerHigh float64 `xml:"PowerHigh,attr"`
	Cadence   int     `xml:"Cadence,attr"`
}

type zwoFreeRide struct {
	Duration int `xml:"Duration,attr"`
	Cadence  int `xml:"Cadence,attr"`
}

func ParseZWO(content []byte) (*ParsedWorkoutData, error) {
	var zwo zwoFile
	err := xml.Unmarshal(content, &zwo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ZWO file: %w", err)
	}

	if zwo.Name == "" {
		return nil, errors.New("workout name is required")
	}

	parsed := &ParsedWorkoutData{
		Name:        zwo.Name,
		Description: zwo.Description,
		Author:      zwo.Author,
		SportType:   zwo.SportType,
		Segments:    []WorkoutSegment{},
	}

	// Parse warmups
	for _, w := range zwo.Workout.Warmups {
		seg := WorkoutSegment{
			Type:      "warmup",
			Duration:  w.Duration,
			PowerLow:  w.PowerLow,
			PowerHigh: w.PowerHigh,
			Cadence:   w.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += w.Duration
	}

	// Parse steady states
	for _, s := range zwo.Workout.SteadyStates {
		seg := WorkoutSegment{
			Type:     "steadystate",
			Duration: s.Duration,
			Power:    s.Power,
			Cadence:  s.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += s.Duration
	}

	// Parse cooldowns
	for _, c := range zwo.Workout.Cooldowns {
		seg := WorkoutSegment{
			Type:      "cooldown",
			Duration:  c.Duration,
			PowerLow:  c.PowerLow,
			PowerHigh: c.PowerHigh,
			Cadence:   c.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += c.Duration
	}

	// Parse intervals
	for _, i := range zwo.Workout.Intervals {
		seg := WorkoutSegment{
			Type:      "interval",
			Duration:  i.Duration,
			PowerLow:  i.PowerLow,
			PowerHigh: i.PowerHigh,
			Cadence:   i.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += i.Duration
	}

	// Parse ramps
	for _, r := range zwo.Workout.Ramps {
		seg := WorkoutSegment{
			Type:      "ramp",
			Duration:  r.Duration,
			PowerLow:  r.PowerLow,
			PowerHigh: r.PowerHigh,
			Cadence:   r.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += r.Duration
	}

	// Parse free rides
	for _, f := range zwo.Workout.FreeRides {
		seg := WorkoutSegment{
			Type:     "freeride",
			Duration: f.Duration,
			Cadence:  f.Cadence,
		}
		parsed.Segments = append(parsed.Segments, seg)
		parsed.TotalDuration += f.Duration
	}

	return parsed, nil
}