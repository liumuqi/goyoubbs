package controller

import (
	"github.com/ego008/goyoubbs/model"
	"github.com/ego008/youdb"
	"goji.io/pat"
	"net/http"
	"strconv"
)

func (h *BaseHandler) CategoryDetail(w http.ResponseWriter, r *http.Request) {
	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if len(score) > 0 {
		_, err := strconv.ParseUint(score, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"score type err"}`))
			return
		}
	}

	cid := pat.Param(r, "cid")
	_, err := strconv.Atoi(cid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

	cmd := "zrscan"
	if btn == "prev" {
		cmd = "zscan"
	}

	db := h.App.Db
	cobj, err := model.CategoryGetById(db, cid)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	if cobj.Hidden && currentUser.Flag < 99 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}
	cobj.Articles = db.Zget("category_article_num", youdb.I2b(cobj.Id)).Uint64()
	pageInfo := model.ArticleList(db, cmd, "category_article_timeline:"+cid, key, score, h.App.Cf.Site.HomeShowNum)

	type pageData struct {
		PageData
		Cobj     model.Category
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = cobj.Name + " - " + h.App.Cf.Site.Name
	evn.Keywords = cobj.Name
	evn.Description = cobj.About
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, h.App.Cf.Site.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, h.App.Cf.Site.CategoryShowNum)

	evn.Cobj = cobj
	evn.PageInfo = pageInfo

	h.Render(w, tpl, evn, "layout.html", "category.html")
}
