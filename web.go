package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type webPageData struct {
	Type        string
	Duration    string
	Step        string
	First       string
	Last        string
	Min         string
	Max         string
	Period      string
	Preview     string
	Error       string
	SampleTitle string
}

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GenX Preview</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 0; background: #0f172a; color: #e2e8f0; }
    main { max-width: 960px; margin: 0 auto; padding: 32px 20px 48px; }
    .card { background: #111827; border: 1px solid #334155; border-radius: 16px; padding: 24px; box-shadow: 0 12px 30px rgba(0,0,0,.25); }
    h1 { margin-top: 0; font-size: 2rem; }
    p { line-height: 1.5; color: #cbd5e1; }
    form { display: grid; gap: 12px; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); margin: 20px 0; }
    label { display: grid; gap: 6px; font-size: .9rem; }
    input { background: #0b1220; color: #e2e8f0; border: 1px solid #475569; border-radius: 10px; padding: 10px 12px; }
    button { grid-column: 1 / -1; background: #38bdf8; color: #082f49; border: 0; border-radius: 10px; padding: 12px 16px; font-weight: 700; cursor: pointer; }
    pre { overflow-x: auto; background: #020617; border-radius: 12px; padding: 16px; border: 1px solid #1e293b; color: #f8fafc; }
    .error { color: #fca5a5; margin-bottom: 12px; }
    .hint { font-size: .95rem; color: #94a3b8; }
  </style>
</head>
<body>
  <main>
    <section class="card">
      <h1>GenX Preview</h1>
      <p>Generate cosine, linear, exponential, or logarithmic sample data in the container, or preview it here in the browser.</p>
      {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
      <form method="get">
        <label>Type<input name="type" value="{{.Type}}"></label>
        <label>Duration<input name="duration" value="{{.Duration}}"></label>
        <label>Step<input name="step" value="{{.Step}}"></label>
        <label>First<input name="first" value="{{.First}}"></label>
        <label>Last<input name="last" value="{{.Last}}"></label>
        <label>Min<input name="min" value="{{.Min}}"></label>
        <label>Max<input name="max" value="{{.Max}}"></label>
        <label>Period<input name="period" value="{{.Period}}"></label>
        <button type="submit">Generate preview</button>
      </form>
      <p class="hint">{{.SampleTitle}}</p>
      <pre>{{.Preview}}</pre>
    </section>
  </main>
</body>
</html>`))

func startWebServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	fmt.Printf("web server listening on :%s\n", port)
	return http.ListenAndServe(":"+port, mux)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := buildWebPageData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = err.Error()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if execErr := pageTemplate.Execute(w, data); execErr != nil {
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
	}
}

func buildWebPageData(r *http.Request) (webPageData, error) {
	query := r.URL.Query()
	data := webPageData{
		Type:        defaultQueryValue(query.Get("type"), "cos"),
		Duration:    defaultQueryValue(query.Get("duration"), "2d"),
		Step:        defaultQueryValue(query.Get("step"), "3h"),
		First:       defaultQueryValue(query.Get("first"), "10"),
		Last:        defaultQueryValue(query.Get("last"), "30"),
		Min:         defaultQueryValue(query.Get("min"), "20"),
		Max:         defaultQueryValue(query.Get("max"), "30"),
		Period:      defaultQueryValue(query.Get("period"), "1d"),
		SampleTitle: "type=" + defaultQueryValue(query.Get("type"), "cos") + ", duration=" + defaultQueryValue(query.Get("duration"), "2d") + ", step=" + defaultQueryValue(query.Get("step"), "3h"),
	}

	durationSeconds, err := secondsFromDuration(data.Duration)
	if err != nil {
		return data, err
	}
	stepSeconds, err := secondsFromDuration(data.Step)
	if err != nil {
		return data, err
	}
	if durationSeconds <= 0 || stepSeconds <= 0 {
		return data, fmt.Errorf("duration and step must be greater than zero")
	}
	if stepSeconds > durationSeconds {
		return data, fmt.Errorf("step must not be greater than duration")
	}

	start := time.Now().Unix()
	itemCount := durationSeconds / stepSeconds
	fn := func(x float64) float64 { return 1 }

	switch data.Type {
	case "linear":
		firstValue, err := strconv.ParseFloat(data.First, 64)
		if err != nil {
			return data, fmt.Errorf("invalid first value")
		}
		lastValue, err := strconv.ParseFloat(data.Last, 64)
		if err != nil {
			return data, fmt.Errorf("invalid last value")
		}
		fn = GetLinear(firstValue, lastValue, start, durationSeconds)
		data.SampleTitle = "linear series preview"
	case "cos":
		minValue, err := strconv.ParseFloat(data.Min, 64)
		if err != nil {
			return data, fmt.Errorf("invalid min value")
		}
		maxValue, err := strconv.ParseFloat(data.Max, 64)
		if err != nil {
			return data, fmt.Errorf("invalid max value")
		}
		fn = GetCosinus(minValue, maxValue, data.Period)
		data.SampleTitle = "cosine series preview"
	case "log":
		fn = GetLog(start)
		data.SampleTitle = "log series preview"
	case "exp":
		fn = GetExp(start, durationSeconds)
		data.SampleTitle = "exp series preview"
	default:
		return data, fmt.Errorf("unknown type %q", data.Type)
	}

	lines := make([]string, 0, itemCount)
	for i := 0; i < itemCount; i++ {
		ts := start + int64(i*stepSeconds)
		lines = append(lines, fmt.Sprintf("%d %.2f", ts, fn(float64(ts))))
	}
	data.Preview = strings.Join(lines, "\n")
	return data, nil
}

func secondsFromDuration(value string) (seconds int, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("invalid duration %q", value)
		}
	}()
	return GetSeconds(value), nil
}

func defaultQueryValue(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}