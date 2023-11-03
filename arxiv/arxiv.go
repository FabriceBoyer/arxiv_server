package arxiv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fabriceboyer/common_go_utils/utils"
)

// TODO use bzip2 version instead
const dataFileName = "arxiv-metadata-oai.json"
const indexFileName = "arxiv-metadata-index.txt"
const indexSeparator = ":"

type ArxivMetadataVersion struct {
	Version string `json:"version"`
	Created string `json:"created"`
}
type ArxivMetadata struct {
	Id            string                 `json:"id"`
	Submitter     string                 `json:"submitter"`
	Authors       string                 `json:"authors"`
	Title         string                 `json:"title"`
	Comments      string                 `json:"comments"`
	JournalRef    string                 `json:"journal-ref"`
	Doi           string                 `json:"doi"`
	Abstract      string                 `json:"abstract"`
	ReportNo      string                 `json:"report-no"`
	Categories    string                 `json:"categories"`
	Versions      []ArxivMetadataVersion `json:"versions"`
	UpdateDate    string                 `json:"update_date"`
	AuthorsParsed [][]string             `json:"authors_parsed"`
}

func (m *ArxivMetadata) String() string {
	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprint("Error:", err)
	}

	return string(jsonBytes)
}

type ArxivMetadataIndex struct {
	seek int64
	id   string
}

type ArxivMetadataManager struct {
	Root_path string
	index     map[string]int64
}

func (mgr *ArxivMetadataManager) InitializeManager() error {
	regenerate := true
	data_info, err := os.Stat(path.Join(mgr.Root_path, dataFileName))
	if err != nil {
		return err
	}
	index_info, err := os.Stat(path.Join(mgr.Root_path, indexFileName))
	if err == nil {
		regenerate = data_info.ModTime().After(index_info.ModTime())
	}
	if regenerate {
		fmt.Print("Index file missing or out of date, regenerating...\n")
		err = mgr.generateArxivMetadataIndex()
		if err != nil {
			return err
		}
	}
	// else {
	// 	fmt.Print("Reusing existing file index\n")
	// }
	err = mgr.readArxivMetadataIndex()
	if err != nil {
		return err
	}
	return nil
}

func (mgr *ArxivMetadataManager) generateArxivMetadataIndex() error {
	fmt.Print("Generating index file\n")
	index := []ArxivMetadataIndex{}

	f, err := os.Open(path.Join(mgr.Root_path, dataFileName))
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)

	i := 0
	for d.More() {
		seek := d.InputOffset()
		elm := &ArxivMetadata{}
		err = d.Decode(elm)
		if err != nil {
			return err
		}
		index = append(index, ArxivMetadataIndex{seek, elm.Id})
		if i%1e5 == 0 {
			fmt.Printf("%s elements parsed\n", humanize.SI(float64(i), ""))
		}
		i++
	}

	fmt.Print("Writing index file\n")
	index_file, err := os.Create(path.Join(mgr.Root_path, indexFileName))
	if err != nil {
		return err
	}
	defer index_file.Close()

	for _, elm := range index {
		index_file.WriteString(fmt.Sprint(elm.seek) + indexSeparator + elm.id + "\n")
	}

	return nil
}

func (mgr *ArxivMetadataManager) GetIndexedArxivMetadata(id string) (*ArxivMetadata, error) {
	seek, found := mgr.index[id]
	if !found {
		return nil, fmt.Errorf("index %v not found", id)
	}

	res, err := mgr.getArxivMetadaFromSeekPosition(seek)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (mgr *ArxivMetadataManager) getArxivMetadaFromSeekPosition(seek int64) (*ArxivMetadata, error) {
	f, err := os.Open(path.Join(mgr.Root_path, dataFileName)) // TODO don't do it on every call, keep it in cache
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(seek, 0)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	elm := &ArxivMetadata{}
	err = d.Decode(elm)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (mgr *ArxivMetadataManager) readArxivMetadataIndex() error {
	//fmt.Print("Reading index file\n")
	mgr.index = make(map[string]int64)
	readFile, err := os.Open(path.Join(mgr.Root_path, indexFileName))
	if err != nil {
		return err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		parts := strings.Split(line, indexSeparator)
		if len(parts) != 2 {
			return fmt.Errorf("expected exactly 2 parts in index, got: [%s]", parts)
		} else {
			seek, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return err
			}
			id := parts[1]
			mgr.index[id] = seek
		}
	}

	return nil
}

type filter func(*ArxivMetadata) bool

func (mgr *ArxivMetadataManager) SearchArxivMetadata(limit int, filt filter) ([]*ArxivMetadata, int, error) {

	res := []*ArxivMetadata{}
	i := 0
	start := time.Now()

	f, err := os.Open(path.Join(mgr.Root_path, dataFileName)) // TODO don't do it on every call, keep it in cache
	if err != nil {
		return res, i, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return res, i, err
	}

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)

	for d.More() {
		elm := &ArxivMetadata{}
		err = d.Decode(elm)
		if err != nil {
			return res, i, err
		}

		if filt != nil && filt(elm) {
			res = append(res, elm)
		}

		// val, _ := json.Marshal(elm)
		// fmt.Printf("%v \n", string(val))
		if i%1e5 == 0 {
			fmt.Printf("%s elements parsed\n", humanize.SI(float64(i), ""))
		}

		if limit > 0 && i >= limit {
			break
		}
		i++
	}

	elapsed := time.Since(start)

	fmt.Printf("Total of [%v] object read.\n", i)
	fmt.Printf("The [%s] is %s long\n", dataFileName, humanize.Bytes(uint64(fi.Size())))
	fmt.Printf("Parsing the file took [%v]\n", elapsed)

	return res, i, nil
}

func (mgr *ArxivMetadataManager) getRandomArxivIds(count int) ([]string, error) {

	keys := mgr.GetMapKeys()
	res := utils.GetRandomSample(keys, count)

	return res, nil

}

func (mgr *ArxivMetadataManager) GetMapKeys() []string {
	m := &mgr.index
	keys := make([]string, 0, len(*m))
	for key := range *m {
		keys = append(keys, key)
	}
	return keys
}
