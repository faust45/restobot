;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; Start of xbindkeys guile configuration ;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This configuration is guile based.
;;   http://www.gnu.org/software/guile/guile.html
;; any functions that work in guile will work here.
;; see EXTRA FUNCTIONS:

;; Version: 1.8.6

;; If you edit this file, do not forget to uncomment any lines
;; that you change.
;; The semicolon(;) symbol may be used anywhere for comments.

;; To specify a key, you can use 'xbindkeys --key' or
;; 'xbindkeys --multikey' and put one of the two lines in this file.

;; A list of keys is in /usr/include/X11/keysym.h and in
;; /usr/include/X11/keysymdef.h
;; The XK_ is not needed.

;; List of modifier:
;;   Release, Control, Shift, Mod1 (Alt), Mod2 (NumLock),
;;   Mod3 (CapsLock), Mod4, Mod5 (Scroll).


;; The release modifier is not a standard X modifier, but you can
;; use it if you want to catch release instead of press events

;; By defaults, xbindkeys does not pay attention to modifiers
;; NumLock, CapsLock and ScrollLock.
;; Uncomment the lines below if you want to use them.
;; To dissable them, call the functions with #f


;;;;EXTRA FUNCTIONS: Enable numlock, scrolllock or capslock usage
;;(set-numlock! #t)
;;(set-scrolllock! #t)
;;(set-capslock! #t)

;;;;; Scheme API reference
;;;;
;; Optional modifier state:
;; (set-numlock! #f or #t)
;; (set-scrolllock! #f or #t)
;; (set-capslock! #f or #t)
;; 
;; Shell command key:
;; (xbindkey key "foo-bar-command [args]")
;; (xbindkey '(modifier* key) "foo-bar-command [args]")
;; 
;; Scheme function key:
;; (xbindkey-function key function-name-or-lambda-function)
;; (xbindkey-function '(modifier* key) function-name-or-lambda-function)
;; 
;; Other functions:
;; (remove-xbindkey key)
;; (run-command "foo-bar-command [args]")
;; (grab-all-keys)
;; (ungrab-all-keys)
;; (remove-all-keys)
;; (debug)


;; Examples of commands:


;; set directly keycode (here control + f with my keyboard)
;; (xbindkey '("m:0x4" "c:41") "xterm")

;; specify a mouse button
;; (xbindkey '(control "b:2") "xterm")

;;(xbindkey '(shift mod2 alt s) "xterm -geom 50x20+20+20")

;; set directly keycode (control+alt+mod2 + f with my keyboard)


; download folder
(remove-all-keys)
(xbindkey '(control alt y) "firefox -new-window youtube.com")
;(xbindkey '(alt "c:23") "rofi -show windowcd")
(xbindkey '(control alt "c:44") "xdotool set_desktop --relative -- -1")
(xbindkey '(control alt "c:47") "xdotool set_desktop --relative -- 1")
;(xbindkey '(control alt "c:57") "xdotool set_desktop --relative -- -1")
(xbindkey '(control alt m) "wmctrl -r :ACTIVE: -b toggle,maximized_vert,maximized_horz;")
;(xbindkey-function '(control alt b)
;									 (lambda ()
;										 (run-command "wmctrl -R 0x02206127")
;										 (run-command "wmctrl -b toggle,shaded -i -r 0x02206127")))

;(xbindkey '(control alt f) "xdg-open ~/Downloads")
(xbindkey '(control alt "c:41") "/home/como/Desktop/Telegram/Telegram")
;; Control+Shift+a  release event starts rxvt
;;(xbindkey '(release control shift a) "rxvt")

;; Control + mouse button 2 release event starts rxvt
;;(xbindkey '(releace control "b:1") "rxvt")


;; Extra features
;; (xbindkey-function '(control a)
;;      	   (lambda ()
;;      	     (display "Hello from Scheme!")
;;      	     (newline)))

;; (xbindkey-function '(shift p)
;;      	   (lambda ()
;;      	     (run-command "xterm")))

;; run-command "xprop -id `xdotool getactivewindow` |awk '/WM_CLASS/{print $4}'"
;; (xbindkey-function '(control e)
;;      		   (lambda ()
;; 		     (let ((b (run-command "xdotool getactivewindow")))
;; 			   ;; (a (run-command (string-append "xprop -id" b  "|awk '/WM_CLASS/{print $4}'"))))
;; 		       (run-command (string-append "zenity --info --text=\"" b "\"")))))

;; Double click test

(let ((count 0))
  (xbindkey-function '(control "c:12")
		     (lambda () 
		       (begin
			 (set! count 0)
			 (run-command "setxkbmap -layout \"us\""))))
  (xbindkey-function '(control space)									  
      		     (lambda ()
		       (begin 
			 (if (= count 0)
			     (begin
			       (set! count 1)
			       (run-command "setxkbmap -layout \"us(cmk_ed_dh)\" -option ctrl:nocaps"))
			     (if (= count 1)
				 (begin
				   (set! count 0)
				   (run-command "setxkbmap -layout \"rulemak\" -option ctrl:nocaps"))))))) 
)
;; Time double click test:
;;  - short double click -> run an xterm
;;  - long  double click -> run an rxvt
;; (xbindkey-function '(shift w)
;;      	   (let ((time (current-time))
;;      		 (count 0))
;;      	     (lambda ()
;;      	       (set! count (+ count 1))
;;      	       (if (> count 1)
;;      		   (begin
;;      		    (if (< (- (current-time) time) 1)
;;      			(run-command "xterm")
;;      			(run-command "rxvt"))
;;      		    (set! count 0)))
;;      	       (set! time (current-time)))))


;; Chording keys test: Start differents program if only one key is
;; pressed or another if two keys are pressed.
;; If key1 is pressed start cmd-k1
;; If key2 is pressed start cmd-k2
;; If both are pressed start cmd-k1-k2 or cmd-k2-k1 following the
;;   release order
;; (define (define-chord-keys key1 key2 cmd-k1 cmd-k2 cmd-k1-k2 cmd-k2-k1)
;;     "Define chording keys"
;;   (let ((k1 #f) (k2 #f))
;;     (xbindkey-function key1 (lambda () (set! k1 #t)))
;;     (xbindkey-function key2 (lambda () (set! k2 #t)))
;;     (xbindkey-function (cons 'release key1)
;;      	       (lambda ()
;;      		 (if (and k1 k2)
;;      		     (run-command cmd-k1-k2)
;;      		     (if k1 (run-command cmd-k1)))
;;      		 (set! k1 #f) (set! k2 #f)))
;;     (xbindkey-function (cons 'release key2)
;;      	       (lambda ()
;;      		 (if (and k1 k2)
;;      		     (run-command cmd-k2-k1)
;;      		     (if k2 (run-command cmd-k2)))
;;      		 (set! k1 #f) (set! k2 #f)))))


;; Example:
;;   Shift + b:1                   start an xterm
;;   Shift + b:3                   start an rxvt
;;   Shift + b:1 then Shift + b:3  start gv
;;   Shift + b:3 then Shift + b:1  start xpdf

;; (define-chord-keys '(shift "b:1") '(shift "b:3")
;;   "xterm" "rxvt" "gv" "xpdf")

;; Here the release order have no importance
;; (the same program is started in both case)
;; (define-chord-keys '(alt "b:1") '(alt "b:3")
;;   "gv" "xpdf" "xterm" "xterm")


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; End of xbindkeys guile configuration ;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
