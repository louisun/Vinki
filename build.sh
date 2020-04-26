REPO=$(cd $(dirname $0); pwd)
COMMIT_SHA=$(git rev-parse --short HEAD)

ASSETS="false"
BINARY="false"

debugInfo () {
  echo "Repo:           $REPO"
  echo "Build assets:   $ASSETS"
  echo "Build binary:   $BINARY"
  echo "Commit:        $COMMIT_SHA"
}

buildAssets () {
  cd "$REPO" || exit
  rm -rf frontend/build
  rm -f statik/statik.go

  export CI=false

  cd "$REPO"/frontend || exit

  yarn install
  yarn run build

  if ! [ -x "$(command -v statik)" ]; then
    export CGO_ENABLED=0
    go get github.com/rakyll/statik
  fi

  cd "$REPO" || exit
  statik -src=frontend/build/  -include="*.html,*.js,*.json,*.css,*.png,*.svg,*.ico,*.woff,*.woff2,*.txt" -f
}

buildBinary () {
  cd "$REPO" || exit
  go build -a -o vinki
}

usage() {
  echo "Usage: $0 [-a] [-b]" 1>&2;
  exit 1;
}

while getopts "bacr:d" o; do
  case "${o}" in
    b)
      ASSETS="true"
      BINARY="true"
      ;;
    a)
      ASSETS="true"
      ;;
    d)
      DEBUG="true"
      ;;
    *)
      usage
      ;;
  esac
done
shift $((OPTIND-1))

if [ "$DEBUG" = "true" ]; then
  debugInfo
fi

if [ "$ASSETS" = "true" ]; then
  buildAssets
fi

if [ "$BINARY" = "true" ]; then
  buildBinary
  RESULT=$?
  if [ $RESULT -eq 0 ]; then
    echo "build binary success"
  else
    echo "build binary failed"
  fi
fi
