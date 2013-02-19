#!/usr/bin/env python
"""Script to boot openstack instances for iPython Notebook"""
import argparse
import ConfigParser
import os
import string
import sys
from novaclient.v1_1 import client

## Parse Command Line Arguments ##
parser = argparse.ArgumentParser(
    description="Boot a VM instance of the iPython notebook.")
parser.add_argument("-u", "--user", dest='user',
    help="The user for which this VM is being started",
    type=str, required=True)
parser.add_argument("-i", "--image", dest="image",
    help="The image to boot", type=str, required=True)
parser.add_argument("-f", "--flavor", dest="flavor",
    help="The flavor ID to boot", type=int, default=25)
parser.add_argument("-s", "--sec-group", dest="sec_group",
    help="The security group to associate with the VM",
    type=str, default="default")

args = parser.parse_args()

## Find configuration file ##
this_dir = os.path.abspath(__file__)
top_dir  = os.path.abspath(this_dir + "/..")
cfg_file = top_dir + "/config/config"
cfg = ConfigParser.ConfigParser()
if not os.path.exists(cfg_file):
    print sys.stderr, "Could not find configuration"
    sys.exit(1)
cfg.read([cfg_file])

## Configuration Options ##
shock_url = cfg.get("shock", "url")
username = cfg.get("openstack", "username")
password = cfg.get("openstack", "password")
tenant   = cfg.get("openstack", "tenant_name")
auth_url = cfg.get("openstack", "auth_url")
vm_key_name = cfg.get("openstack", "vm_key_name")
if not (shock_url and username and password and tenant and auth_url):
    print sys.stderr, "Missing configuration in config file"
    sys.exit(1)

nova = client.Client(
    username=username,
    password=password,
    tenant_name=tenant,
    auth_url=auth_url, 
    insecure=True,
)
## Validate that image and flavor exist ##
os_image  = nova.images.get(args.image)
os_flavor = nova.flavors.get(args.flavor)
os_sec_group = nova.security_groups.get(args.sec_group)
os_key = nova.keypairs.get(vm_key_name)

## Generate startup script from template & options ##
template = string.Template(
"""\
#!/bin/bash
SHOCK_USER='${user}'
SHOCK_URL='${shock_url}'

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
IPY_PATH=`which ipython`
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
    --ShockNotebookManager.shock_url=$$gSHOCK_URL --ShockNotebookManager.shock_user=$SHOCK_USER -c "${PRE_PYTHON}"
""")
script = template.safe_substitute(user=args.user, shock_url=shock_url)
script_file = tempfile.NamedTemporaryFile(delete=False)  
script_file.write(script)
script_file.close()

## Boot the instance ##
server_name = "ipnb_"+args.user
nova.servers.create(server_name, os_image, os_flavor,
    security_groups=[os_sec_group],
    userdata=script_file.name,
    key_name=vm_key_name,
)


