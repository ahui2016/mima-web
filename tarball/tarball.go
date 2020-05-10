package tarball

// 参考 https://gist.github.com/maximilien/328c9ac19ab0a158a8df

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha512"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ahui2016/mima-go/util"
)

// Create 把 filePaths 里的全部文件打包压缩, 新建文件 tarballFilePath.
func Create(tarballFilePath string, filePaths []string) error {
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return err
	}
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	var allErrors []error
	for _, filePath := range filePaths {
		if err := addFileToTar(filePath, tarWriter); err != nil {
			allErrors = append(allErrors, err)
			break
		}
	}
	allErrors = append(allErrors, tarWriter.Close(), gzipWriter.Close(), file.Close())
	return util.WrapErrors(allErrors...)
}

func addFileToTar(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name:    filepath.Base(filePath),
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}
	return nil
}

// Reader 用来帮助读取 tarball 内容.
// 主要是为了更方便地关闭资源.
type Reader struct {
	file       *os.File
	gzipReader *gzip.Reader
	tarReader  *tar.Reader
}

// NewReader 返回一个新的 *tarball.Reader.
// 包括打开文件以及相关 reader.
func NewReader(filePath string) (*Reader, error) {
	tr := new(Reader)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	tr.file = file

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, util.WrapErrors(file.Close(), err)
	}
	tr.gzipReader = gzipReader
	tr.tarReader = tar.NewReader(gzipReader)
	return tr, nil
}

// Sha512 返回一个 tarball 里面全部文件的 SHA512 checksum.
// 只适用于 tarball 里只有文件, 没有文件夹的情况.
func (tr Reader) Sha512() (checksums [][]byte, err error) {
	for {
		_, err := tr.tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 错误处理强迫症, 不想漏掉任何一个错误信息.
			return nil, util.WrapErrors(tr.Close(), err)
		}

		var data []byte
		data, err = ioutil.ReadAll(tr.tarReader)
		if err != nil {
			return nil, util.WrapErrors(tr.Close(), err)
		}
		sum := sha512.Sum512(data)
		checksums = append(checksums, sum[:])
	}
	return
}

// Close 依次关闭 tr 里的 各个 reader, 并把错误合并后返回.
func (tr Reader) Close() error {
	return util.WrapErrors(tr.gzipReader.Close(), tr.file.Close())
}
