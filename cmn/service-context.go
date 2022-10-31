package cmn

import (
	"errors"
	"fmt"
	"sync"
)

type ModuleAuthor struct {
	Name  string `json:"name"`
	Tel   string `json:"tel"`
	Email string `json:"email"`
	Addi  string `json:"addi"`
}
type ServeEndPoint struct {
	Path      string
	Name      string
	Developer *ModuleAuthor
}

var (
	Services     = make(map[string]*ServeEndPoint)
	serviceMutex sync.Mutex
)

func AddService(ep *ServeEndPoint) (err error) {
	for {
		if ep == nil {
			err = errors.New("ep is nil")
			break
		}

		if ep.Path == "" {
			err = errors.New("ep.path empty")
			break
		}

		//if ep.PathPattern == "" {
		//	ep.PathPattern = fmt.Sprintf(`(?i)^%s(/.*)?$`, ep.Path)
		//}
		//ep.PathMatcher = regexp.MustCompile(ep.PathPattern)
		//
		//if ep.IsFileServe {
		//	if ep.DocRoot == "" {
		//		err = errors.New("must specify docRoot when ep.isFileServe equal true")
		//		break
		//	}
		//
		//	if ep.Fn == nil {
		//		ep.Fn = WebFS
		//	}
		//} else {
		//	if ep.Fn == nil {
		//		err = errors.New("must specify fn when ep.isFileServe equal false")
		//		break
		//	}
		//
		//	if !rIsAPI.MatchString(ep.Path) {
		//		ep.Path = strings.ReplaceAll("/api/"+ep.Path, "//", "/")
		//	}
		//}

		if ep.Name == "" {
			err = errors.New("must specify apiName")
			break
		}

		_, ok := Services[ep.Path]
		if ok {
			err = errors.New(fmt.Sprintf("%s[%s] already exists", ep.Path, ep.Name))
		}
		break
	}

	if err != nil {
		//z.Error(err.Error())
		return
	}

	//z.Info(ep.Name + " added")

	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	Services[ep.Path] = ep
	return
}
