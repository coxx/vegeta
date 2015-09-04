package vegeta

import (
	"encoding/json"
	"fmt"
	"io"
)

// CSVDumper implements the Reporter interface by reporting a Result as a CSV
// record with six columns. The columns are: unix timestamp in ns since epoch,
// http status code, request latency in ns, bytes out, bytes in, and lastly the error.
type CSVDumper struct{ Result }

// Report partially implements the Reporter interface.
func (d CSVDumper) Report(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d,%d,%d,%d,%d,\"%s\"\n",
		d.Timestamp.UnixNano(),
		d.Code,
		d.Latency.Nanoseconds(),
		d.BytesOut,
		d.BytesIn,
		d.Error,
	)
	return err
}

// JSONDumper implements the Reporter interface by reporting a Result as a JSON
// object.
type JSONDumper struct{ Result }

// Report partially implements the Reporter interface.
func (d JSONDumper) Report(w io.Writer) error {
	return json.NewEncoder(w).Encode(d.Result)
}
