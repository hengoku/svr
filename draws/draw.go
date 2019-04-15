package draws

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	"os"
)

func DrawWith(filename string, series ...chart.Series) error {
	if err := os.Mkdir("img", os.ModePerm); err != nil {
		return err
	}
	filename = "img/" + filename
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},

		Series: series,
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := buffer.WriteTo(f); err != nil {
		return err
	}

	return nil
}
