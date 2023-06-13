package gocache

import (
	"MiniArch/gocache/consistenthash"
	"MiniArch/gocache/gcachepb"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const defaultBasePath = "/gCache/"
const defaultReplicas = 50

var defaultTransport = http.DefaultTransport

func NewHTTPPool(self string) *HTTPPool {
	server := NewHTTPPoolOption(self, nil)
	http.Handle(server.config.BasePath, server)
	return server
}
func NewHTTPPoolOption(self string, config *HTTPPoolOptions) *HTTPPool {
	server := &HTTPPool{
		self:        self,
		httpGetters: make(map[string]*httpGetter),
		ctx:         context.Background(),
	}
	if config != nil {
		server.config = config
	} else {
		server.config = new(HTTPPoolOptions)
	}
	if server.config.BasePath == "" {
		server.config.BasePath = defaultBasePath
	}
	if server.config.Replicas == 0 {
		server.config.Replicas = defaultReplicas
	}
	server.transport = defaultTransport
	server.hashRing = consistenthash.NewConsistentHash(server.config.Replicas, nil)
	return server
}

type HTTPPool struct {
	mu sync.Mutex

	ctx context.Context

	self string

	transport http.RoundTripper

	hashRing *consistenthash.ConsistentHash

	httpGetters map[string]*httpGetter

	config *HTTPPoolOptions
}

type HTTPPoolOptions struct {
	BasePath string

	HashFunc consistenthash.HashFunc

	Replicas int
}

func (h *HTTPPool) SetTransportFunc(fn http.RoundTripper) error {
	if h == nil {
		return errors.New("you should init HTTPPool before call SetTransportFunc")
	}
	h.transport = fn
	return nil
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.hashRing == nil {
		h.hashRing = consistenthash.NewConsistentHash(h.config.Replicas, h.config.HashFunc)
	}
	h.hashRing.Add(peers...)
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.config.BasePath, transport: h.transport}
	}
}

// ServeHTTP implements interface Handler
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("accept message from" + r.URL.String())

	// parse request to get cache
	// except that request's url is basePath/groupName/key
	if !strings.HasPrefix(r.URL.Path, h.config.BasePath) {
		http.Error(w, "serving unexpected path", http.StatusBadRequest)
		return
	}
	// 抛弃basePath后，解析后面的/groupName/key
	parts := strings.SplitN(r.URL.Path[len(h.config.BasePath):], "/", 2)
	if len(parts) != 2 {

	}
	// parts[0]=groupName part[1]=key
	groupName := parts[0]
	key := parts[1]

	var err error
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "group which is queried is not exist", http.StatusBadRequest)
	}
	value, err := group.Get(h.ctx, key)
	if err != nil {
		http.Error(w, "group get value occur error", http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(&gcachepb.GetResponse{Value: value.String()})
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *HTTPPool) Run() error {
	log.Println("server is running at " + h.self)
	return http.ListenAndServe(h.self[7:], h)
}

func (h *HTTPPool) PickPromoteGetter(key string) (PromoteGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.hashRing == nil || h.hashRing.IsEmpty() {
		return nil, false
	}
	groupAddr := h.hashRing.Get(key)
	if groupAddr == h.self {
		return nil, false
	}
	getter := h.httpGetters[groupAddr]
	if getter == nil {
		return nil, false
	}
	return getter, true
}

type httpGetter struct {
	baseURL   string
	transport http.RoundTripper
}

func (h *httpGetter) GetFromPromote(ctx context.Context, req *gcachepb.GetRequest, resp *gcachepb.GetResponse) error {
	url := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		req.Group,
		req.Key)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.WithContext(ctx)
	log.Println("send request to " + request.URL.String())
	response, err := h.transport.RoundTrip(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("server return bad status code" + strconv.Itoa(response.StatusCode))
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return err
	}
	return nil
}
