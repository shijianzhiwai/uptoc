package engine

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/saltbo/gopkg/fileutil"

	"uptoc/uploader"
)

// Config provides core configuration for the engine.
type Config struct {
	SaveRoot  string   `yaml:"save_root"`
	VisitHost string   `yaml:"visit_host"`
	ForceSync bool     `yaml:"force_sync"`
	Excludes  []string `yaml:"excludes"`
}

// Engine provides the core logic to finish the feature
type Engine struct {
	echo
	conf Config

	uploader uploader.Driver
}

// New returns a new engine.
func New(conf Config, ud uploader.Driver) *Engine {

	// windows path compatible.
	conf.SaveRoot = pathCompatible(conf.SaveRoot)
	for k, v := range conf.Excludes {
		conf.Excludes[k] = pathCompatible(v)
	}

	return &Engine{
		conf:     conf,
		uploader: ud,
	}
}

// TailRun run the core logic with every path.
func (e *Engine) TailRun(paths ...string) {
	for _, path := range paths {
		path = pathCompatible(path)
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatalln(err)
		}

		if stat.IsDir() {
			e.uploadDirectory(path)
			continue
		}

		e.uploadFile(path, path2ObjectKey(filepath.Join(e.conf.SaveRoot, stat.Name())))
	}
}

func (e *Engine) uploadDirectory(dirPath string) {
	objects, err := e.loadLocalObjects(dirPath)
	if err != nil {
		log.Fatalln(err)
	}

	// directory sync
	if e.conf.ForceSync {
		s := NewSyncer(e.uploader)
		if err := s.Sync(objects, e.conf.SaveRoot, e.conf.Excludes); err != nil {
			log.Fatalln(err)
		}
		return
	}

	// directory normal upload
	for _, obj := range objects {
		e.uploadFile(obj.FilePath, obj.Key)
	}
}

func (e *Engine) uploadFile(filePath, object string) {
	if err := e.uploader.Upload(object, filePath); err != nil {
		e.Failed(filePath, err)
		return
	}

	e.Success(e.conf.VisitHost, object)
}

func (e *Engine) loadLocalObjects(dirPath string) ([]uploader.Object, error) {
	dirPath = addDirSuffix(dirPath)

	localObjects := make([]uploader.Object, 0)
	visitor := func(filePath string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() || shouldExclude(dirPath, filePath, e.conf.Excludes) {
			return nil
		}

		localPath := strings.TrimPrefix(filePath, dirPath)
		localObjects = append(localObjects, uploader.Object{
			Key:      path2ObjectKey(filepath.Join(e.conf.SaveRoot, localPath)),
			ETag:     fileutil.MD5Hex(filePath),
			FilePath: filePath,
		})
		return nil
	}

	if err := filepath.Walk(dirPath, visitor); err != nil {
		return nil, err
	}

	return localObjects, nil
}
