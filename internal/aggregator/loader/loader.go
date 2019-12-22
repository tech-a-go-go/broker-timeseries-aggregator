package loader

import (
	"bufio"
	"compress/gzip"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/log"
)

var logger = log.GetLogger()

// DataLoader ファイルに保存されたデータをロードし一行ずつchannelに送信します
type DataLoader struct {
	dataPath string      // ブローカーのデータが保存されているディレクトリのパス
	dataCh   chan []byte // データを一行ずつ送信するchannel
	endCh    chan interface{}
}

func NewDataLoader(dataPath string) *DataLoader {
	loader := &DataLoader{
		dataPath: dataPath,
		dataCh:   make(chan []byte),
		endCh:    make(chan interface{}),
	}
	return loader
}

func (d *DataLoader) GetDataCh() chan []byte {
	return d.dataCh
}

func (d *DataLoader) GetEndCh() chan interface{} {
	return d.endCh
}

func (d *DataLoader) Load() error {
	filenames, err := d.GetDataFilenames()
	if err != nil {
		return errors.Wrap(err, "Failed to get filenames")
	}
	for _, filename := range filenames {
		// fmt.Println(filename, "--------------------------------------------------")
		err := d.read(filename)
		if err != nil {
			close(d.dataCh)
			return errors.Wrap(err, "Failed to read file: "+filename)
		}
	}
	close(d.endCh)
	close(d.dataCh)
	return nil
}

func (d *DataLoader) read(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "Failed to oepn open file")
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		return errors.Wrap(err, "Failed to create a gzip reader")
	}

	scanner := bufio.NewScanner(gz)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// scanner.Bytes() は戻り値の []byte のメモリ領域を毎回再利用しているので次にscanner.Scan()された時には
		// メモリ領域が新しいデータで上書きされているので呼び出し元でキャッシュなどしないように注意すること
		d.dataCh <- scanner.Bytes()
	}
	gz.Close()
	file.Close()
	return nil
}

func (d *DataLoader) GetDataFilenames() ([]string, error) {
	return filepath.Glob(d.dataPath)
}
