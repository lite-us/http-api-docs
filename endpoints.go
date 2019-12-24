// Package docs can be used to gather go-ipfs commands and automatically
// generate documentation or tests.
package docs

import (
	"fmt"
	"sort"

	jsondoc "github.com/Stebalien/go-json-doc"
	cid "github.com/ipfs/go-cid"
	//config "github.com/ipfs/go-ipfs"
	//config "github.com/TRON-US/go-btfs"
	//cmds "github.com/ipfs/go-ipfs-cmds"
	cmds "github.com/TRON-US/go-btfs-cmds"
	//corecmds "github.com/ipfs/go-ipfs/core/commands"
	corecmds "github.com/TRON-US/go-btfs/core/commands"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var JsondocGlossary = jsondoc.NewGlossary().
	WithSchema(new(cid.Cid), jsondoc.Object{"/": "<cid-string>"}).
	WithName(new(multiaddr.Multiaddr), "multiaddr-string").
	WithName(new(peer.ID), "peer-id").
	WithSchema(new(peerstore.PeerInfo),
		jsondoc.Object{"ID": "peer-id", "Addrs": []string{"<multiaddr-string>"}})

// A map of single endpoints to be skipped (subcommands are processed though).
var IgnoreEndpoints = map[string]bool{}

// How much to indent when generating the response schemas
const IndentLevel = 4

// Failsafe when traversing objects containing objects of the same type
const MaxIndent = 20

// Endpoint defines an IPFS RPC API endpoint.
type Endpoint struct {
	Name        string
	Arguments   []*Argument
	Options     []*Argument
	Description string
	Response    string
	Group       string
}

// Argument defines an IPFS RPC API endpoint argument.
type Argument struct {
	Name        string
	Description string
	Type        string
	Required    bool
	Default     string
}

type sorter []*Endpoint

func (a sorter) Len() int           { return len(a) }
func (a sorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

const APIPrefix = "/api/vXX"

// AllEndpoints gathers all the endpoints from go-ipfs.
func AllEndpoints() []*Endpoint {
	return Endpoints(APIPrefix, corecmds.Root)
}

func IPFSVersion() string {
	return "config.CurrentVersionNumber"
}

// Endpoints receives a name and a go-ipfs command and returns the endpoints it
// defines] (sorted). It does this by recursively gathering endpoints defined by
// subcommands. Thus, calling it with the core command Root generates all
// the endpoints.
func Endpoints(name string, cmd *cmds.Command) (endpoints []*Endpoint) {
	var arguments []*Argument
	var options []*Argument

	ignore := cmd.Run == nil || IgnoreEndpoints[name]
	if !ignore { // Extract arguments, options...
		for _, arg := range cmd.Arguments {
			argType := "string"
			if arg.Type == cmds.ArgFile {
				argType = "file"
			}
			arguments = append(arguments, &Argument{
				Name:        arg.Name,
				Type:        argType,
				Required:    arg.Required,
				Description: arg.Description,
			})
		}

		for _, opt := range cmd.Options {
			def := fmt.Sprint(opt.Default())
			if def == "<nil>" {
				def = ""
			}
			options = append(options, &Argument{
				Name:        opt.Names()[0],
				Type:        opt.Type().String(),
				Description: opt.Description(),
				Default:     def,
			})
		}

		res := buildResponse(cmd.Type)

		endpoints = []*Endpoint{
			&Endpoint{
				Name:        name,
				Description: cmd.Helptext.Tagline,
				Arguments:   arguments,
				Options:     options,
				Response:    res,
			},
		}
	}

	for n, cmd := range cmd.Subcommands {
		endpoints = append(endpoints,
			Endpoints(fmt.Sprintf("%s/%s", name, n), cmd)...)
	}
	sort.Sort(sorter(endpoints))
	return endpoints
}

func buildResponse(res interface{}) string {
	// Commands with a nil type return text. This is a bad thing.
	if res == nil {
		return "This endpoint returns a `text/plain` response body."
	}
	desc, err := JsondocGlossary.Describe(res)
	if err != nil {
		panic(err)
	}
	return desc
}
