all: sass_files
	
sass_files:
	sass login.scss ../server_root/CSS/login.css
	sass not_logged.scss ../server_root/CSS/not_logged.css
	sass login_form.scss ../server_root/CSS/login_form.css
	- cd ../server_root/CSS/ && rm -rf *.css.map

clean:
	-cd ../server_root/CSS/ && rm -rf *.css
	-rm -rf *.css.map
	-rm -rf *~
	-rm -rf .sass-cache
