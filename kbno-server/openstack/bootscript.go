package openstack

import (
	"encoding/base64"
	//"fmt"
	"bytes"
	"github.com/MG-RAST/KBNO/kbno-server/conf"
	"io/ioutil"
	"text/template"
)

type bs struct {
	User string
	Url  string
}

const script = `#!/bin/bash
SHOCK_USER='{{.User}}'
SHOCK_URL='{{.Url}}'

if [ -n "$$g1" ]; then
    if [[ "$$g1" == *"-h"* ]]; then
        echo "Usage: start_service SHOCK_USER SHOCK_URL"
        exit
    else
        SHOCK_USER=$$g1
    fi
fi
if [ -n "$$g2" ]; then
    SHOCK_URL=$$g2
fi
SERVICE=analysis_book
DEPLOY_DIR=/kb/deployment
SERVICE_DIR=$$gDEPLOY_DIR/services/$SERVICE
KB_BIN=/kb/runtime/bin
IPY_USER=ipython
IPY_PATH=` + "`which ipython`" + `
R_HOME=/usr/local/lib/R
LD_LIBRARY_PATH=$$gR_HOME/lib:$LD_LIBRARY_PATH
PATH=$$gKB_BIN:$DEPLOY_DIR/bin:$PATH

PRE_PYTHON="import rpy2.robjects as ro; from IPython.core.display import Image, SVG; import pylab; \
from ipyMKMQ import retina, flotplot; from ipyMKMQ.ipyTools import *; \
from ipyMKMQ.analysis import get_analysis_set, Analysis, AnalysisSet; \
from ipyMKMQ.project import Project; from ipyMKMQ.collection import get_collection, Collection; \
from ipyMKMQ.metagenome import Metagenome; from ipyMKMQ.qc import QC, Drisee, NucleoProfile, Kmer, Rarefaction, merge_drisee_profile"

daemonize -v -u $$gIPY_USER -c $SERVICE_DIR/notebook -p $SERVICE_DIR/service.pid -e $SERVICE_DIR/log/start_service.log \
    -o $$gSERVICE_DIR/log/start_service.log -E R_HOME=/usr/local/lib/R -E LD_LIBRARY_PATH=$R_HOME/lib:$LD_LIBRARY_PATH \
    -E PATH=$$gKB_BIN:$DEPLOY_DIR/bin:$PATH $IPY_PATH notebook --user=$IPY_USER --pylab=inline --no-browser --port=7051 \
    --ip='*' --ipython-dir=$$gSERVICE_DIR/ipython --notebook-dir=$SERVICE_DIR/notebook --NotebookApp.verbose_crash=True \
    --NotebookApp.notebook_manager_class='IPython.frontend.html.notebook.shocknbmanager.ShockNotebookManager' \
    --ShockNotebookManager.shock_url=$$gSHOCK_URL --ShockNotebookManager.shock_user=$SHOCK_USER -c "${PRE_PYTHON}"`

func BootScript(u string) (encoding string, err error) {
	buf := bytes.NewBuffer(nil)
	t := template.Must(template.New("bs").Parse(script))
	err = t.Execute(buf, bs{User: u, Url: conf.Shock["url"].Str})
	temp, _ := ioutil.ReadAll(buf)
	return base64.StdEncoding.EncodeToString(temp), nil
}
