package vo

import (
	"interastral-peace.com/alnitak/internal/domain/model"
)

type MenuResp struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Desc      string     `json:"desc"`
	Sort      uint       `json:"sort"`
	Component string     `json:"component"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuResp `json:"children"`
	ParentId  uint       `json:"parentId"`
}

type MenuMeta struct {
	Title     string `json:"title"`
	KeepAlive bool   `json:"keepAlive"`
	Icon      string `json:"icon"`
	Hidden    bool   `json:"hidden"`
}

func MenuListToMenuTree(menus []model.Menu, parentId uint) []MenuResp {
	// newMenus := make([]MenuResp, len(menus))
	// for i := 0; i < len(menus); i++ {
	// 	newMenus[i] = MenuToMenuResp(menus[i])
	// }

	// return newMenus

	res := make([]MenuResp, 0)
	for _, v := range menus {
		if v.ParentId == parentId {
			newMenus := MenuToMenuResp(v)
			newMenus.Children = MenuListToMenuTree(menus, v.ID)
			res = append(res, newMenus)
		}
	}
	return res
}

func MenuListToMenuResp(menus []model.Menu) []MenuResp {
	newMenus := make([]MenuResp, len(menus))
	for i := 0; i < len(menus); i++ {
		newMenus[i] = MenuToMenuResp(menus[i])
	}

	return newMenus
}

func MenuToMenuResp(menu model.Menu) MenuResp {
	newMenus := MenuResp{
		ID:        menu.ID,
		Name:      menu.Name,
		Desc:      menu.Desc,
		Path:      menu.Path,
		Component: menu.Component,
		Sort:      menu.Sort,
		ParentId:  menu.ParentId,
		Meta: MenuMeta{
			Title:     menu.Title,
			KeepAlive: menu.KeepAlive,
			Icon:      menu.Icon,
			Hidden:    menu.Hidden,
		},
	}

	for i := 0; i < len(menu.Children); i++ {
		newMenus.Children = append(newMenus.Children, MenuToMenuResp(menu.Children[i]))
	}

	return newMenus
}
