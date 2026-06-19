package base

import "github.com/Se7enSe7enSe7en/tenant-manager/internal/web/component/sidebar"

func GetSidebarSections() []sidebar.SidebarSection {
	return []sidebar.SidebarSection{
		{
			Title: "Components",
			Items: []sidebar.SidebarItem{
				{Title: "Dashboard", Href: "/dashboard"},
				{Title: "Properties", Href: "/property"},
				// {Title: "Breadcrumb", Href: "/components/breadcrumb"},
			},
		},
	}
}

func getCurrentPath(currentPage string) string {
	switch currentPage {
	case "docs":
		return "/docs"
	case "components":
		return "/components"
	default:
		return "/"
	}
}
