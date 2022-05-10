package dirscan

import (
	"chan_psw/pkg/chantools"
	"chan_psw/pkg/metrics"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	Base = "scanDir"
	Size = Base + ".Size"
)

var (
	// Define the metrics we use
	mnames = []string{Size}
)

func recoverFromClosedChan() {
	chantools.RecoverFromClosedChan()
}

// DirScan holds a directry to scan and a Channel to put the filenames onto
type DirScan struct {
	Dir      string
	ScanTime time.Duration
	ch       chan string
	ctx      context.Context
	can      context.CancelFunc
	once     sync.Once
	met      *metrics.Metrics
}

// setNumFiles will set a metric if defined to the number of files in the current directory
func (d *DirScan) setNumFiles(s int) {
	if d.met != nil {
		d.met.Set(Size, s)
	}
}

// scanDir
// This is Blocking on the channel write
func (d *DirScan) scanDir() error {
	// This is to recover if we write to a closed channel, that is not a problem so recover from a panic
	defer recoverFromClosedChan()

	entries, err := os.ReadDir(d.Dir)
	if err != nil {
		return err
	}

	path, err := filepath.Abs(d.Dir)
	if err != nil {
		return err
	}

	d.setNumFiles(len(entries))
	for _, f := range entries {
		// if it is not a directory and is not .prefixed
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			select {
			case d.ch <- filepath.Join(path, f.Name()):
			case <-d.ctx.Done():
				return nil
			}
		}
	}
	return nil
}

// loop until we receive a stop on the run channel
func (d *DirScan) loop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		d.scanDir()
		select {
		case <-time.After(d.ScanTime):
		case <-d.ctx.Done():
			return
		}
	}
}

// -----  Public Methods

// OutputChan returns the output Channel as ReadOnly
func (d *DirScan) OutputChan() <-chan string {
	return d.ch
}

// Close will close the data channel
func (d *DirScan) Close() {
	d.can()
	d.once.Do(func() {
		close(d.ch)
	})
}

// SetMetric set the metrics var and adds our metric names to it
func (d *DirScan) SetMetric(met *metrics.Metrics) (names []string) {
	d.met = met
	d.met.AddMetric(mnames...)

	return mnames
}

// New creates a new dir scanner and starts a scanning loop to send filenames to a channel
// Must pass a WaitGroup it as we create a go routine for the scanner
// As a writter we assume we own the channel we return, we will close it when our Close() is called
func New(wg *sync.WaitGroup, dir string, scantime time.Duration, chanSize int) (*DirScan, error) {
	if wg == nil {
		return nil, errors.New("must provide a valid waitgroup")
	}

	if fileInfo, err := os.Stat(dir); err != nil {
		return nil, errors.New("name not found: " + dir)
	} else {
		if !fileInfo.IsDir() {
			return nil, errors.New("name is not a directory: " + dir)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := DirScan{Dir: dir, ch: make(chan string, chanSize), ScanTime: scantime, ctx: ctx, can: cancel}

	wg.Add(1)
	go d.loop(wg)

	return &d, nil
}
