package cms

import (
  "github.com/Lunkov/lib-ui"
)

type CMS struct {
  Conf  ConfigInfo
  U    *ui.UI
}

func New() *CMS {
  return &CMS{}
}

func (c *CMS)_InitUI() {
  c.U = ui.NewUI(c.Conf.UI.PathTemplates, &c.Conf.UI.CacheForms, &c.Conf.UI.CacheViews, &c.Conf.UI.CachePages, &c.Conf.UI.CacheRenders) 
  c.U.Init(c.Conf.ConfigPath, false, false)
}
