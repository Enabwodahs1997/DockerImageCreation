package main

import (
    "flag"
    "fmt"
    "time"
)

func main() {
	typePtr := flag.String("type", "cos", "type of curve")
	durationPtr := flag.String("duration", "1d", "duration of the generation")
	stepPtr := flag.String("step", "1h", "step / sampling period")
	webMode := flag.Bool("web", false, "run the web interface")
	portPtr := flag.String("port", "8080", "web server port")

	linearFirstValue := flag.Float64("first", 0, "first value for linear type")
	linearLastValue := flag.Float64("last", 1, "last value for linear type")

	cosMinValue := flag.Float64("min", 10, "min value for cos type")
	cosMaxValue := flag.Float64("max", 25, "max value for cos type")
	cosPeriod := flag.String("period", "1d", "period for cos type")

	flag.Parse()

	if *webMode {
		if err := startWebServer(*portPtr); err != nil {
			panic(err)
		}
		return
	}

	durationSeconds := GetSeconds(*durationPtr)
	stepSeconds := GetSeconds(*stepPtr)
	if durationSeconds <= 0 {
		panic("duration must be greater than zero")
	}
	if stepSeconds <= 0 {
		panic("step must be greater than zero")
	}
	if stepSeconds > durationSeconds {
		panic("step must not be greater than duration")
	}
	itemNbr := durationSeconds / stepSeconds

	fn := func(x float64) float64 { return 1 }
	start := time.Now().Unix()

	switch *typePtr {
	case "linear":
		fn = GetLinear(*linearFirstValue, *linearLastValue, start, durationSeconds)
	case "cos":
		fn = GetCosinus(*cosMinValue, *cosMaxValue, *cosPeriod)
	case "log":
		fn = GetLog(start)
	case "exp":
		fn = GetExp(start, durationSeconds)
	default:
		panic("uncorrect function type")
	}

	for i := 0; i < itemNbr; i++ {
		ts := start + int64(i*stepSeconds)
		fmt.Printf("%d %.2f\n", ts, fn(float64(ts)))
	}
}
