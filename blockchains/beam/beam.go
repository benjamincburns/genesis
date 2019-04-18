package beam

import (
	db "../../db"
	ssh "../../ssh"
	state "../../state"
	util "../../util"
	helpers "../helpers"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

var conf *util.Config

func init() {
	conf = util.GetConfig()
}

const port int = 10000

func Build(details *db.DeploymentDetails, servers []db.Server, clients []*ssh.Client,
	buildState *state.BuildState) ([]string, error) {

	beamConf, err := NewConf(details.Params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	buildState.SetBuildSteps(0 + (details.Nodes * 4))

	buildState.SetBuildStage("Setting up the wallets")
	/**Set up wallets**/
	ownerKeys := make([]string, details.Nodes)
	secretMinerKeys := make([]string, details.Nodes)
	mux := sync.Mutex{}
	// walletIDs := []string{}
	err = helpers.AllNodeExecCon(servers, buildState, func(serverNum int, localNodeNum int, absoluteNodeNum int) error {

		clients[serverNum].DockerExec(localNodeNum, "beam-wallet --command init --pass password") //ign err

		res1, _ := clients[serverNum].DockerExec(localNodeNum, "beam-wallet --command export_owner_key --pass password") //ign err

		buildState.IncrementBuildProgress()

		re := regexp.MustCompile(`(?m)^Owner([A-z|0-9|\s|\:|\/|\+|\=])*$`)
		ownKLine := re.FindAllString(res1, -1)[0]

		mux.Lock()
		ownerKeys[absoluteNodeNum] = strings.Split(ownKLine, " ")[3]
		mux.Unlock()

		res2, _ := clients[serverNum].DockerExec(localNodeNum, "beam-wallet --command export_miner_key --subkey=1 --pass password") //ign err

		re = regexp.MustCompile(`(?m)^Secret([A-z|0-9|\s|\:|\/|\+|\=])*$`)
		secMLine := re.FindAllString(res2, -1)[0]

		mux.Lock()
		secretMinerKeys[absoluteNodeNum] = strings.Split(secMLine, " ")[3]
		mux.Unlock()

		buildState.IncrementBuildProgress()
		return nil
	})

	ips := []string{}
	for _, server := range servers {
		for _, ip := range server.Ips {
			ips = append(ips, ip)
		}
	}
	buildState.SetBuildStage("Creating node configuration files")
	/**Create node config files**/

	err = helpers.CreateConfigs(servers, clients, buildState, "/beam/beam-node.cfg",
		func(serverNum int, localNodeNum int, absoluteNodeNum int) ([]byte, error) {
			beam_node_config, err := makeNodeConfig(beamConf, ownerKeys[absoluteNodeNum],
				secretMinerKeys[absoluteNodeNum], details, absoluteNodeNum)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			for _, ip := range append(ips[:absoluteNodeNum], ips[absoluteNodeNum+1:]...) {
				beam_node_config += fmt.Sprintf("peer=%s:%d\n", ip, port)
			}
			return []byte(beam_node_config), nil
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = helpers.CreateConfigs(servers, clients, buildState, "/beam/beam-wallet.cfg",
		func(serverNum int, localNodeNum int, absoluteNodeNum int) ([]byte, error) {
			beam_wallet_config := []string{
				"# Emission.Value0=800000000",
				"# Emission.Drop0=525600",
				"# Emission.Drop1=2102400",
				"Maturity.Coinbase=1",
				"# Maturity.Std=0",
				"# MaxBodySize=0x100000",
				"DA.Target_s=1",
				"# DA.MaxAhead_s=900",
				"# DA.WindowWork=120",
				"# DA.WindowMedian0=25",
				"# DA.WindowMedian1=7",
				"DA.Difficulty0=100",
				"# AllowPublicUtxos=0",
				"# FakePoW=0",
			}
			return []byte(util.CombineConfig(beam_wallet_config)), nil
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	buildState.SetBuildStage("Starting beam")
	err = helpers.AllNodeExecCon(servers, buildState, func(serverNum int, localNodeNum int, absoluteNodeNum int) error {
		defer buildState.IncrementBuildProgress()
		miningFlag := ""
		if absoluteNodeNum >= int(beamConf.Validators) {
			miningFlag = " --mining_threads 1"
		}
		_, err := clients[serverNum].DockerExecd(localNodeNum, fmt.Sprintf("beam-node%s", miningFlag))
		if err != nil {
			log.Println(err)
			return err
		}
		return clients[serverNum].DockerExecdLog(localNodeNum, fmt.Sprintf("beam-wallet --command listen -n 0.0.0.0:%d --pass password", port))
	})

	return nil, err
}

func Add(details *db.DeploymentDetails, servers []db.Server, clients []*ssh.Client,
	newNodes map[int][]string, buildState *state.BuildState) ([]string, error) {
	return nil, nil
}
