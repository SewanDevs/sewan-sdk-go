# Generate pdf doc files with LaTeX and Kile editor
## Prerequisites
* Install Kile LaTeX editor and LaTeX : https://kile.sourceforge.io/index.php (other editors can be used : https://en.wikibooks.org/wiki/LaTeX/Installation#Editors)
* Notice the use of an uml package : tikz-uml (package file available on sequences_diagrams folder), http://www.ensta-paristech.fr/~kielbasi/tikzuml/index.php?lang=en
* install cm-super package if not already installed, it is required to prevent the fonts from beeing used as bitmap in pdf :
<pre>
$ tlmgr install cm-super
</pre>

## Build steps with Kile
* Open a .tex file with Kile
* Short cut : alt-2 -> generate .aux and .dvi files
* Short cut : alt-6 -> generate pdf file

## Docs refs
* Sequence diagrams : https://www.ibm.com/developerworks/rational/library/3101.html
