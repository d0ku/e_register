all: sass_files
	
sass_files:
	sass login.scss ../server_root/CSS/login.css --style compressed
	sass not_logged.scss ../server_root/CSS/not_logged.css --style compressed
	sass login_form.scss ../server_root/CSS/login_form.css --style compressed
	sass choose_school.scss ../server_root/CSS/choose_school.css --style compressed
	sass main_teacher.scss ../server_root/CSS/main_teacher.css --style compressed
	- cd ../server_root/CSS/ && rm -rf *.css.map

clean:
	-cd ../server_root/CSS/ && rm -rf *.css
	-rm -rf *.css.map
	-rm -rf *~
	-rm -rf .sass-cache
