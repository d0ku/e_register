package core

//Page is basic abstraction over HTTP webpage.
type Page struct {
	cssFiles []string

	title       string
	author      string
	viewport    string
	description string
	language    string
	charset     string

	body string
}

//GetPage returns empty (not filled) page template to fill by yourself.
//Reset CSS is added to CSS files by default.
func GetPage() *Page {
	var defaultLanguage = "en"
	var defaultCharset = "utf-8"
	var defaultViewport = "width=device-width, initial-scale=1.0"

	page := &Page{language: defaultLanguage, charset: defaultCharset, viewport: defaultViewport}

	page.AddCSSFile("reset.css")

	return page
}

//SetCharSet sets value of meta tag charset.
func (page *Page) SetCharSet(charset string) {
	page.charset = charset
}

//SetTitle sets value of meta tag title.
func (page *Page) SetTitle(title string) {
	page.title = title
}

//SetAuthor sets value of meta tag author.
func (page *Page) SetAuthor(author string) {
	page.author = author
}

//SetViewPort sets value of meta tag viewport.
func (page *Page) SetViewPort(viewport string) {
	page.viewport = viewport
}

//SetDescription sets value of meta tag description.
func (page *Page) SetDescription(description string) {
	page.description = description
}

//SetLanguage sets value written in top html tag.
func (page *Page) SetLanguage(language string) {
	page.language = language
}

//AddCSSFile adds provided path to css files supported in generator.
func (page *Page) AddCSSFile(name string) {
	page.cssFiles = append(page.cssFiles, "CSS/"+name)
}

//SetBody sets value of html body.
func (page *Page) SetBody(body string) {
	page.body = body
}

//GetHTMLString return source code of generated page.
func (page *Page) GetHTMLString() string {
	var html string

	html = "<!DOCTYPE html>\n"

	html += "<html lang=\"" + page.language + "\">\n"

	html += "<head>\n"

	html += "  <meta charset=\"" + page.charset + "\">\n"
	html += "  <title>" + page.title + "</title>\n"
	html += "  <meta name=\"description\" content=\"" + page.description + "\">\n"
	html += "  <meta name=\"author\" content=\"" + page.author + "\">\n"
	html += "  <meta name=\"viewport\" content=\"" + page.viewport + "\">\n"

	for _, cssFile := range page.cssFiles {
		html += "  <link rel=\"stylesheet\" type=\"text/css\" href=\"" + cssFile + "\">\n"
	}

	html += "</head>\n"

	html += "<body>\n"
	html += page.body + "\n"

	html += "</body>\n"

	html += "</html>\n"

	return html
}
