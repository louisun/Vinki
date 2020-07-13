usage() {
  echo "                 [Vinki Build Scripts]"
  echo "[Params]  default='--all'"
  echo
  echo "--all     build all project (frontend end backend)"
  echo "--front   build web files (including generate static go file system)"
  echo "--back    build go project"
}

if [[ $# == 0 ]];then
  usage
  exit 0
fi  

REPO=$(cd $(dirname $0); pwd)

buildFrontend () {
 cd "$REPO" || exit
 rm -rf frontend/build
 rm -f statik/statik.go

 export CI=false

 cd "$REPO"/frontend || exit

 yarn install
 echo "[INFO]    yarn install finished"
 yarn run build
 echo "[INFO]    yarn build finished"

 if ! [ -x "$(command -v statik)" ]; then
   export CGO_ENABLED=0
   go get github.com/rakyll/statik
 fi

 cd "$REPO" || exit
 statik -src=frontend/build/  -include="*.html,*.js,*.json,*.css,*.png,*.svg,*.ico,*.woff,*.woff2,*.txt" -f
 echo "[INFO]    static/statik.go generated"
 echo "[SUCCESS] build frontend done"
}

buildBackend () {
 cd "$REPO" || exit
 echo "[INFO]    go build -a -o vinki"
 go build -a -o vinki
 echo "[INFO]    vinki binary generated"
 echo "[SUCCESS] build backend done"
}


buildAll () {
  buildFrontend
  buildBackend
}

case $1 in
  -a|--all)
    buildAll
    ;;
  -f | --front)
    buildFrontend
    ;;
  -b | --back)
    buildBackend
    ;;
  *)
    usage
    ;;  
esac
