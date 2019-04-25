package beam

import (
	db "../../db"
	util "../../util"
	helpers "../helpers"
	"github.com/Whiteblock/mustache"
	"io/ioutil"
)

type BeamConf struct {
	Validators int64 `json:"validators"`
	TxNodes    int64 `json:"txNodes"`
	NilNodes   int64 `json:"nilNodes"`
}

func NewConf(data map[string]interface{}) (*BeamConf, error) {
	out := new(BeamConf)

	err := util.GetJSONInt64(data, "validators", &out.Validators)
	if err != nil {
		return nil, err
	}

	err = util.GetJSONInt64(data, "txNodes", &out.TxNodes)
	if err != nil {
		return nil, err
	}

	err = util.GetJSONInt64(data, "nilNodes", &out.NilNodes)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func GetParams() string {
	dat, err := ioutil.ReadFile("./resources/beam/params.json")
	if err != nil {
		panic(err) //Missing required files is a fatal error
	}
	return string(dat)
}

func GetDefaults() string {
	dat, err := ioutil.ReadFile("./resources/beam/defaults.json")
	if err != nil {
		panic(err) //Missing required files is a fatal error
	}
	return string(dat)
}

func GetServices() []util.Service {
	return nil
}

func makeNodeConfig(bconf *BeamConf, keyOwner string, keyMine string, details *db.DeploymentDetails, node int) (string, error) {

	filler := util.ConvertToStringMap(map[string]interface{}{
		"keyOwner": keyOwner,
		"keyMine":  keyMine,
	})
	dat, err := helpers.GetBlockchainConfig("beam", node, "beam-node.cfg.mustache", details)
	if err != nil {
		return "", err
	}
	data, err := mustache.Render(string(dat), filler)
	return data, err
}
