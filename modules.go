package summer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-wrapper"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
)

type (
	// Func is alias for map[string]func(c *gin.Context)
	Func map[string]func(c *gin.Context)
	// WebFunc is alias for map[string]func(c *websocket.Conn)
	WebFunc map[string]func(ws *websocket.Conn)

	//Module struct
	Module struct {
		Panel
		Collection *mgo.Collection
		Settings   *ModuleSettings
	}

	//ModuleSettings struct
	ModuleSettings struct {
		Name             string
		Menu             *Menu
		MenuTitle        string
		MenuOrder        int
		PageRouteName    string
		AjaxRouteName    string
		SocketsRouteName string
		Title            string
		CollectionName   string
		TemplateName     string
		AllowGroups      []string
		AllowRoles       []string
		AllowUsers       []uint64
		Ajax             Func
		Websockets       WebFunc
	}

	// Simple module interface
	Simple interface {
		Init(settings *ModuleSettings, panel *Panel)
		Page(c *gin.Context)
		Ajax(c *gin.Context)
		Websockets(c *gin.Context)
		GetSettings() *ModuleSettings
	}
)

var (
	modulesList   = map[string]Simple{}
	modulesListMu = sync.Mutex{}
)

// Ajax  is default module's ajax method
func (m *Module) Ajax(c *gin.Context) {
	found := false
	for ajaxRoute, ajaxFunc := range m.Settings.Ajax {
		if strings.ToLower(c.Param("method")) == ajaxRoute {
			ajaxFunc(c)
			found = true
			break
		}
	}

	if !found {
		c.String(400, `Method not found in module "`+m.Settings.Name+`"!`)
	}
}

// Websockets  is default module's websockets method
func (m *Module) Websockets(c *gin.Context) {
	fmt.Println("...")
	fmt.Println(strings.ToLower(c.Param("method")))
	fmt.Println(m.Settings.Websockets)
	found := false
	for websocketsRoute, websocketsFunc := range m.Settings.Websockets {
		fmt.Println("...", strings.ToLower(c.Param("method")), websocketsRoute)
		if strings.ToLower(c.Param("method")) == websocketsRoute {
			handler := websocket.Handler(websocketsFunc)
			handler.ServeHTTP(c.Writer, c.Request)
			found = true
			break
		}
	}

	if !found {
		c.String(400, `Method not found in module "`+m.Settings.Name+`"!`)
	}
}

// Page is default module's page rendering method
func (m *Module) Page(c *gin.Context) {
	c.HTML(200, m.Settings.TemplateName+".html", gin.H{
		"title": m.Settings.Title,
		"user":  c.MustGet("user"),
	})
}

// Init is default module's initial method
func (m *Module) Init(settings *ModuleSettings, panel *Panel) {
	m.Settings = settings
	m.Panel = *panel
	if m.Collection == nil {
		m.Collection = mongo.DB(panel.DBName).C(settings.CollectionName)
	}
}

// GetSettings needs for correct settings getting from module struct
func (m *Module) GetSettings() *ModuleSettings {
	return m.Settings
}